package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/proto"
)

func (db *_database) CreateUpstream(u *proto.Upstream) (err error) {
	create := `
	INSERT INTO upstreams (
		ident,
		upstream_name,
		upstream_description,
		name,
		description,
		url,
		metadata,
	)
	VALUES ($1, $2, $3, $4, $5)
	`
	md, err := json.Marshal(u.Metadata)
	if err != nil {
		return nil
	}
	_, err = db.conn.Exec(
		create,
		u.Ident,
		u.UpstreamName,
		u.UpstreamDescription,
		u.Name,
		u.Description,
		u.Url,
		&md,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *_database) GetUpstream(ident string) (*proto.Upstream, error) {
	q := `
	SELECT upstreams.uid,
		upstreams.upstream_name,
		upstreams.upstream_description,
		upstreams.name,
		upstreams.description,
		upstreams.url,
		upstreams.metadata,
		upstreams.geo_coords::varchar,
		COALESCE(users.username, ''),
		COALESCE(l.count, 0),
		COALESCE(u.count, 0)
	FROM upstreams
	LEFT JOIN users ON users.uid = upstreams.default_user AND users.deleted_at IS NULL
	LEFT OUTER JOIN (
		SELECT upstreams.uid uid,
			COUNT(laws.uid) count
		FROM laws
		INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
		WHERE upstreams.deleted_at IS NULL
		GROUP BY upstreams.uid
	) l ON l.uid = upstreams.uid
	LEFT OUTER JOIN (
		SELECT upstreams.uid uid,
			COUNT(users.uid) count
		FROM users
		INNER JOIN upstreams ON upstreams.uid = users.upstream
		WHERE upstreams.deleted_at IS NULL
			AND users.deleted_at IS NULL
		GROUP BY upstreams.uid
	) u ON u.uid = upstreams.uid
	WHERE ident = $1
		AND upstreams.deleted_at IS NULL
	`
	var raw []byte
	var u proto.Upstream
	u.Ident = ident
	err := db.conn.QueryRow(
		q,
		ident,
	).Scan(
		&u.Uid,
		&u.UpstreamName,
		&u.UpstreamDescription,
		&u.Name,
		&u.Description,
		&u.Url,
		&raw,
		&u.GeoCoords,
		&u.DefaultUser,
		&u.Laws,
		&u.Users,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	var md proto.UpstreamMetadata
	err = json.Unmarshal(raw, &md)
	if err != nil {
		return nil, err
	}

	u.Metadata = &md
	return &u, nil
}

func (db *_database) ListUpstreams() (us []*proto.Upstream, err error) {
	q := `
	SELECT ident,
		upstream_name,
		upstream_description,
		name,
		description,
		url,
		COALESCE(l.count, 0),
		COALESCE(u.count, 0)
	FROM upstreams
	LEFT OUTER JOIN (
		SELECT upstreams.uid uid,
			COUNT(laws.uid) count
		FROM laws
		INNER JOIN upstreams ON upstreams.uid = laws.upstream_id
		WHERE upstreams.deleted_at IS NULL
		GROUP BY upstreams.uid
	) l ON l.uid = upstreams.uid
	LEFT OUTER JOIN (
		SELECT upstreams.uid uid,
			COUNT(users.uid) count
		FROM users
		INNER JOIN upstreams ON upstreams.uid = users.upstream
		WHERE upstreams.deleted_at IS NULL
			AND users.deleted_at IS NULL
		GROUP BY upstreams.uid
	) u ON u.uid = upstreams.uid
	WHERE upstreams.deleted_at IS NULL
	`
	rows, err := db.conn.Query(q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := new(proto.Upstream)
		err = rows.Scan(
			&u.Ident,
			&u.UpstreamName,
			&u.UpstreamDescription,
			&u.Name,
			&u.Description,
			&u.Url,
			&u.Laws,
			&u.Users,
		)
		if err != nil {
			return
		}
		us = append(us, u)
	}
	err = rows.Err()
	return
}

//func (db *_database) ListUpstreamTags(upstream string) (tags []*proto.LawTag, err error) {
//q := `
//SELECT upstream_tags.ident,
//upstream_tags.name,
//upstream_tags.description,
//upstream_tags.ranking,
//upstream_tags.number_type
//FROM upstream_tags
//INNER JOIN upstreams ON upstreams.uid = upstream_tags.upstream_id AND upstreams.deleted_at IS NULL
//WHERE upstreams.ident = $1
//AND upstream_tags.deleted_at IS NULL
//`
//rows, err := db.conn.Query(q, upstream)
//if err != nil {
//return
//}
//defer rows.Close()
//for rows.Next() {
//t := new(proto.LawTag)
//err = rows.Scan(
//&t.Ident,
//&t.Name,
//&t.Description,
//&t.Ranking,
//&t.NumberType,
//)
//if err != nil {
//return
//}
//tags = append(tags, t)
//}
//err = rows.Err()
//return
//}

func (db *_database) UpdateUpstream(upstream *proto.Upstream) (err error) {
	// get current
	u, err := db.GetUpstream(upstream.UpstreamName)
	if err != nil {
		return err
	}

	copyUpstream(upstream, u)

	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// delete old
	q := `
	UPDATE upstreams
	SET deleted_at = $2
	WHERE uid = $1
		AND	deleted_at IS NULL
	`
	res, err := tx.Exec(
		q,
		upstream.Uid,
		time.Now(),
	)
	rows, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows == 0 {
		tx.Rollback()
		return err
	}

	md, err := json.Marshal(u.Metadata)
	if err != nil {
		tx.Rollback()
		return err
	}
	// create new
	q = `
	INSERT INTO upstreams (
		ident,
		upstream_name,
		upstream_description,
		name,
		description,
		url,
		metadata
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(
		q,
		u.Ident,
		u.UpstreamName,
		u.UpstreamDescription,
		u.Name,
		u.Description,
		u.Url,
		md,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return
}

func copyUpstream(from, to *proto.Upstream) {
	if from.Ident != "" {
		to.Ident = from.Ident
	}
	if from.UpstreamName != "" {
		to.UpstreamName = from.UpstreamName
	}
	if from.UpstreamDescription != "" {
		to.UpstreamDescription = from.UpstreamDescription
	}
	if from.Name != "" {
		to.Name = from.Name
	}
	if from.Description != "" {
		to.Description = from.Description
	}
	if from.Url != "" {
		to.Url = from.Url
	}
	if from.Metadata != nil {
		to.Metadata = from.Metadata
	}
}
