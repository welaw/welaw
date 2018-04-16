package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/proto"
)

// new user
func (db *_database) CreateUser(u *proto.User) (*proto.User, error) {
	q := `
	INSERT INTO users (
		username,
		full_name,
		full_name_private,
		email,
		email_private,
		biography,
		picture_url,
		provider,
		provider_id,
		upstream
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, (
		SELECT uid FROM providers WHERE ident = $8 AND deleted_at IS NULL
	), $9, (
		SELECT uid FROM upstreams WHERE ident = $10 AND deleted_at IS NULL
	))
	RETURNING uid
	`
	var uid string
	err := db.conn.QueryRow(
		q,
		u.Username,
		u.FullName,
		u.FullNamePrivate,
		u.Email,
		u.EmailPrivate,
		u.Biography,
		u.PictureUrl,
		u.Provider,
		u.ProviderId,
		u.Upstream,
	).Scan(&uid)
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			return nil, errs.ErrConflict
		}
		return nil, err
	}
	u.Uid = uid
	u.PictureUrl = db.avatarURL(uid)
	return u, nil
}

// new user
func (db *_database) CreateUserWithId(u *proto.User) (*proto.User, error) {
	q := `
	INSERT INTO users (
		uid,
		username,
		full_name,
		full_name_private,
		email,
		email_private,
		biography,
		picture_url,
		provider,
		provider_id,
		upstream
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, (
		SELECT uid FROM providers WHERE ident = $9 AND deleted_at IS NULL
	), $10, (
		SELECT uid FROM upstreams WHERE ident = $11 AND deleted_at IS NULL
	))
	RETURNING uid
	`
	var uid string
	err := db.conn.QueryRow(
		q,
		u.Uid,
		u.Username,
		u.FullName,
		u.FullNamePrivate,
		u.Email,
		u.EmailPrivate,
		u.Biography,
		u.PictureUrl,
		u.Provider,
		u.ProviderId,
		u.Upstream,
	).Scan(&uid)
	if err != nil {
		return nil, err
	}
	u.Uid = uid
	u.PictureUrl = db.avatarURL(uid)
	return u, nil
}

func (db *_database) createUpstreamUser(u *proto.User) (*proto.User, error) {
	q := `
	INSERT INTO users (
		username,
		full_name,
		full_name_private,
		email,
		email_private,
		biography,
		picture_url,
		upstream,
		provider
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, (
		SELECT uid FROM upstreams WHERE ident = $8 AND deleted_at IS NULL
	), (
		SELECT uid FROM providers WHERE ident = $9 AND deleted_at IS NULL
	))
	RETURNING uid
	`
	var uid string
	err := db.conn.QueryRow(
		q,
		u.Username,
		u.FullName,
		u.FullNamePrivate,
		u.Email,
		u.EmailPrivate,
		u.Biography,
		u.PictureUrl,
		u.Upstream,
		"welaw",
	).Scan(&uid)
	if err != nil {
		return nil, err
	}
	u.Uid = uid
	u.PictureUrl = db.avatarURL(uid)
	return u, nil
}

func (db *_database) DeleteUser(username string) (err error) {
	q := `
	UPDATE users
		SET deleted_at = $2
	WHERE username = $1
		AND deleted_at IS NULL
	`
	_, err = db.conn.Exec(q, username, time.Now())
	return
}

// only internal?
func (db *_database) GetUserById(uid string, full bool) (*proto.User, error) {
	db.logger.Log("method", "get_user_by_id", "uid", uid, "full", full)
	q := `
	SELECT users.key,
		users.provider_id,
		users.username,
		CASE WHEN $2 THEN users.full_name
			WHEN $1 = users.uid THEN users.full_name
			WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		users.full_name_private,
		CASE WHEN $2 THEN users.email
			WHEN $1 = users.uid THEN users.email
			WHEN users.email_private THEN ''
			ELSE users.email
		END,
		users.email_private,
		users.picture_url,
		users.biography,
		COALESCE(upstreams.ident, ''),
		providers.ident
	FROM users
	INNER JOIN providers ON providers.uid = users.provider AND providers.deleted_at IS NULL
	LEFT JOIN upstreams ON upstreams.uid = users.upstream AND upstreams.deleted_at IS NULL
	WHERE users.uid = $1
		AND	users.deleted_at IS NULL
	`
	u := &proto.User{Uid: uid}
	err := db.conn.QueryRow(q, uid, full).Scan(
		&u.Key,
		&u.ProviderId,
		&u.Username,
		&u.FullName,
		&u.FullNamePrivate,
		&u.Email,
		&u.EmailPrivate,
		&u.PictureUrl,
		&u.Biography,
		&u.Upstream,
		&u.Provider,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	u.PictureUrl = db.avatarURL(uid)
	return u, nil
}

// GetUserByProviderId is used by auth and will return the latest user, deleted or not.
func (db *_database) GetUserByProviderId(pid string, full bool) (*proto.User, error) {
	db.logger.Log("method", "get_user_by_provider_id", "pid", pid)
	if pid == "" {
		return nil, errs.ErrNotFound
	}
	q := `
	SELECT users.key, 
		users.uid,
		users.username,
		CASE WHEN $2 THEN users.full_name
			WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		users.full_name_private,
		CASE WHEN $2 THEN users.email
			WHEN users.email_private THEN ''
			ELSE users.email
		END,
		users.email_private,
		users.picture_url,
		users.biography,
		COALESCE(upstreams.ident, ''),
		providers.ident,
		users.deleted_at IS NOT NULL,
		CASE WHEN users.password IS NULL OR users.password = '' THEN TRUE
			ELSE FALSE
		END
	FROM users
	INNER JOIN providers ON providers.uid = users.provider AND providers.deleted_at IS NULL
	LEFT JOIN upstreams ON upstreams.uid = users.upstream AND upstreams.deleted_at IS NULL
	WHERE users.provider_id = $1
	ORDER BY users.created_at DESC
	LIMIT 1
	`
	u := &proto.User{ProviderId: pid}
	var t bool
	err := db.conn.QueryRow(q, pid, full).Scan(
		&u.Key,
		&u.Uid,
		&u.Username,
		&u.FullName,
		&u.FullNamePrivate,
		&u.Email,
		&u.EmailPrivate,
		&u.PictureUrl,
		&u.Biography,
		&u.Upstream,
		&u.Provider,
		&t,
		&u.HasPassword,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	if t {
		u.DeletedAt = "true"
	} else {
		u.DeletedAt = ""
	}
	u.PictureUrl = db.avatarURL(u.Uid)
	return u, nil
}

func (db *_database) GetUserByUsername(username string, full bool) (*proto.User, error) {
	db.logger.Log("method", "get_user_by_username", "username", username, "full", full)

	q := `
	SELECT users.key,
		users.uid,
		CASE WHEN $2 THEN users.full_name
			WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		users.full_name_private,
		CASE WHEN $2 THEN users.email
			WHEN users.email_private THEN ''
			ELSE users.email
		END,
		users.email_private,
		users.picture_url,
		users.biography,
		COALESCE(upstreams.ident, ''),
		providers.ident,
		users.provider_id,
		last_login,
		CASE WHEN users.password IS NULL OR users.password = '' THEN TRUE
			ELSE FALSE
		END
	FROM users
	INNER JOIN providers ON providers.uid = users.provider
	LEFT JOIN upstreams ON upstreams.uid = users.upstream
	WHERE LOWER(users.username) = LOWER($1)
		AND	users.deleted_at IS NULL
	`
	u := proto.User{Username: username}
	err := db.conn.QueryRow(q, username, full).Scan(
		&u.Key,
		&u.Uid,
		&u.FullName,
		&u.FullNamePrivate,
		&u.Email,
		&u.EmailPrivate,
		&u.PictureUrl,
		&u.Biography,
		&u.Upstream,
		&u.Provider,
		&u.ProviderId,
		&u.LastLogin,
		&u.HasPassword,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.NotFound("user: %s", username)
	case err != nil:
		return nil, err
	}
	u.PictureUrl = db.avatarURL(u.Uid)
	return &u, nil
}

func (db *_database) FilterAllUsers(pageSize, pageNum int32, full bool, search string) (us []*proto.User, total int, err error) {
	if pageNum < 0 {
		return nil, 0, fmt.Errorf("bad pageNum: %v", pageNum)
	}
	if pageSize < 0 {
		return nil, 0, fmt.Errorf("bad pageSize: %v", pageSize)
	}
	offset := pageSize * pageNum
	q := `
	SELECT users.uid,
		users.username,
		users.full_name,
		CASE WHEN $3 = true 
			THEN users.email
			ELSE COALESCE((SELECT email FROM users WHERE u1.uid = uid AND email_private = true AND deleted_at IS NULL), '')
		END,
		users.email_private,
		users.biography,
		users.picture_url,
		COALESCE(upstreams.ident, ''),
		LEVENSHTEIN(users.username, $3),
		LEVENSHTEIN(users.full_name, $3),
		COUNT(*) OVER() AS total
	FROM users u1
	LEFT JOIN upstreams ON upstreams.uid = users.upstream
	WHERE deleted_at IS NULL
	ORDER BY users.full_name
	OFFSET $1
	LIMIT $2
	`
	rows, err := db.conn.Query(q, offset, pageSize, full)
	if err != nil {
		return
	}
	defer rows.Close()
	var usernameLeven, fullNameLeven int
	for rows.Next() {
		u := new(proto.User)
		err = rows.Scan(
			&u.Uid,
			&u.Username,
			&u.FullName,
			&u.Email,
			&u.EmailPrivate,
			&u.Biography,
			&u.PictureUrl,
			&u.Upstream,
			&usernameLeven,
			&fullNameLeven,
			&total,
		)
		if err != nil {
			return
		}
		u.PictureUrl = db.avatarURL(u.Uid)
		us = append(us, u)
	}
	err = rows.Err()
	return
}

func (db *_database) ListAllUsers(pageSize, pageNum int32, full bool) (us []*proto.User, total int, err error) {
	if pageNum < 0 {
		return nil, 0, fmt.Errorf("bad pageNum: %v", pageNum)
	}
	if pageSize < 0 {
		return nil, 0, fmt.Errorf("bad pageSize: %v", pageSize)
	}
	q := `
	SELECT users.uid,
		users.username,
		CASE WHEN $3 THEN users.full_name
			WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		CASE WHEN $3 THEN users.email
			WHEN users.email_private THEN ''
			ELSE users.email
		END,
		users.email_private,
		users.biography,
		users.picture_url,
		COALESCE(upstreams.ident, ''),
		COUNT(*) OVER() AS total
	FROM users
	LEFT JOIN upstreams ON upstreams.uid = users.upstream
	WHERE users.deleted_at IS NULL
	ORDER BY users.full_name
	OFFSET $1
	LIMIT $2
	`
	rows, err := db.conn.Query(q, pageSize*pageNum, pageSize, full)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := new(proto.User)
		err = rows.Scan(
			&u.Uid,
			&u.Username,
			&u.FullName,
			&u.Email,
			&u.EmailPrivate,
			&u.Biography,
			&u.PictureUrl,
			&u.Upstream,
			&total,
		)
		if err != nil {
			return
		}
		u.PictureUrl = db.avatarURL(u.Uid)
		us = append(us, u)
	}
	err = rows.Err()
	return
}

func (db *_database) ListPublicUsers(pageSize, pageNum int32, full bool) (us []*proto.User, total int, err error) {
	if pageNum < 0 {
		return nil, 0, fmt.Errorf("bad pageNum: %v", pageNum)
	}
	if pageSize < 0 {
		return nil, 0, fmt.Errorf("bad pageSize: %v", pageSize)
	}
	q := `
	SELECT uid,
		username,
		CASE WHEN $3 THEN full_name
			WHEN full_name_private THEN ''
			ELSE full_name
		END,
		CASE WHEN $3 THEN email
			WHEN email_private THEN ''
			ELSE email
		END,
		biography,
		picture_url,
		last_login,
		COUNT(*) OVER() AS total
	FROM users
	WHERE deleted_at IS NULL
		AND upstream IS NULL
	ORDER BY full_name
	OFFSET $1
	LIMIT $2
	`
	rows, err := db.conn.Query(q, pageSize*pageNum, pageSize, full)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := new(proto.User)
		err = rows.Scan(
			&u.Uid,
			&u.Username,
			&u.FullName,
			&u.Email,
			&u.Biography,
			&u.PictureUrl,
			&u.LastLogin,
			&total,
		)
		if err != nil {
			return
		}
		u.PictureUrl = db.avatarURL(u.Uid)
		us = append(us, u)
	}
	err = rows.Err()
	return
}

func (db *_database) ListUpstreamUsers(upstream string, pageSize, pageNum int32) (users []*proto.User, total int, err error) {
	if pageNum < 0 {
		return nil, 0, fmt.Errorf("bad pageNum: %v", pageNum)
	}
	if pageSize < 0 {
		return nil, 0, fmt.Errorf("bad pageSize: %v", pageSize)
	}

	q := `
	SELECT users.uid,
		users.username,
		CASE WHEN users.full_name_private THEN ''
			ELSE users.full_name
		END,
		CASE WHEN users.email_private THEN ''
			ELSE users.email
		END,
		users.biography,
		users.picture_url,
		(
			SELECT COUNT(votes.value)
			FROM votes
			WHERE votes.user_id = users.uid
				AND votes.deleted_at IS NULL
		),
		(
			SELECT COUNT(versions.uid)
			FROM versions
			INNER JOIN branches ON branches.uid = versions.branch_id AND branches.deleted_at IS NULL
			WHERE versions.user_id = users.uid
				AND branches.name <> (SELECT uid FROM users WHERE username = 'master')::varchar
				AND versions.deleted_at IS NULL
		),
		COUNT(*) OVER() AS total
	FROM users
	LEFT JOIN user_roles ON user_roles.user_id = users.uid AND user_roles.deleted_at IS NULL
	LEFT JOIN roles ON roles.uid = user_roles.role_id AND roles.deleted_at IS NULL
	WHERE users.upstream = (SELECT uid FROM upstreams WHERE ident = $1 AND deleted_at IS NULL)
		AND roles.name IS NULL
		AND users.deleted_at IS NULL
	ORDER BY users.full_name
	OFFSET $2
	LIMIT $3
	`
	rows, err := db.conn.Query(q, upstream, pageSize*pageNum, pageSize)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := new(proto.User)
		u.Upstream = upstream
		err = rows.Scan(
			&u.Uid,
			&u.Username,
			&u.FullName,
			&u.Email,
			&u.Biography,
			&u.PictureUrl,
			&u.VoteCount,
			&u.LawCount,
			&total,
		)
		if err != nil {
			return
		}
		u.PictureUrl = db.avatarURL(u.Uid)
		users = append(users, u)
	}
	err = rows.Err()
	return
}

func (db *_database) UpdateLastLogin(uid string) (err error) {
	q := `
	UPDATE users
	SET last_login = $2,
		updated_at = $2
	WHERE uid = $1
		AND	deleted_at IS NULL
	`
	res, err := db.conn.Exec(q, uid, time.Now())
	if err != nil {
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return
	}
	if rows == 0 {
		return errs.ErrNotFound
	}
	return
}

func (db *_database) SetPassword(uid string, password string) (err error) {
	q := `
	UPDATE users
	SET password = crypt($2, gen_salt('bf')),
		updated_at = $3
	WHERE uid = $1
	AND	deleted_at IS NULL
	`
	res, err := db.conn.Exec(q, uid, password, time.Now())
	if err != nil {
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return
	}
	if rows == 0 {
		return errs.ErrNotFound
	}
	return
}

func (db *_database) UpdateUser(uid string, user *proto.User) (u *proto.User, err error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return
	}

	// delete old
	q := `
	UPDATE users
	SET deleted_at = $2
	WHERE uid = $1
		AND	deleted_at IS NULL
	RETURNING key
	`
	var key uuid.UUID
	err = tx.QueryRow(
		q,
		user.Uid,
		time.Now(),
	).Scan(&key)
	if err != nil {
		tx.Rollback()
		return
	}

	// create new
	q = `
	INSERT INTO	users (
		uid,
		provider,
		provider_id,
		username,
		full_name,
		full_name_private,
		email,
		email_private,
		biography,
		upstream,
		picture_url
	)
	VALUES (
		(
			SELECT CASE WHEN $2 = '' THEN uid
				ELSE $2::uuid
			END
			FROM users WHERE key = $1
		), (
			SELECT uid FROM providers WHERE ident = $3 AND deleted_at IS NULL
		), (
			SELECT CASE WHEN $4 = '' THEN provider_id
				ELSE $4
			END
			FROM users WHERE key = $1
		), (
			SELECT CASE WHEN $5 = '' THEN username
				ELSE $5
			END
			FROM users WHERE key = $1
		), (
			SELECT CASE WHEN $6 = '' THEN full_name
				ELSE $6
			END
			FROM users WHERE key = $1
		), (
			SELECT CASE WHEN $7 THEN NOT full_name_private
				ELSE full_name_private
			END
			FROM users WHERE key = $1
		), (
			SELECT CASE WHEN $8 = '' THEN email
				ELSE $8
			END
			FROM users WHERE key = $1
		), (
			SELECT CASE WHEN $9 THEN NOT email_private
				ELSE email_private
			END
			FROM users WHERE key = $1
		), (
			SELECT CASE WHEN $10 = '' THEN biography
				ELSE $10
			END
			FROM users WHERE key = $1
		), (
			SELECT uid FROM upstreams WHERE ident = $11 AND deleted_at IS NULL
		), (
			SELECT CASE WHEN $12 = '' THEN picture_url
				ELSE $12
			END
			FROM users WHERE key = $1
		)
	)
	`
	res, err := tx.Exec(
		q,
		key,
		user.Uid,
		user.Provider,
		user.ProviderId,
		user.Username,
		user.FullName,
		user.FullNamePrivate,
		user.Email,
		user.EmailPrivate,
		user.Biography,
		user.Upstream,
		user.PictureUrl,
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
	user.PictureUrl = db.avatarURL(user.Uid)
	return user, nil
}

//func copyUser(from, to *proto.User) error {
//if to.Key != "" {
//from.Key = to.Key
//}
//if to.Uid != "" {
//from.Uid = to.Uid
//}
//if to.Username != "" {
//from.Username = to.Username
//}
//if to.Email != "" {
//from.Email = to.Email
//}
//if to.EmailPrivate {
//from.EmailPrivate = to.EmailPrivate
//}
//if to.FullName != "" {
//from.FullName = to.FullName
//}
//if to.FullNamePrivate {
//from.FullNamePrivate = to.FullNamePrivate
//}
//if to.Biography != "" {
//from.Biography = to.Biography
//}
//if to.Provider != "" {
//from.Provider = to.Provider
//}
//if to.ProviderId != "" {
//from.ProviderId = to.ProviderId
//}
//return nil
//}

func (db *_database) avatarURL(uid string) string {
	return fmt.Sprintf("%s/%s.jpg", db.opts.AvatarURL, uid)
}
