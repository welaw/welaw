package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/pkg/errs"
)

func (db *_database) CreateAnnotation(ann *apiv1.Annotation) (string, error) {
	db.logger.Log("method", "create_annotation", "annotation", fmt.Sprintf("%+v", ann))

	q := `
	INSERT INTO	annotations (
		comment_id,
		text,
		quote,
		ranges
	)
	VALUES (
		$1, $2, $3, $4
	)
	RETURNING id
	`
	s, err := json.Marshal(ann.Ranges)
	if err != nil {
		return "", err
	}
	var id uuid.UUID
	err = db.conn.QueryRow(
		q,
		ann.CommentId,
		ann.Text,
		ann.Quote,
		json.RawMessage(s),
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (db *_database) DeleteAnnotationById(id string) error {
	q := `
	UPDATE annotations
	SET deleted_at = $2
	WHERE id = $1
	`
	res, err := db.conn.Exec(
		q,
		id,
		time.Now(),
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.NotFound("annotation not found: %s", id)
	}
	return nil
}

func (db *_database) ListAnnotations(comment_id string) ([]*apiv1.Annotation, int, error) {
	db.logger.Log("method", "comment_search")

	q := `
	SELECT text,
	  	quote,
		ranges,
		id,
		COUNT(*) OVER() AS total
	FROM annotations
	WHERE comment_id = $1
		AND deleted_at IS NULL
	`
	rows, err := db.conn.Query(q, comment_id)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var annotations []*apiv1.Annotation
	var total int
	for rows.Next() {
		a := new(apiv1.Annotation)
		b := []byte{}
		ranges := []*apiv1.AnnotationRange{}
		rows.Scan(
			&a.Text,
			&a.Quote,
			&b,
			&a.Id,
			&total,
		)
		err = json.Unmarshal(b, &ranges)
		if err != nil {
			return nil, 0, err
		}
		a.Ranges = ranges
		annotations = append(annotations, a)
	}
	return annotations, total, err
}

func (db *_database) CreateComment(user_id string, c *apiv1.Comment) (*apiv1.Comment, error) {
	db.logger.Log("method", "create_comment", "comment", fmt.Sprintf("%+v", c))

	q := `
	INSERT INTO	comments (
	  	user_id,
		version_id,
		comment
	)
	VALUES (
		$1,
		(
			SELECT versions.uid
			FROM versions
			INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
			INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
			INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
			WHERE upstreams.ident = $2
				AND laws.ident = $3
				AND branches.name = (SELECT uid FROM users WHERE username = $4 AND deleted_at IS NULL)::varchar
				AND versions.number = $5
				AND versions.deleted_at IS NULL
		), $6
	)
	RETURNING uid
	`
	var uid uuid.UUID
	err := db.conn.QueryRow(
		q,
		user_id,
		c.Upstream,
		c.Ident,
		c.Branch,
		c.Version,
		c.Comment,
	).Scan(&uid)
	if err != nil {
		return nil, err
	}
	c.Uid = uid.String()
	return c, err
}

func (db *_database) DeleteComment(uid string) error {
	db.logger.Log("method", "delete_comment", "uid", uid)
	if uid == "" {
		return errs.BadRequest("uid not found")
	}
	q := `
	UPDATE comments
	SET deleted_at = $2
	WHERE uid = $1
		AND deleted_at IS NULL
	`
	res, err := db.conn.Exec(q, uid, time.Now())
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (db *_database) GetCommentByUid(uid string) (*apiv1.Comment, error) {
	db.logger.Log("method", "get_comment_by_uid", "uid", uid)
	if uid == "" {
		return nil, errs.BadRequest("uid not found")
	}
	q := `
	SELECT comments.comment,
		comments.disabled,
		users.username,
		users.uid,
		CASE WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		COUNT(annotations.uid)
	FROM comments
	LEFT JOIN annotations ON annotations.comment_id = comments.uid AND annotations.deleted_at IS NULL
	INNER JOIN users ON users.uid = comments.user_id AND users.deleted_at IS NULL
	WHERE comments.uid = $1
		AND comments.deleted_at IS NULL
	GROUP BY comments.comment,
		comments.disabled,
		users.username,
		users.uid,
		users.full_name,
		users.full_name_private
	`
	var c apiv1.Comment
	var u apiv1.User
	err := db.conn.QueryRow(q, uid).Scan(
		&c.Comment,
		&c.Disabled,
		&u.Username,
		&u.Uid,
		&u.FullName,
		&c.AnnotationCount,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	c.Uid = uid
	u.PictureUrl = db.avatarURL(u.Uid)
	c.User = &u
	return &c, nil
}

func (db *_database) GetCommentByUserVersion(username, upstream, ident, branch string, version int32) (*apiv1.Comment, error) {
	db.logger.Log("method", "get_comment", "upstream", upstream, "ident", ident, "username", username)
	q := `
	SELECT comments.uid,
	  	comments.comment,
		comments.disabled,
		users.username,
		users.uid,
		CASE WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		COUNT(annotations.uid)
	FROM comments
	LEFT JOIN annotations ON annotations.comment_id = comments.uid AND annotations.deleted_at IS NULL
	INNER JOIN versions ON versions.uid = comments.version_id AND versions.deleted_at IS NULL
	INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	INNER JOIN users ON users.uid = comments.user_id AND users.deleted_at IS NULL
	WHERE upstreams.ident = $1
		AND laws.ident = $2
		AND branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
		AND versions.number = $4
		AND users.username = $5
		AND comments.deleted_at IS NULL
	GROUP BY comments.uid,
		comments.comment,
		comments.disabled,
		users.username,
		users.uid,
		users.full_name,
		users.full_name_private
	`
	var uid uuid.UUID
	var c apiv1.Comment
	var u apiv1.User
	err := db.conn.QueryRow(q, upstream, ident, branch, version, username).Scan(
		&uid,
		&c.Comment,
		&c.Disabled,
		&u.Username,
		&u.Uid,
		&u.FullName,
		&c.AnnotationCount,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	u.PictureUrl = db.avatarURL(u.Uid)
	c.Uid = uid.String()
	c.User = &u
	return &c, nil
}

func (db *_database) UpdateComment(comment *apiv1.Comment) (c *apiv1.Comment, err error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return
	}

	// delete old
	q := `
	UPDATE comments
	SET deleted_at = $2
	WHERE uid = $1
		AND	deleted_at IS NULL
	RETURNING key
	`
	var key uuid.UUID
	err = tx.QueryRow(q, comment.Uid, time.Now()).Scan(&key)
	if err != nil {
		tx.Rollback()
		return
	}

	// create new
	q = `
	INSERT INTO	comments (
		uid,
		version_id,
		user_id,
		comment,
		disabled
	)
	VALUES ($2, (
		SELECT version_id FROM comments WHERE key = $1
	), (
		SELECT user_id FROM comments WHERE key = $1
	), (
		SELECT CASE WHEN $3 = '' THEN comment
			ELSE $3
		END
		FROM comments WHERE key = $1
	), false)
	`
	res, err := tx.Exec(
		q,
		key,
		comment.Uid,
		comment.Comment,
	)
	if err != nil {
		tx.Rollback()
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if count == 0 {
		tx.Rollback()
		return nil, errs.ErrNotFound
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}
	return comment, nil
}

func (db *_database) LikeComment(comment_id, user_id string) error {
	q := `
	INSERT INTO	comment_likes (
		comment_id,
	  	user_id
	)
	VALUES (
		$1, $2
	)
	`
	res, err := db.conn.Exec(
		q,
		comment_id,
		user_id,
	)
	if err != nil && strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
		return errs.ErrConflict
	} else if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.ErrConflict
	}
	return nil
}

func (db *_database) ListCommentsByUsername(userID, username string) ([]*apiv1.Comment, int, error) {
	db.logger.Log("method", "list_comments_by_username", "user_id", userID, "username", username)

	q := `
	SELECT comments.uid,
		comments.comment,
		COUNT(annotations.uid),
		COUNT(*) OVER() AS total
	FROM comments
	INNER JOIN users ON users.uid = comments.user_id AND users.deleted_at IS NULL
	LEFT JOIN annotations ON annotations.comment_id = comments.uid AND annotations.deleted_at IS NULL
	WHERE users.username = $1
		AND comments.disabled = false
		AND comments.deleted_at IS NULL
	GROUP BY comments.uid,
		comments.comment
	`
	rows, err := db.conn.Query(q, username)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var comments []*apiv1.Comment
	var total int
	var annotationCount int32
	for rows.Next() {
		c := new(apiv1.Comment)
		rows.Scan(
			&c.Uid,
			&c.Comment,
			&annotationCount,
			&total,
		)
		c.AnnotationCount = annotationCount
		comments = append(comments, c)
	}
	return comments, total, err
}

func (db *_database) ListCommentsByVersion(userID, upstream, ident, branch, orderBy string, version, pageSize, pageNum int32, desc bool) ([]*apiv1.Comment, int, error) {
	db.logger.Log("method", "list_comments_by_version",
		"user_id", userID,
		"upstream", upstream,
		"ident", ident,
		"branch", branch,
		"version", version,
	)
	if pageNum < 0 {
		return nil, 0, fmt.Errorf("bad pageNum: %v", pageNum)
	}
	if pageSize < 1 {
		return nil, 0, fmt.Errorf("bad pageSize: %v", pageSize)
	}
	q := `
	SELECT comments.uid,
		comments.comment,
		users.username,
		users.uid,
		CASE WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		(
			SELECT COUNT(comment_likes.uid)
			FROM comment_likes
			WHERE comment_likes.comment_id = comments.uid
				AND comment_likes.deleted_at IS NULL
		) AS likes,
		(
			SELECT 1
			FROM comment_likes
			INNER JOIN users ON comment_likes.user_id = users.uid AND users.deleted_at IS NULL
			WHERE users.uid::varchar = $5
				AND comment_likes.comment_id = comments.uid
				AND comment_likes.deleted_at IS NULL
		),
		COUNT(annotations.uid) AS annotation_count,
		COUNT(*) OVER() AS total
	FROM comments
	LEFT JOIN annotations ON annotations.comment_id = comments.uid AND annotations.deleted_at IS NULL
	INNER JOIN versions ON versions.uid = comments.version_id AND versions.deleted_at IS NULL
	INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
	INNER JOIN laws ON laws.uid = branches.law_id AND laws.deleted_at IS NULL
	INNER JOIN upstreams ON upstreams.uid = laws.upstream_id AND upstreams.deleted_at IS NULL
	INNER JOIN users ON users.uid = comments.user_id AND users.deleted_at IS NULL
	WHERE upstreams.ident = $1
		AND laws.ident = $2
		AND branches.name = (SELECT uid FROM users WHERE username = $3 AND deleted_at IS NULL)::varchar
		AND versions.number = $4
		AND comments.disabled = false
		AND comments.deleted_at IS NULL
	GROUP BY comments.comment,
		comments.uid,
		comments.created_at,
		users.username,
		users.uid,
		users.full_name,
		users.full_name_private
	`
	switch orderBy {
	case "date":
		q += `ORDER BY comments.created_at`
		if desc {
			q += ` DESC`
		}
	case "votes":
		if desc {
			q += `ORDER BY dislikes`
		} else {
			q += `ORDER BY likes`
		}
	default:
		q += `ORDER BY created_at `
		if desc {
			q += ` DESC`
		}
	}
	q += `
	OFFSET $6
	LIMIT $7
	`
	rows, err := db.conn.Query(q, upstream, ident, branch, version, userID, pageSize*pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var comments []*apiv1.Comment
	var total int
	var annotationCount int32
	for rows.Next() {
		c := new(apiv1.Comment)
		u := new(apiv1.User)
		rows.Scan(
			&c.Uid,
			&c.Comment,
			&u.Username,
			&u.Uid,
			&u.FullName,
			&c.Likes,
			&c.Liked,
			&annotationCount,
			&total,
		)
		u.PictureUrl = db.avatarURL(u.Uid)
		c.User = u
		c.AnnotationCount = annotationCount
		comments = append(comments, c)
	}
	return comments, total, err
}
