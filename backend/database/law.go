package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/pkg/errs"
)

const (
	latestLabel = "latest"
	masterLabel = "master"
)

func (db *_database) CreateBranch(tx *sql.Tx, upstream, ident, name, username string) (uuid.UUID, error) {
	db.logger.Log("method", "create_branch",
		"upstream", upstream,
		"ident", ident,
		"name", name,
		"username", username)

	q := `
	INSERT INTO	branches (
	  law_id,
	  user_id,
	  name
	)
	VALUES (
		(SELECT uid FROM laws WHERE upstream_id = (
			SELECT uid FROM upstreams WHERE upstreams.ident = $1 AND deleted_at IS NULL
		) AND ident = $2 AND deleted_at IS NULL),
		(SELECT uid FROM users WHERE username = $4 AND deleted_at IS NULL),
		(SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
	)
	RETURNING branches.uid
	`
	var uid uuid.UUID
	err := tx.QueryRow(
		q,
		upstream,
		ident,
		name,
		username,
	).Scan(&uid)
	if err != nil {
		return uid, err
	}
	return uid, nil
}

func (db *_database) CreateLaw(tx *sql.Tx, set *apiv1.LawSet) (uuid.UUID, error) {
	db.logger.Log(
		"method", "create_law",
		"short_title", set.Law.ShortTitle,
		"upstream_group", set.Version.UpstreamGroup,
		"author", fmt.Sprintf("%+v", set.Author),
	)

	q := `
	INSERT INTO	laws (
		upstream_id,
		user_id,
		ident,
		title,
		short_title,
		description,
		published_at,
		upstream_group_id
	)
	VALUES (
		(SELECT uid FROM upstreams WHERE ident = $1 AND deleted_at IS NULL),
		(SELECT uid FROM users WHERE username = $2 AND deleted_at IS NULL),
		$3, $4, $5, $6, $7, (
			SELECT upstream_groups.uid
			FROM upstream_groups
			INNER JOIN upstreams ON upstreams.uid = upstream_groups.upstream_id AND upstreams.deleted_at IS NULL
			WHERE upstreams.ident = $1
				AND upstream_groups.ident = LOWER($8)
				AND upstream_groups.deleted_at IS NULL
		)
	)
	RETURNING laws.uid
	`
	var when time.Time
	var err error
	var uid uuid.UUID
	if set.Version.PublishedAt == nil {
		when = time.Now()
	} else {
		when, err = ptypes.Timestamp(set.Version.PublishedAt)
		if err != nil {
			return uid, err
		}
	}
	err = tx.QueryRow(
		q,
		set.Law.Upstream,
		set.Author.Username,
		set.Law.Ident,
		set.Law.Title,
		set.Law.ShortTitle,
		set.Law.Description,
		when,
		set.Version.UpstreamGroup,
	).Scan(&uid)
	if err != nil {
		return uid, err
	}

	// TODO insert law groups

	return uid, nil
}

func (db *_database) CreateFirstVersion(set *apiv1.LawSet) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	// insert law in law table
	_, err = db.CreateLaw(tx, set)
	if err != nil {
		tx.Rollback()
		return err
	}
	//set.Law.Uid = lawUid.String()

	// create the master branch
	set.Branch.Name = "master"
	uid, err := db.CreateBranch(
		tx,
		set.Law.Upstream,
		set.Law.Ident,
		set.Branch.Name,
		set.Author.Username,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	set.Branch.Uid = uid.String()
	_, err = db.CreateVersion(tx, set)
	if err != nil {
		tx.Rollback()
		return err
	}

	// create the user's branch
	set.Branch.Name = set.Author.Username
	uid, err = db.CreateBranch(
		tx,
		set.Law.Upstream,
		set.Law.Ident,
		set.Branch.Name,
		set.Author.Username,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	set.Branch.Uid = uid.String()
	_, err = db.CreateVersion(tx, set)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (db *_database) UpdateVersion(tx *sql.Tx, set *apiv1.LawSet) (*apiv1.LawSet, error) {
	db.logger.Log("method", "update_version",
		"uid", set.Version.Uid,
		"hash", set.Version.Hash,
	)

	//q := `
	//UPDATE versions AS v
	//SET hash = $1,
	//updated_at = $6
	//FROM versions
	//INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
	//INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	//INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	//WHERE upstreams.ident = $2
	//AND laws.ident = $3
	//AND branches.name = (SELECT uid FROM users WHERE username = $4 AND deleted_at IS NULL)::varchar
	//AND versions.number = $5
	//AND versions.deleted_at IS NULL
	//`
	q := `
	UPDATE versions
	SET hash = $1,
		updated_at = $3
	WHERE uid = $2
	`
	res, err := tx.Exec(
		q,
		set.Version.Hash,
		set.Version.Uid,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errs.ErrNotFound
	}
	return set, nil
}

func (db *_database) CreateVersion(tx *sql.Tx, set *apiv1.LawSet) (*apiv1.LawSet, error) {
	db.logger.Log("method", "create_version",
		"upstream", set.Law.Upstream,
		"ident", set.Law.Ident,
		"branch", set.Branch.Name,
		"username", set.Author.Username,
		"hash", set.Version.Hash,
		"msg", set.Version.Msg,
		"when", set.Version.PublishedAt,
		"upstream_group", set.Version.UpstreamGroup,
	)

	q := `
	INSERT INTO	versions (
		branch_id,
		user_id,
		hash,
		message,
		published_at,
		number,
		tag_1,
		tag_2,
		tag_3,
		tag_4,
		upstream_group_id
	) VALUES (
		(
			SELECT CASE WHEN $13 = '' THEN (
				SELECT branches.uid
				FROM branches
				INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
				INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
				WHERE upstreams.ident = $1
					AND branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
					AND	laws.ident = $2
					AND	branches.deleted_at IS NULL
				)
			ELSE uuid($13)
			END
		),
		(SELECT users.uid FROM users WHERE username = $4 AND deleted_at IS NULL),
		$5,
		$6,
		$7,
		(
			SELECT COALESCE((
				SELECT versions.number
				FROM versions
				INNER JOIN branches ON versions.branch_id = branches.uid AND branches.deleted_at IS NULL
				INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
				INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
				WHERE branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
					AND upstreams.ident = $1
					AND laws.ident = $2
					AND versions.deleted_at IS NULL
				ORDER BY versions.number DESC
				LIMIT 1
			), 0) + 1
		), $8, $9, $10, $11, (
			SELECT upstream_groups.uid
			FROM upstream_groups
			INNER JOIN upstreams ON upstreams.uid = upstream_groups.upstream_id AND upstreams.deleted_at IS NULL
			WHERE upstreams.ident = $1
				AND upstream_groups.ident = LOWER($12)
				AND upstream_groups.deleted_at IS NULL
				
		)

	)
	RETURNING versions.number,
		versions.uid
	`
	var version uint32
	var err error
	var when time.Time
	var uid uuid.UUID
	if set.Version.PublishedAt == nil {
		when = time.Now()
	} else {
		when, err = ptypes.Timestamp(set.Version.PublishedAt)
		if err != nil {
			return nil, err
		}
	}
	err = tx.QueryRow(
		q,
		set.Law.Upstream,
		set.Law.Ident,
		set.Branch.Name,
		set.Author.Username,
		set.Version.Hash,
		set.Version.Msg,
		when,
		set.Version.Tag_1,
		set.Version.Tag_2,
		set.Version.Tag_3,
		set.Version.Tag_4,
		set.Version.UpstreamGroup,
		set.Branch.Uid,
	).Scan(&version, &uid)
	if err != nil {
		return nil, err
	}
	set.Version.Version = version
	set.Version.Uid = uid.String()

	err = db.RefreshSearchView(tx)
	if err != nil {
		return nil, err
	}

	return set, nil
}

func (db *_database) RefreshSearchView(tx *sql.Tx) error {
	q := `
	REFRESH MATERIALIZED VIEW search_index
	`
	_, err := tx.Exec(q)
	return err
}

func (db *_database) GetUserVersionsCount(username string) (int, error) {
	q := `
	SELECT COUNT(versions.uid)
	FROM versions
	INNER JOIN users ON users.uid = versions.user_id AND users.deleted_at IS NULL
	WHERE users.username = $1
		AND versions.deleted_at IS NULL
	`
	var c int
	err := db.conn.QueryRow(q, username).Scan(&c)
	switch {
	case err == sql.ErrNoRows:
		return 0, errs.ErrNotFound
	case err != nil:
		return 0, err
	}
	return c, err
}

/*

	(
		SELECT 1
		FROM votes
		INNER JOIN users ON votes.user_id = $4
		WHERE votes.version_id = versions.uid
		AND users.deleted_at IS NULL
	)

*/

func (db *_database) GetVersionByLatest(user_id, upstream, ident, branch string) (*apiv1.LawSet, error) {
	db.logger.Log("method", "get_version_by_latest", "user_id", user_id, "upstream", upstream, "ident", ident, "branch", branch)
	q := `
	SELECT laws.title,
		laws.short_title,
		laws.description,
		users.username,
		CASE WHEN users.email_private = true THEN ''
			ELSE users.email
		END,
		CASE WHEN users.full_name_private = true THEN ''
			ELSE users.full_name
		END,
		users.uid,
		versions.hash,
		versions.number,
		versions.message,
		versions.published_at,
		versions.tag_1,
		versions.tag_2,
		versions.tag_3,
		versions.tag_4,
		(
			SELECT COUNT(votes.value)
			FROM votes
			WHERE votes.version_id = versions.uid
				AND votes.value = 'YES'
				AND votes.deleted_at IS NULL
		),
		(
			SELECT COUNT(votes.value)
			FROM votes
			WHERE votes.version_id = versions.uid
				AND votes.value = 'NO'
				AND votes.deleted_at IS NULL
		),
		COALESCE((
			SELECT TRUE
			FROM votes
			INNER JOIN users ON votes.user_id::varchar = $4
			WHERE votes.version_id = versions.uid
				AND votes.deleted_at IS NULL
				AND users.deleted_at IS NULL
			LIMIT 1
		), FALSE),
		COALESCE(c.count, 0)
	FROM laws
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	INNER JOIN branches ON branches.law_id = laws.uid AND branches.deleted_at IS NULL
	INNER JOIN versions ON versions.branch_id = branches.uid AND versions.deleted_at IS NULL
	INNER JOIN users ON users.uid = versions.user_id AND users.deleted_at IS NULL
	LEFT OUTER JOIN (
			SELECT COUNT(comments.uid) AS count,
				versions.uid AS uid
			FROM comments
			INNER JOIN versions on versions.uid = comments.version_id AND versions.deleted_at IS NULL
			WHERE comments.disabled = false
				AND comments.deleted_at IS NULL
			GROUP BY versions.uid
	) as c ON versions.uid = c.uid
	WHERE laws.ident = $2
		AND upstreams.ident = $1
		AND branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
		AND laws.deleted_at IS NULL
	ORDER BY versions.number DESC
	LIMIT 1
	`
	set := makeLawSet()
	set.Law.Upstream = upstream
	set.Law.Ident = ident
	set.Branch.Name = branch
	var t time.Time
	err := db.conn.QueryRow(q, upstream, ident, branch, user_id).Scan(
		&set.Law.Title,
		&set.Law.ShortTitle,
		&set.Law.Description,
		&set.Author.Username,
		&set.Author.Email,
		&set.Author.FullName,
		&set.Author.Uid,
		&set.Version.Hash,
		&set.Version.Version,
		&set.Version.Msg,
		&t,
		&set.Version.Tag_1,
		&set.Version.Tag_2,
		&set.Version.Tag_3,
		&set.Version.Tag_4,
		&set.Version.Yays,
		&set.Version.Nays,
		&set.Version.HasVoted,
		&set.Version.CommentCount,
	)
	if err == sql.ErrNoRows {
		return nil, errs.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	set.Author.PictureUrl = db.avatarURL(set.Author.Uid)
	s := int64(t.Unix())
	n := int32(t.Nanosecond())
	set.Version.PublishedAt = &timestamp.Timestamp{Seconds: s, Nanos: n}
	return set, nil
}

func (db *_database) GetVersion(user_id, upstream, ident, branch, version string) (*apiv1.LawSet, error) {
	if branch == "" {
		branch = masterLabel
	}
	if version == "" {
		version = latestLabel
	}
	if version == latestLabel {
		return db.GetVersionByLatest(user_id, upstream, ident, branch)
	}
	num, err := strconv.Atoi(version)
	if err != nil {
		return nil, err
	}
	return db.GetVersionByNumber(user_id, upstream, ident, branch, uint32(num))
}

func (db *_database) GetVersionByNumber(user_id, upstream, ident, branch string, version uint32) (*apiv1.LawSet, error) {
	db.logger.Log(
		"method", "get_version_by_number",
		"user_id", upstream,
		"upstream", upstream,
		"ident", ident,
		"branch", branch,
		"version", version,
	)
	q := `
	SELECT laws.title,
		laws.short_title,
		laws.description,
		users.username,
		CASE WHEN users.email_private = true THEN ''
			ELSE users.email
		END,
		CASE WHEN users.full_name_private = true THEN ''
			ELSE users.full_name
		END,
		users.uid,
		versions.hash,
		versions.message,
		versions.published_at,
		versions.tag_1,
		versions.tag_2,
		versions.tag_3,
		versions.tag_4,
		(
			SELECT COUNT(votes.value)
			FROM votes
			WHERE votes.version_id = versions.uid
			AND votes.value = 'YES'
			AND votes.deleted_at IS NULL
		),
		(
			SELECT COUNT(votes.value)
			FROM votes
			WHERE votes.version_id = versions.uid
			AND votes.value = 'NO'
			AND votes.deleted_at IS NULL
		),
		COALESCE((
			SELECT TRUE
			FROM votes
			INNER JOIN users ON votes.user_id::varchar = $5
			WHERE votes.version_id = versions.uid
				AND votes.deleted_at IS NULL
				AND users.deleted_at IS NULL
			LIMIT 1
		), FALSE),
		COALESCE(c.count, 0)
	FROM versions
	INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	INNER JOIN users ON users.uid = versions.user_id AND users.deleted_at IS NULL
	LEFT OUTER JOIN (
			SELECT COUNT(comments.uid) AS count,
				versions.uid AS uid
			FROM comments
			INNER JOIN versions on versions.uid = comments.version_id AND versions.deleted_at IS NULL
			WHERE comments.disabled = false
				AND comments.deleted_at IS NULL
			GROUP BY versions.uid
	) as c ON versions.uid = c.uid
	WHERE upstreams.ident = $1
		AND	laws.ident = $2
		AND branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
		AND versions.number = $4
		AND versions.deleted_at IS NULL
	`
	set := makeLawSet()
	var t time.Time
	set.Law.Upstream = upstream
	set.Law.Ident = ident
	set.Branch.Name = branch
	set.Version.Version = version
	err := db.conn.QueryRow(q, upstream, ident, branch, version, user_id).Scan(
		&set.Law.Title,
		&set.Law.ShortTitle,
		&set.Law.Description,
		&set.Author.Username,
		&set.Author.Email,
		&set.Author.FullName,
		&set.Author.Uid,
		&set.Version.Hash,
		&set.Version.Msg,
		&t,
		&set.Version.Tag_1,
		&set.Version.Tag_2,
		&set.Version.Tag_3,
		&set.Version.Tag_4,
		&set.Version.Yays,
		&set.Version.Nays,
		&set.Version.HasVoted,
		&set.Version.CommentCount,
	)
	if err == sql.ErrNoRows {
		return nil, errs.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	s := int64(t.Unix())
	n := int32(t.Nanosecond())
	set.Author.PictureUrl = db.avatarURL(set.Author.Uid)
	set.Version.PublishedAt = &timestamp.Timestamp{Seconds: s, Nanos: n}
	return set, nil
}

func (db *_database) ListBranchVersions(upstream, ident, branch string) ([]*apiv1.Version, error) {
	db.logger.Log("method", "list_branch_versions",
		"upstream", upstream,
		"ident", ident,
		"branch", branch)

	q := `
	SELECT versions.hash,
		versions.number,
		versions.tag_1,
		versions.tag_2,
		versions.tag_3,
		versions.tag_4
	FROM versions
	INNER JOIN branches ON versions.branch_id = branches.uid AND branches.deleted_at IS NULL
	INNER JOIN laws ON branches.law_id = laws.uid AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON laws.upstream_id = upstreams.uid AND upstreams.deleted_at IS NULL
	WHERE upstreams.ident = $1
		AND laws.ident = $2
		AND	branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::"varchar"
		AND	versions.deleted_at IS NULL
	ORDER BY versions.number
	`
	rows, err := db.conn.Query(q, upstream, ident, branch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vs []*apiv1.Version
	for rows.Next() {
		v := new(apiv1.Version)
		rows.Scan(
			&v.Hash,
			&v.Version,
			&v.Tag_1,
			&v.Tag_2,
			&v.Tag_3,
			&v.Tag_4,
		)
		vs = append(vs, v)
	}
	return vs, nil
}

func (db *_database) ListUpstreamLaws(upstream, orderBy string, desc bool, pageSize, pageNum int32) (sets []*apiv1.LawSet, total int32, err error) {
	db.logger.Log("method", "list_upstream_laws", "upstream", upstream, "pageSize", pageSize, "pageNum", pageNum)
	if pageNum < 0 {
		return nil, 0, fmt.Errorf("bad pageNum: %v", pageNum)
	}
	if pageSize < 0 {
		return nil, 0, fmt.Errorf("bad pageSize: %v", pageSize)
	}

	q := `
	SELECT laws.uid,
		laws.ident,
		laws.title,
		laws.short_title,
		users.username,
		users.email,
		users.full_name,
		users.uid,
		versions.number,
		versions.published_at,
		versions.tag_1,
		versions.tag_2,
		versions.tag_3,
		versions.tag_4,
		(
			SELECT COUNT(votes.value)
			FROM votes
			WHERE votes.version_id = versions.uid
			AND votes.value = 'YES'
			AND votes.deleted_at IS NULL
		),
		(
			SELECT COUNT(votes.value)
			FROM votes
			WHERE votes.version_id = versions.uid
			AND votes.value = 'NO'
			AND votes.deleted_at IS NULL
		),
		COALESCE(c.count, 0),
		COUNT(*) OVER() AS total
	FROM versions
	INNER JOIN branches ON branches.uid = versions.branch_id
	INNER JOIN laws ON laws.uid = branches.law_id
	INNER JOIN users ON users.uid = branches.user_id
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
    INNER JOIN (
        SELECT MAX(versions.number) as latest,
        	versions.branch_id
        FROM versions
        INNER JOIN branches ON branches.uid = versions.branch_id
        INNER JOIN laws ON laws.uid = branches.law_id
        INNER JOIN users ON users.uid = branches.user_id
        INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
        WHERE upstreams.ident = $1
		AND	branches.name = (SELECT uid FROM users WHERE username = 'master')::varchar
        GROUP BY branch_id
    ) as v ON v.branch_id = versions.branch_id AND v.latest = versions.number
	LEFT OUTER JOIN (
			SELECT COUNT(comments.uid) AS count,
				versions.uid AS uid
			FROM comments
			INNER JOIN versions on versions.uid = comments.version_id AND versions.deleted_at IS NULL
			WHERE comments.disabled = false
				AND comments.deleted_at IS NULL
			GROUP BY versions.uid
	) as c ON versions.uid = c.uid
	WHERE upstreams.ident = $1
		AND	branches.name = (SELECT uid FROM users WHERE username = 'master')::varchar
	`
	switch orderBy {
	case "date":
		q += `ORDER BY versions.published_at`
	case "title":
		q += `ORDER BY laws.short_title`
	case "author":
		q += `ORDER BY users.full_name`
	default:
		q += `ORDER BY versions.published_at`
	}
	if desc {
		q += ` DESC`
	}
	q += `
	OFFSET $2
	LIMIT $3
	`
	rows, err := db.conn.Query(q, upstream, pageSize*pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		set := makeLawSet()
		set.Law.Upstream = upstream

		t := time.Time{}
		err = rows.Scan(
			&set.Law.Uid,
			&set.Law.Ident,
			&set.Law.Title,
			&set.Law.ShortTitle,
			&set.Author.Username,
			&set.Author.Email,
			&set.Author.FullName,
			&set.Author.Uid,
			&set.Version.Version,
			&t,
			&set.Version.Tag_1,
			&set.Version.Tag_2,
			&set.Version.Tag_3,
			&set.Version.Tag_4,
			&set.Version.Yays,
			&set.Version.Nays,
			&set.Version.CommentCount,
			&total,
		)
		if err != nil {
			return nil, 0, err
		}
		s := int64(t.Unix())
		n := int32(t.Nanosecond())
		set.Author.PictureUrl = db.avatarURL(set.Author.Uid)
		set.Version.PublishedAt = &timestamp.Timestamp{Seconds: s, Nanos: n}

		sets = append(sets, set)
	}
	err = rows.Err()
	return sets, total, err
}

func (db *_database) FilterUpstreamLaws(upstream, orderBy string, desc bool, pageSize, pageNum int32, search string) (sets []*apiv1.LawSet, total int32, err error) {
	db.logger.Log("method", "filter_upstream_laws", "upstream", upstream, "orderBy", orderBy, "desc", desc, "pageSize", pageSize, "pageNum", pageNum, "search", search)
	if pageNum < 0 {
		return nil, 0, fmt.Errorf("bad pageNum: %v", pageNum)
	}
	if pageSize < 0 {
		return nil, 0, fmt.Errorf("bad pageSize: %v", pageSize)
	}

	q := `
	SELECT laws.uid,
		laws.ident,
		laws.title,
		laws.short_title,
		users.username,
		users.email,
		users.full_name,
		users.uid,
		versions.number,
		versions.published_at,
		versions.tag_1,
		versions.tag_2,
		versions.tag_3,
		versions.tag_4,
		(
			SELECT COUNT(votes.value)
			FROM votes
			WHERE votes.version_id = versions.uid
				AND votes.value = 'YES'
				AND votes.deleted_at IS NULL
		),
		(
			SELECT COUNT(votes.value)
			FROM votes
			WHERE votes.version_id = versions.uid
				AND votes.value = 'NO'
				AND votes.deleted_at IS NULL
		),
		COUNT(*) OVER() AS total
	FROM versions
	INNER JOIN search_index ON search_index.uid = versions.uid
	INNER JOIN branches ON branches.uid = versions.branch_id
	INNER JOIN laws ON laws.uid = branches.law_id
	INNER JOIN users ON users.uid = branches.user_id
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
    INNER JOIN (
        SELECT MAX(versions.number) as latest,
        	versions.branch_id
        FROM versions
        INNER JOIN branches ON branches.uid = versions.branch_id
        INNER JOIN laws ON laws.uid = branches.law_id
        INNER JOIN users ON users.uid = branches.user_id
        INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
        WHERE upstreams.ident = $1
			AND	branches.name = (SELECT uid FROM users WHERE username = 'master')::varchar
        GROUP BY branch_id
    ) as v ON v.branch_id = versions.branch_id AND v.latest = versions.number
	WHERE upstreams.ident = $1
		AND	branches.name = (SELECT uid FROM users WHERE username = 'master')::varchar
		AND search_index.doc @@ plainto_tsquery($4)
	ORDER BY ts_rank(search_index.doc, plainto_tsquery($4)) DESC
	OFFSET $2
	LIMIT $3
	`
	//`
	//switch orderBy {
	//case "date":
	//q += `ORDER BY versions.published_at`
	//case "title":
	//q += `ORDER BY laws.short_title`
	//case "author":
	//q += `ORDER BY users.full_name`
	//default:
	//q += `ORDER BY versions.published_at`
	//}
	//if desc {
	//q += ` DESC`
	//}
	//q += `
	rows, err := db.conn.Query(q, upstream, pageSize*pageNum, pageSize, search)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		set := makeLawSet()
		set.Law.Upstream = upstream

		t := time.Time{}
		err = rows.Scan(
			&set.Law.Uid,
			&set.Law.Ident,
			&set.Law.Title,
			&set.Law.ShortTitle,
			&set.Author.Username,
			&set.Author.Email,
			&set.Author.FullName,
			&set.Author.Uid,
			&set.Version.Version,
			&t,
			&set.Version.Tag_1,
			&set.Version.Tag_2,
			&set.Version.Tag_3,
			&set.Version.Tag_4,
			&set.Version.Yays,
			&set.Version.Nays,
			&total,
		)
		if err != nil {
			return nil, 0, err
		}
		s := int64(t.Unix())
		n := int32(t.Nanosecond())
		set.Author.PictureUrl = db.avatarURL(set.Author.Uid)
		set.Version.PublishedAt = &timestamp.Timestamp{Seconds: s, Nanos: n}
		sets = append(sets, set)
	}
	err = rows.Err()
	return sets, total, err
}

// ListLawBranches returns a list of all branches of a law.
func (db *_database) ListLawBranches(upstream, ident string) ([]*apiv1.LawSet, error) {
	db.logger.Log("method", "list_law_branches", "upstream", upstream, "ident", ident)
	q := `
	SELECT (
		SELECT username FROM users WHERE uid::varchar = branches.name AND deleted_at IS NULL
	),
	MAX(versions.number),
	(
		SELECT 			COUNT(votes.value)
		FROM 			laws
		INNER JOIN 		branches ON branches.law_id = laws.uid
		INNER JOIN 		versions ON versions.branch_id = branches.uid
		FULL JOIN 		votes ON votes.version_id = versions.uid
		WHERE			laws.upstream_id = (SELECT uid FROM upstreams WHERE ident = $1)
			AND 			laws.ident = $2
			AND				laws.deleted_at IS NULL
			AND 			branches.deleted_at IS NULL
			AND				versions.deleted_at IS NULL
			AND 			votes.deleted_at IS NULL
			AND 			votes.value = 'YES'
	),
	(
		SELECT 			COUNT(votes.value)
		FROM 			laws
		INNER JOIN 		branches ON branches.law_id = laws.uid
		INNER JOIN 		versions ON versions.branch_id = branches.uid
		FULL JOIN 		votes ON votes.version_id = versions.uid
		WHERE			laws.upstream_id = (SELECT uid FROM upstreams WHERE ident = $1)
			AND 			laws.ident = $2
			AND				laws.deleted_at IS NULL
			AND 			branches.deleted_at IS NULL
			AND				versions.deleted_at IS NULL
			AND 			votes.deleted_at IS NULL
			AND 			votes.value = 'NO'
	)
	FROM laws
	INNER JOIN branches ON branches.law_id = laws.uid AND branches.deleted_at IS NULL
	INNER JOIN versions ON versions.branch_id = branches.uid AND versions.deleted_at IS NULL
	FULL JOIN votes ON votes.version_id = versions.uid AND votes.deleted_at IS NULL
	WHERE laws.upstream_id = (SELECT uid FROM upstreams WHERE ident = $1)
		AND laws.ident = $2
		AND	laws.deleted_at IS NULL
	GROUP BY branches.name
	`
	//q := `
	//SELECT users.username,
	//SELECT username FROM users WHERE uid::varchar = branches.name AND deleted_at IS NULL
	//),
	//MAX(versions.number),
	//(
	//SELECT 			COUNT(votes.value)
	//FROM 			laws
	//INNER JOIN 		branches ON branches.law_id = laws.uid
	//INNER JOIN 		versions ON versions.branch_id = branches.uid
	//FULL JOIN 		votes ON votes.version_id = versions.uid
	//WHERE			laws.upstream_id = (SELECT uid FROM upstreams WHERE ident = $1)
	//AND 			laws.ident = $2
	//AND				laws.deleted_at IS NULL
	//AND 			branches.deleted_at IS NULL
	//AND				versions.deleted_at IS NULL
	//AND 			votes.deleted_at IS NULL
	//AND 			votes.value = 'YES'
	//),
	//(
	//SELECT 			COUNT(votes.value)
	//FROM 			laws
	//INNER JOIN 		branches ON branches.law_id = laws.uid
	//INNER JOIN 		versions ON versions.branch_id = branches.uid
	//FULL JOIN 		votes ON votes.version_id = versions.uid
	//WHERE			laws.upstream_id = (SELECT uid FROM upstreams WHERE ident = $1)
	//AND 			laws.ident = $2
	//AND				laws.deleted_at IS NULL
	//AND 			branches.deleted_at IS NULL
	//AND				versions.deleted_at IS NULL
	//AND 			votes.deleted_at IS NULL
	//AND 			votes.value = 'NO'
	//)
	//FROM laws
	//INNER JOIN branches ON branches.law_id = laws.uid
	//INNER JOIN versions ON versions.branch_id = branches.uid
	//FULL JOIN votes ON votes.version_id = versions.uid
	//WHERE laws.upstream_id = (SELECT uid FROM upstreams WHERE ident = $1)
	//AND laws.ident = $2
	//AND	laws.deleted_at IS NULL
	//AND	versions.deleted_at IS NULL
	//AND votes.deleted_at IS NULL
	//GROUP BY branches.name
	//`
	rows, err := db.conn.Query(q, upstream, ident)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sets []*apiv1.LawSet
	for rows.Next() {
		set := &apiv1.LawSet{
			Law: &apiv1.Law{
				Upstream: upstream,
			},
			Branch:  &apiv1.Branch{},
			Version: &apiv1.Version{},
			Author:  &apiv1.Author{},
		}
		rows.Scan(
			&set.Branch.Name,
			&set.Version.Version,
			&set.Version.Yays,
			&set.Version.Nays,
		)
		sets = append(sets, set)
	}
	return sets, nil
}

func (db *_database) ListUserLaws(username, orderBy string, desc bool, pageSize, pageNum int32) ([]*apiv1.LawSet, int32, error) {
	q := `
	SELECT laws.ident,
		laws.title,
		laws.short_title,
		laws.description,
		users.username,
		users.uid,
		upstreams.ident,
		versions.number,
		versions.published_at,
		versions.tag_1,
		versions.tag_2,
		versions.tag_3,
		versions.tag_4,
		COUNT(*) OVER() AS total
	FROM laws
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
	INNER JOIN branches ON laws.uid = branches.law_id AND branches.deleted_at IS NULL
	INNER JOIN versions ON branches.uid = versions.branch_id AND versions.deleted_at IS NULL
	INNER JOIN users ON users.uid = versions.user_id AND users.deleted_at IS NULL
	INNER JOIN (
		SELECT MAX(versions.number) as latest,
			versions.branch_id
		FROM versions
		INNER JOIN users ON users.uid = versions.user_id
		INNER JOIN branches ON branches.uid = versions.branch_id
		WHERE LOWER(users.username) = LOWER($1)
			AND branches.name <> (SELECT uid FROM users WHERE username = 'master')::varchar
		GROUP BY branch_id
	) as v ON v.branch_id = versions.branch_id AND v.latest = versions.number
	WHERE LOWER(users.username) = LOWER($1)
		AND	laws.deleted_at IS NULL
	`
	switch orderBy {
	case "date":
		q += `ORDER BY published_at`
	default:
		q += `ORDER BY published_at`
	}
	if desc {
		q += ` DESC`
	}
	q += `
	OFFSET $2
	LIMIT $3
	`
	rows, err := db.conn.Query(q, username, pageSize*pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var sets []*apiv1.LawSet
	var total int32
	for rows.Next() {
		set := makeLawSet()
		var t time.Time
		rows.Scan(
			&set.Law.Ident,
			&set.Law.Title,
			&set.Law.ShortTitle,
			&set.Law.Description,
			&set.Author.Username,
			&set.Author.Uid,
			&set.Law.Upstream,
			&set.Version.Version,
			&t,
			&set.Version.Tag_1,
			&set.Version.Tag_2,
			&set.Version.Tag_3,
			&set.Version.Tag_4,
			&total,
		)
		s := int64(t.Unix())
		n := int32(t.Nanosecond())
		set.Version.PublishedAt = &timestamp.Timestamp{Seconds: s, Nanos: n}
		set.Author.Username = username
		set.Author.PictureUrl = db.avatarURL(set.Author.Uid)
		sets = append(sets, set)
	}
	return sets, total, nil
}

func (db *_database) UpdateLaw(lm *apiv1.Law) (err error) {
	return
}

func makeLawSet() *apiv1.LawSet {
	return &apiv1.LawSet{
		Law:     &apiv1.Law{},
		Branch:  &apiv1.Branch{},
		Version: &apiv1.Version{},
		Author:  &apiv1.Author{},
	}
}

func hasBlank(vars ...string) bool {
	for _, v := range vars {
		if v == "" {
			return true
		}
	}
	return false
}
