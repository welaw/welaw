package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/pkg/errs"
)

const (
	voteYes          = "YES"
	voteNo           = "NO"
	votePresent      = "PRESENT"
	voteNotPresent   = "NOT_PRESENT"
	voteResultPassed = "PASSED"
	voteResultFailed = "FAILED"
)

func (db *_database) CreateVoteResult(vr *apiv1.VoteResult) error {
	q := `
	INSERT INTO	vote_results (
		law_id,
		upstream_group_id,
		result,
		published_at
	) VALUES (
		(
			SELECT laws.uid
			FROM laws
			INNER JOIN upstreams ON uptreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
			WHERE laws.ident = $2
				AND upstreams.ident = $1
		),
		(
			SELECT upstream_groups.uid
			FROM upstream_groups
			INNER JOIN upstreams ON uptreams.uid = upstream_groups.upstream_id AND upstreams.deleted_at IS NULL
			WHERE upstream_groups.ident = $3
				AND upstreams.ident = $1
				AND deleted_at IS NULL
		),
		$4::vote_result,
		$5
	)
	RETURNING uid
	`
	var uid uuid.UUID
	err := db.conn.QueryRow(
		q,
		vr.Upstream,
		vr.Ident,
		vr.UpstreamGroup,
		vr.Result,
		vr.PublishedAt,
	).Scan(&uid)
	switch {
	case err != nil:
		return err
	}
	return nil
}

func (db *_database) CreateVoteByLatest(vote *apiv1.Vote) (*apiv1.Vote, error) {
	db.logger.Log("method", "create_vote_for_latest",
		"username", vote.Username,
		"upstream", vote.Upstream,
		"ident", vote.Ident,
		"branch", vote.Branch,
		"vote", vote.Vote,
		"comment", vote.Comment,
	)
	q := `
	INSERT INTO votes (
		user_id,
		version_id,
		value,
		comment
	) VALUES (
		( SELECT users.uid FROM	users WHERE	users.username = $1 AND users.deleted_at IS NULL ),
		(	SELECT versions.uid
			FROM versions
			INNER JOIN branches ON branches.uid = versions.branch_id
			INNER JOIN laws ON laws.uid = branches.law_id
			INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
			WHERE laws.ident = $3
			AND branches.name = (SELECT uid FROM users WHERE username = $4 AND deleted_at IS NULL)::varchar
			AND upstreams.ident = $2
			AND versions.number = (
				SELECT MAX(versions.number)
				FROM versions 
				INNER JOIN branches ON branches.uid = versions.branch_id
				INNER JOIN laws ON laws.uid = branches.law_id
				INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
				WHERE laws.ident = $3
				AND branches.name = (SELECT uid FROM users WHERE username = $4 AND deleted_at IS NULL)::varchar
				AND upstreams.ident = $2
			)
		),
		$5::vote_value,
		$6
	)
	RETURNING votes.uid
	`
	v, err := sanitizeVote(vote.Vote)
	if err != nil {
		return nil, err
	}
	var uid uuid.UUID
	err = db.conn.QueryRow(
		q,
		vote.Username,
		vote.Upstream,
		vote.Ident,
		vote.Branch,
		v,
		vote.Comment,
	).Scan(&uid)
	switch {
	case err != nil && err.Error() == "pq: null value in column \"version_id\" violates not-null constraint":
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	return vote, nil
}

func (db *_database) CreateVoteByVersion(vote *apiv1.Vote) (*apiv1.Vote, error) {
	db.logger.Log("method", "create_vote_by_version",
		"username", vote.Username,
		"upstream", vote.Upstream,
		"ident", vote.Ident,
		"branch", vote.Branch,
		"version", vote.Version,
		"vote", vote.Vote,
		"comment", vote.Comment,
	)
	v, err := sanitizeVote(vote.Vote)
	if err != nil {
		return nil, err
	}
	q := `
	INSERT INTO	votes (
		user_id,
		version_id,
		value,
		comment
	) VALUES (
		( SELECT users.uid FROM	users WHERE users.username = $1 AND users.deleted_at IS NULL ),
		(
			SELECT versions.uid
			FROM versions
			INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
			INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
			INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
			WHERE laws.ident = $3
				AND branches.name = (SELECT uid FROM users WHERE username = $4 AND deleted_at IS NULL)::varchar
				AND versions.number = $5
				AND upstreams.ident = $2
				AND versions.deleted_at IS NULL
		),
		$6::vote_value,
		$7
	)
	RETURNING votes.uid
	`
	var uid uuid.UUID
	err = db.conn.QueryRow(
		q,
		vote.Username,
		vote.Upstream,
		vote.Ident,
		vote.Branch,
		vote.Version,
		v,
		vote.Comment,
	).Scan(&uid)
	switch {
	case err != nil && err.Error() == "pq: null value in column \"version_id\" violates not-null constraint":
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	return vote, nil
}

func (db *_database) DeleteVote(username, upstream, ident, branch string, version uint32) error {
	q := `
	UPDATE votes
	SET votes.deleted_at = $6
	INNER JOIN users ON users.uid = votes.user_id AND users.deleted_at IS NULL
	INNER JOIN versions ON versions.uid = votes.version_id AND versions.deleted_at IS NULL
	INNER JOIN branches ON branches.oid = versions.branch_id AND branches.deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstream.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	WHERE users.username = $1
		AND upstreams.ident = $2
		AND laws.ident = $3
		AND branches.name = (SELECT uid FROM users WHERE username = $4 AND deleted_at IS NULL)::varchar
		AND versions.number = $5
		AND votes.deleted_at IS NULL
	`
	res, err := db.conn.Exec(
		q,
		username,
		upstream,
		ident,
		branch,
		version,
		time.Now(),
	)
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (db *_database) DeleteVoteById(lawId, userId string) error {
	q := `
	UPDATE votes
	SET deleted_at = $3
	WHERE law_id = $1
		AND user_id = $2
		AND deleted_at IS NULL
	`
	res, err := db.conn.Exec(
		q,
		lawId,
		userId,
		time.Now(),
	)
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (db *_database) GetLawVoteSummary(upstream, ident, branch string, version uint32) (*apiv1.VoteSummary, error) {
	q := `
	SELECT COUNT(votes.uid),
		votes.value
	FROM votes
	INNER JOIN versions ON versions.uid = votes.version_id AND versions.deleted_at IS NULL
	INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	INNER JOIN users ON users.uid = votes.user_id AND users.deleted_at IS NULL
	WHERE upstreams.ident = $1
		AND	laws.ident = $2
		AND branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
		AND versions.number = $4
		AND users.upstream IS NULL
		AND votes.deleted_at IS NULL
	GROUP BY votes.value
	`
	rows, err := db.conn.Query(q, upstream, ident, branch, version)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	defer rows.Close()

	var resp apiv1.VoteSummary
	var c int32
	var vote string
	for rows.Next() {
		rows.Scan(
			&c,
			&vote,
		)
		switch vote {
		case voteYes:
			resp.Yays = c
		case voteNo:
			resp.Nays = c
		case votePresent:
			resp.Present = c
		case voteNotPresent:
			resp.NotPresent = c
		default:
			return nil, fmt.Errorf("bad vote: %s", vote)
		}
	}

	q = `
	SELECT COUNT(votes.uid),
		votes.value
	FROM votes
	INNER JOIN versions ON versions.uid = votes.version_id AND versions.deleted_at IS NULL
	INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND usptreams.deleted_at IS NULL
	INNER JOIN users ON users.uid = votes.user_id AND users.deleted_at IS NULL
	WHERE upstreams.ident = $1
		AND	laws.ident = $2
		AND branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
		AND versions.number = $4
		AND users.upstream IS NOT NULL
		AND votes.deleted_at IS NULL
	GROUP BY votes.value
	`
	rows, err = db.conn.Query(q, upstream, ident, branch, version)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&c,
			&vote,
		)
		switch vote {
		case voteYes:
			resp.UpstreamYays = c
		case voteNo:
			resp.UpstreamNays = c
		case votePresent:
			resp.UpstreamPresent = c
		case voteNotPresent:
			resp.UpstreamNotPresent = c
		default:
			return nil, fmt.Errorf("bad vote: %s", vote)
		}
	}

	return &resp, err
}

func (db *_database) GetUserVoteSummary(username string) (*apiv1.VoteSummary, error) {
	q := `
	SELECT COUNT(votes.uid),
		votes.value
	FROM votes
	INNER JOIN users ON users.uid = votes.user_id AND users.deleted_at IS NULL
	WHERE users.username = $1
		AND votes.deleted_at IS NULL
	GROUP BY votes.value
	`
	rows, err := db.conn.Query(q, username)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	defer rows.Close()
	var resp apiv1.VoteSummary
	var c int32
	var vote string
	for rows.Next() {
		rows.Scan(
			&c,
			&vote,
		)
		switch vote {
		case voteYes:
			resp.Yays = c
		case voteNo:
			resp.Nays = c
		case votePresent:
			resp.Present = c
		case voteNotPresent:
			resp.NotPresent = c
		default:
			return nil, fmt.Errorf("bad vote: %s", vote)
		}
	}
	return &resp, err
}

func (db *_database) GetVoteByLatest(username, upstream, ident, branch string) (*apiv1.Vote, error) {
	q := `
	SELECT DISTINCT ON (
		votes.uid,
		votes.comment,
		votes.value,
		votes.cast_at,
		votes.version_id
	)
		votes.uid,
		votes.comment,
		votes.value,
		votes.cast_at,
		votes.version_id,
		MAX(versions.number),
		users.full_name,
		users.uid
	FROM votes
	INNER JOIN users ON users.uid = votes.user_id AND users.deleted_at IS NULL
	INNER JOIN versions ON versions.uid = votes.version_id AND versions.deleted_at IS NULL
	INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	WHERE users.username = $1
		AND upstreams.ident = $2
		AND	laws.ident = $3
		AND branches.name = (SELECT uid FROM users WHERE username = $4 AND deleted_at IS NULL)::varchar
		AND votes.deleted_at IS NULL
	GROUP BY votes.uid,
		votes.comment,
		votes.value,
		votes.cast_at,
		votes.version_id,
		users.full_name,
		users.uid
	`
	v := &apiv1.Vote{
		Username: username,
		Upstream: upstream,
		Ident:    ident,
		Branch:   branch,
	}
	var t time.Time
	var u apiv1.User
	u.Username = username
	err := db.conn.QueryRow(q, username, upstream, ident, branch).Scan(
		&v.Uid,
		&v.Comment,
		&v.Vote,
		&t,
		&v.VersionId,
		&v.Version,
		&u.FullName,
		&u.Uid,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	s := int64(t.Unix())
	n := int32(t.Nanosecond())
	v.CastAt = &timestamp.Timestamp{Seconds: s, Nanos: n}
	u.PictureUrl = db.avatarURL(u.Uid)
	v.User = &u
	return v, err
}

func (db *_database) GetVoteByVersion(username, upstream, ident, branch string, version uint32) (*apiv1.Vote, error) {
	q := `
	SELECT votes.uid,
		votes.version_id,
		votes.value,
		votes.comment,
		votes.cast_at,
		CASE WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		users.uid
	FROM votes
	INNER JOIN users ON users.uid = votes.user_id AND users.deleted_at IS NULL
	INNER JOIN versions ON versions.uid = votes.version_id AND versions.deleted_at IS NULL
	INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	WHERE users.username = $1
		AND upstreams.ident = $2
		AND	laws.ident = $3
		AND branches.name = (SELECT uid FROM users WHERE username = $4 AND deleted_at IS NULL)::varchar
		AND versions.number = $5
		AND votes.deleted_at IS NULL
	`
	v := &apiv1.Vote{
		Username: username,
		Upstream: upstream,
		Ident:    ident,
		Branch:   branch,
		Version:  version,
	}
	var t time.Time
	var u apiv1.User
	u.Username = username
	err := db.conn.QueryRow(q, username, upstream, ident, branch, version).Scan(
		&v.Uid,
		&v.VersionId,
		&v.UserId,
		&v.Vote,
		&v.Comment,
		&t,
		&u.FullName,
		&u.Uid,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	s := int64(t.Unix())
	n := int32(t.Nanosecond())
	v.CastAt = &timestamp.Timestamp{Seconds: s, Nanos: n}
	u.PictureUrl = db.avatarURL(u.Uid)
	v.User = &u
	return v, err
}

func (db *_database) ListVersionVotes(upstream, ident, branch string, version, pageNum, pageSize uint32) ([]*apiv1.Vote, int32, error) {
	db.logger.Log("method", "list_version_votes",
		"upstream", upstream,
		"ident", ident,
		"branch", branch,
		"version", version,
	)
	if pageNum < 0 {
		return nil, 0, fmt.Errorf("bad pageNum: %v", pageNum)
	}
	if pageSize < 0 {
		return nil, 0, fmt.Errorf("bad pageSize: %v", pageSize)
	}
	offset := pageSize * pageNum
	q := `
	SELECT votes.value,
		votes.comment,
		votes.cast_at,
		branches.name,
		users.username,
		CASE WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		users.uid
	FROM votes
	INNER JOIN users ON users.uid = votes.user_id AND users.deleted_at IS NULL
	INNER JOIN versions ON versions.uid = votes.version_id AND versions.deleted_at IS NULL
	INNER JOIN branches ON branches.uid = versions.branch_id AND branches. deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	WHERE upstreams.ident = $1
		AND laws.ident = $2
		AND branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
		AND versions.number = $4
		AND votes.deleted_at IS NULL
	ORDER BY votes.cast_at DESC
	OFFSET $5
	LIMIT $6
	`
	rows, err := db.conn.Query(q, upstream, ident, branch, version, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var votes []*apiv1.Vote
	var t time.Time
	var u *apiv1.User
	for rows.Next() {
		u = new(apiv1.User)
		v := &apiv1.Vote{
			Upstream: upstream,
			Ident:    ident,
			Branch:   branch,
			Version:  version,
		}
		err = rows.Scan(
			&v.Vote,
			&v.Comment,
			&t,
			&v.Branch,
			&u.Username,
			&u.FullName,
			&u.Uid,
		)
		if err != nil {
			return nil, 0, err
		}
		s := int64(t.Unix())
		n := int32(t.Nanosecond())
		v.CastAt = &timestamp.Timestamp{Seconds: s, Nanos: n}
		u.PictureUrl = db.avatarURL(u.Uid)
		v.User = u
		votes = append(votes, v)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return votes, int32(len(votes)), err
}

func (db *_database) ListUserVotes(username string, pageNum, pageSize uint32) ([]*apiv1.Vote, int32, error) {
	db.logger.Log("method", "list_user_votes", "username", username, "pageNum", pageNum, "pageSize", pageSize)
	if pageNum < 0 {
		return nil, 0, fmt.Errorf("bad pageNum: %v", pageNum)
	}
	if pageSize < 0 {
		return nil, 0, fmt.Errorf("bad pageSize: %v", pageSize)
	}
	q := `
	SELECT votes.value,
		votes.comment,
		votes.cast_at,
		upstreams.ident,
		laws.ident,
		laws.short_title,
		laws.title,
		versions.tag_1,
		versions.tag_2,
		versions.tag_3,
		versions.tag_4,
		versions.number,
		COALESCE((
			SELECT users.username
			FROM users
			WHERE uid::varchar = branches.name
			AND deleted_at IS NULL
		), '')
	FROM votes
	INNER JOIN users ON users.uid = votes.user_id AND users.deleted_at IS NULL
	INNER JOIN versions ON versions.uid = votes.version_id AND versions.deleted_at IS NULL
	INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	WHERE LOWER(users.username) = LOWER($1)
		AND votes.deleted_at IS NULL
	OFFSET $2
	LIMIT $3
	`
	rows, err := db.conn.Query(q, username, pageSize*pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var votes []*apiv1.Vote
	var t time.Time
	for rows.Next() {
		v := new(apiv1.Vote)
		l := makeLawSet()
		u := new(apiv1.User)
		v.User = u
		v.Law = l
		v.Username = username
		err = rows.Scan(
			&v.Vote,
			&v.Comment,
			&t,
			&l.Law.Upstream,
			&l.Law.Ident,
			&l.Law.ShortTitle,
			&l.Law.Title,
			&l.Version.Tag_1,
			&l.Version.Tag_2,
			&l.Version.Tag_3,
			&l.Version.Tag_4,
			&l.Version.Version,
			&l.Branch.Name,
		)
		if err != nil {
			return nil, 0, err
		}
		s := int64(t.Unix())
		n := int32(t.Nanosecond())
		v.CastAt = &timestamp.Timestamp{Seconds: s, Nanos: n}
		votes = append(votes, v)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return votes, int32(len(votes)), nil
}

//func (db *_database) UpdateVoteByVersion(username, upstream, ident, branch string, version uint32, vote *apiv1.Vote) (v *apiv1.Vote, err error) {
func (db *_database) UpdateVote(uid string, vote *apiv1.Vote) (v *apiv1.Vote, err error) {
	voteValue, err := sanitizeVote(vote.Vote)
	if err != nil {
		return
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return
	}

	// delete old
	q := `
	UPDATE votes
		SET	deleted_at = $2
	WHERE uid = $1
		AND deleted_at IS NULL
	`
	res, err := tx.Exec(
		q,
		uid,
		time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return
	}
	if rows == 0 {
		tx.Rollback()
		return nil, errs.ErrNotFound
	}

	// create new
	q = `
	INSERT INTO votes (
		version_id,
		user_id,
		value,
		comment
	)
	VALUES ($1, $2, $3, $4)
	RETURNING uid
	`
	var newUid uuid.UUID
	err = tx.QueryRow(
		q,
		vote.VersionId,
		vote.UserId,
		voteValue,
		vote.Comment,
	).Scan(&newUid)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	vote.Uid = newUid.String()
	return vote, nil
}

func sanitizeVote(v string) (string, error) {
	v = strings.ToLower(v)
	switch v {
	case "yea", "yay", "yes", "aye":
		return voteYes, nil
	case "nay", "no":
		return voteNo, nil
	case "present", "not voting":
		return votePresent, nil
	case "not-present":
		return voteNotPresent, nil
	default:
		return "", fmt.Errorf("bad vote: %s", v)
	}
}

func copyVote(to, from *apiv1.Vote) {
	if from.LawId != "" {
		to.LawId = to.LawId
	}
	if from.UserId != "" {
		to.UserId = to.UserId
	}
	if from.Vote != "" {
		to.Vote = to.Vote
	}
	if from.Comment != "" {
		to.Comment = from.Comment
	}
	if from.Upstream != "" {
		to.Upstream = from.Upstream
	}
	if from.Ident != "" {
		to.Ident = from.Ident
	}
	if from.Branch != "" {
		to.Branch = from.Branch
	}
	if from.Version != 0 {
		to.Version = from.Version
	}
	if from.Username != "" {
		to.Username = from.Username
	}
	return
}
