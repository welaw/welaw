package permissions

import "database/sql"

// _migrations/1020_roles_init.up.sql

const (
	// admin
	OpReposSave = "repos_save"
	OpReposLoad = "repos_load"
	// ballot
	OpVoteCreate  = "vote_create"
	OpVotesCreate = "votes_create"
	OpVoteDelete  = "vote_delete"
	OpVoteUpdate  = "vote_update"
	// comment
	OpCommentCreate = "comment_create"
	OpCommentDelete = "comment_delete"
	OpCommentUpdate = "comment_update"
	// user
	OpUserCreate = "user_create"
	OpUserDelete = "user_delete"
	OpUserView   = "user_view"
	OpUserUpdate = "user_update"
	OpUserList   = "user_list"
	// roles
	OpRolesList = "roles_list"
	// law
	OpLawCreate = "law_create"
	OpLawView   = "law_view"
	OpLawUpdate = "law_update"
	OpLawDelete = "law_delete"
	// upstream
	OpUpstreamCreate = "upstream_create"
	OpUpstreamDelete = "upstream_delete"
	OpUpstreamView   = "upstream_view"
	OpUpstreamUpdate = "upstream_update"
)

const (
	RoleAdmin         = "admin"
	RoleUpstreamAdmin = "upstream-admin"
)

type Scope interface {
	HasPermision(o string) (bool, error)
	UserId() string
	Roles() []string
	Upstreams(tx *sql.Tx) []string
}
