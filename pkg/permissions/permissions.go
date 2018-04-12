package permissions

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
	// law
	OpLawCreate = "law_create"
	OpLawView   = "law_view"
	// upstream
	OpUpstreamCreate = "upstream_create"
	OpUpstreamView   = "upstream_view"
)
