package database

import (
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	apiv1 "github.com/welaw/welaw/api/v1"
)

type Database interface {
	// admin
	Begin() (*sql.Tx, error)
	GetServerStats() (*apiv1.ServerStats, error)
	// auth
	AuthorizeUser(username, password, operation string) (bool, error)
	CreateUserRoles(username string, role []string) error
	DeleteUserRoles(username string, role []string) error
	ListUserRoles(username string) ([]*apiv1.UserRole, error)
	HasRole(string, string) (bool, error)
	HasPermission(string, string) (bool, error)
	// ballot
	CreateVoteResult(*apiv1.VoteResult) error
	CreateVoteByLatest(*apiv1.Vote) (*apiv1.Vote, error)
	CreateVoteByVersion(*apiv1.Vote) (*apiv1.Vote, error)
	DeleteVote(username, upstream, ident, branch string, version uint32) error
	GetLawVoteSummary(upstream, ident, branch string, version uint32) (*apiv1.VoteSummary, error)
	GetVoteByLatest(username, upstream, ident, branch string) (*apiv1.Vote, error)
	GetVoteByVersion(username, upstream, ident, branch string, version uint32) (*apiv1.Vote, error)
	ListVersionVotes(upstream, ident, branch string, version, pageSize, pageNum uint32) ([]*apiv1.Vote, int32, error)
	ListUserVotes(username string, pageSize, pageNum uint32) ([]*apiv1.Vote, int32, error)
	UpdateVote(uid string, vote *apiv1.Vote) (*apiv1.Vote, error)
	// laws
	CreateAnnotation(*apiv1.Annotation) (string, error)
	DeleteAnnotationById(string) error
	ListAnnotations(string) ([]*apiv1.Annotation, int, error)
	CreateComment(string, *apiv1.Comment) (*apiv1.Comment, error)
	GetCommentByUserVersion(username, upstream, ident, branch string, version int32) (*apiv1.Comment, error)
	GetCommentByUid(uid string) (*apiv1.Comment, error)
	UpdateComment(*apiv1.Comment) (*apiv1.Comment, error)
	LikeComment(commentID, userID string) error
	CreateLaw(*sql.Tx, *apiv1.LawSet) (uuid.UUID, error)
	CreateBranch(tx *sql.Tx, up, id, br, user string) (uuid.UUID, error)
	CreateVersion(*sql.Tx, *apiv1.LawSet) (*apiv1.LawSet, error)
	CreateFirstVersion(*apiv1.LawSet) error
	DeleteComment(string) error
	FilterUpstreamLaws(upstream, orderBy string, desc bool, pageSize, pageNum int32, search string) ([]*apiv1.LawSet, int32, error)
	GetVersion(userID, upstream, ident, branch, version string) (*apiv1.LawSet, error)
	GetVersionByLatest(userID, upstream, ident, branch string) (*apiv1.LawSet, error)
	GetVersionByNumber(userID, upstream, ident, branch string, version uint32) (*apiv1.LawSet, error)
	ListCommentsByUsername(userID, username string) ([]*apiv1.Comment, int, error)
	ListCommentsByVersion(userID, upstream, ident, branch, orderBy string, version, pageSize, pageNum int32, desc bool) ([]*apiv1.Comment, int, error)
	ListBranchVersions(upstream, ident, branch string) ([]*apiv1.Version, error)
	ListLawBranches(upstream, ident string) ([]*apiv1.LawSet, error)
	ListUpstreamLaws(upstream, orderBy string, desc bool, pageSize, pageNum int32) ([]*apiv1.LawSet, int32, error)
	ListUserLaws(username, orderBy string, desc bool, pageSize, pageNum int32) ([]*apiv1.LawSet, int32, error)
	UpdateLaw(*apiv1.Law) error
	UpdateVersion(*sql.Tx, *apiv1.LawSet) (*apiv1.LawSet, error)
	// upstreams
	CreateUpstream(*apiv1.Upstream) error
	GetUpstream(string) (*apiv1.Upstream, error)
	ListUpstreams() ([]*apiv1.Upstream, error)
	ListUpstreamTags(string) ([]*apiv1.LawTag, error)
	UpdateUpstream(*apiv1.Upstream) error
	// users
	CreateUser(*apiv1.User) (*apiv1.User, error)
	CreateUserWithId(*apiv1.User) (*apiv1.User, error)
	DeleteUser(string) error
	FilterAllUsers(pageSize, pageNum int32, admin bool, search string) ([]*apiv1.User, int, error)
	GetUserById(string, bool) (*apiv1.User, error)
	GetUserByProviderId(string, bool) (*apiv1.User, error)
	GetUserByUsername(string, string, bool) (*apiv1.User, error)
	ListAllUsers(pageSize, pageNum int32, admin bool) ([]*apiv1.User, int, error)
	ListPublicUsers(pageSize, pageNum int32, admin bool) ([]*apiv1.User, int, error)
	ListUpstreamUsers(username string, pageSize, pageNum int32) ([]*apiv1.User, int, error)
	SetPassword(uid, password string) error
	UpdateLastLogin(uid string) error
	UpdateUser(string, *apiv1.User) (*apiv1.User, error)
	// db
	Close() error
}

type DatabaseConfigOptions struct {
	AvatarURL string
}

type _database struct {
	conn   *sql.DB
	logger log.Logger
	opts   *DatabaseConfigOptions
}

func (db *_database) Close() error {
	return db.conn.Close()
}

func (db *_database) Begin() (*sql.Tx, error) {
	return db.conn.Begin()
}

func NewDatabase(connStr string, logger log.Logger, opts *DatabaseConfigOptions) Database {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &_database{
		conn:   db,
		logger: logger,
		opts:   opts,
	}
}
