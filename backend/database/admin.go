package database

import (
	"github.com/welaw/welaw/proto"
)

func (db *_database) GetServerStats() (*proto.ServerStats, error) {
	q := `
	SELECT (
		SELECT COUNT(uid)
		FROM upstreams
		WHERE deleted_at IS NULL
	), (
		SELECT COUNT(uid)
		FROM users
		WHERE deleted_at IS NULL
	), (
		SELECT COUNT(uid)
		FROM branches
		WHERE name <> (SELECT uid FROM users WHERE username = 'master' AND deleted_at IS NULL)::varchar
		AND deleted_at IS NULL
	), (
		SELECT COUNT(versions.uid)
		FROM versions
		INNER JOIN branches ON branches.uid = versions.branch_id
		WHERE branches.name <> (SELECT uid FROM users WHERE username = 'master' AND deleted_at IS NULL)::varchar
		AND versions.deleted_at IS NULL
	), (
		SELECT COUNT(uid)
		FROM laws
		WHERE deleted_at IS NULL
	),
	COUNT (uid)
	FROM votes
	WHERE deleted_at IS NULL
	`
	stats := new(proto.ServerStats)
	err := db.conn.QueryRow(q).Scan(
		&stats.Upstreams,
		&stats.Users,
		&stats.Branches,
		&stats.Versions,
		&stats.Laws,
		&stats.Votes,
	)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
