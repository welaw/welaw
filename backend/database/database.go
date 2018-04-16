package database

import (
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/welaw/welaw/proto"
)

type Database interface {
	// admin
	Begin() (*sql.Tx, error)
	GetServerStats() (*proto.ServerStats, error)
	// auth
	AuthorizeUser(username, password, operation string) (bool, error)
	CreateUserRoles(username string, role []string) error
	DeleteUserRoles(username string, role []string) error
	ListUserRoles(username string) ([]*proto.UserRole, error)
	HasRole(string, string) (bool, error)
	GetUserAuthScope(string, string) ([]string, error)
	HasPermission(string, string) (bool, error)
	// ballot
	CreateVoteResult(*proto.VoteResult) error
	CreateVoteByLatest(*proto.Vote) (*proto.Vote, error)
	CreateVoteByVersion(*proto.Vote) (*proto.Vote, error)
	DeleteVote(username, upstream, ident, branch string, version uint32) error
	GetLawVoteSummary(upstream, ident, branch string, version uint32) (*proto.VoteSummary, error)
	GetVoteByLatest(username, upstream, ident, branch string) (*proto.Vote, error)
	GetVoteByVersion(username, upstream, ident, branch string, version uint32) (*proto.Vote, error)
	ListVersionVotes(upstream, ident, branch string, version, pageSize, pageNum uint32) ([]*proto.Vote, int32, error)
	ListUserVotes(username string, pageSize, pageNum uint32) ([]*proto.Vote, int32, error)
	UpdateVote(uid string, vote *proto.Vote) (*proto.Vote, error)
	// laws
	CreateAnnotation(*proto.Annotation) (string, error)
	DeleteAnnotationById(string) error
	ListAnnotations(string) ([]*proto.Annotation, int, error)
	CreateComment(string, *proto.Comment) (*proto.Comment, error)
	GetCommentByUserVersion(username, upstream, ident, branch string, version int32) (*proto.Comment, error)
	GetCommentByUid(uid string) (*proto.Comment, error)
	UpdateComment(*proto.Comment) (*proto.Comment, error)
	LikeComment(commentID, userID string) error
	CreateLaw(*sql.Tx, *proto.LawSet) (uuid.UUID, error)
	CreateBranch(tx *sql.Tx, up, id, br, user string) (uuid.UUID, error)
	CreateVersion(*sql.Tx, *proto.LawSet) (*proto.LawSet, error)
	CreateFirstVersion(*proto.LawSet) error
	DeleteComment(string) error
	FilterUpstreamLaws(upstream, orderBy string, desc bool, pageSize, pageNum int32, search string) ([]*proto.LawSet, int32, error)
	GetLaw(upstream, ident string) (*proto.Law, error)
	GetVersion(userID, upstream, ident, branch, version string) (*proto.LawSet, error)
	GetVersionByLatest(userID, upstream, ident, branch string) (*proto.LawSet, error)
	GetVersionByNumber(userID, upstream, ident, branch string, version uint32) (*proto.LawSet, error)
	ListCommentsByUsername(userID, username string) ([]*proto.Comment, int, error)
	ListCommentsByVersion(userID, upstream, ident, branch, orderBy string, version, pageSize, pageNum int32, desc bool) ([]*proto.Comment, int, error)
	ListBranchVersions(upstream, ident, branch string) ([]*proto.Version, error)
	ListLawBranches(upstream, ident string) ([]*proto.LawSet, error)
	ListUpstreamLaws(upstream, orderBy string, desc bool, pageSize, pageNum int32) ([]*proto.LawSet, int32, error)
	ListUserLaws(username, orderBy string, desc bool, pageSize, pageNum int32) ([]*proto.LawSet, int32, error)
	UpdateLaw(*proto.Law) error
	UpdateVersion(*sql.Tx, *proto.LawSet) (*proto.LawSet, error)
	// upstreams
	CreateUpstream(*proto.Upstream) error
	GetUpstream(string) (*proto.Upstream, error)
	ListUpstreams() ([]*proto.Upstream, error)
	//ListUpstreamTags(string) ([]*proto.LawTag, error)
	UpdateUpstream(*proto.Upstream) error
	// users
	CreateUser(*proto.User) (*proto.User, error)
	CreateUserWithId(*proto.User) (*proto.User, error)
	DeleteUser(string) error
	FilterAllUsers(pageSize, pageNum int32, admin bool, search string) ([]*proto.User, int, error)
	GetUserById(string, bool) (*proto.User, error)
	GetUserByProviderId(string, bool) (*proto.User, error)
	GetUserByUsername(string, bool) (*proto.User, error)
	ListAllUsers(pageSize, pageNum int32, admin bool) ([]*proto.User, int, error)
	ListPublicUsers(pageSize, pageNum int32, admin bool) ([]*proto.User, int, error)
	ListUpstreamUsers(username string, pageSize, pageNum int32) ([]*proto.User, int, error)
	SetPassword(uid, password string) error
	UpdateLastLogin(uid string) error
	UpdateUser(string, *proto.User) (*proto.User, error)
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
