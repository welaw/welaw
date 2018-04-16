package services

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/securecookie"
	"github.com/welaw/welaw/backend/database"
	"github.com/welaw/welaw/backend/filesystem"
	"github.com/welaw/welaw/backend/vcs"
	"github.com/welaw/welaw/pkg/oauth"
	"github.com/welaw/welaw/proto"
)

type Service interface {
	// admin
	GetServerStats(context.Context) (*proto.ServerStats, error)
	LoadRepos(context.Context, *proto.LoadReposOptions) (*proto.LoadReposReply, error)
	SaveRepos(context.Context, *proto.SaveReposOptions) (*proto.SaveReposReply, error)
	// auth
	LoggedInCheck(context.Context) (*proto.User, error)
	Login(ctx context.Context, returnURL, provider string) (*http.Cookie, string, error)
	LoginCallback(ctx context.Context, code, state string) (*http.Cookie, *http.Cookie, string, error)
	LoginAs(context.Context, *proto.User) error
	Logout(context.Context) error
	// ballot
	CreateVote(context.Context, *proto.Vote, *proto.CreateVoteOptions) (*proto.Vote, error)
	CreateVotes(context.Context, []*proto.Vote, *proto.CreateVotesOptions) ([]*proto.Vote, error)
	DeleteVote(ctx context.Context, upstream, ident string, opts *proto.DeleteVoteOptions) error
	GetVote(ctx context.Context, upstream, ident string, opts *proto.GetVoteOptions) (*proto.Vote, error)
	ListVotes(context.Context, *proto.ListVotesOptions) (*proto.ListVotesReply, error)
	UpdateVote(context.Context, *proto.Vote, *proto.UpdateVoteOptions) (*proto.Vote, error)
	// law
	CreateAnnotation(context.Context, *proto.Annotation) (string, error)
	DeleteAnnotation(context.Context, string) error
	ListAnnotations(context.Context, *proto.ListAnnotationsOptions) ([]*proto.Annotation, int, error)
	ListComments(context.Context, *proto.ListCommentsOptions) ([]*proto.Comment, int, error)
	CreateComment(context.Context, *proto.Comment) (*proto.Comment, error)
	DeleteComment(context.Context, string) error
	GetComment(ctx context.Context, opts *proto.GetCommentOptions) (*proto.Comment, error)
	UpdateComment(context.Context, *proto.Comment) (*proto.Comment, error)
	LikeComment(context.Context, *proto.LikeCommentOptions) error
	CreateLaw(context.Context, *proto.LawSet, *proto.CreateLawOptions) (*proto.CreateLawReply, error)
	CreateLaws(context.Context, []*proto.LawSet, *proto.CreateLawsOptions) ([]*proto.LawSet, error)
	DeleteLaw(ctx context.Context, upstream, ident string, opts *proto.DeleteLawOptions) error
	DiffLaws(ctx context.Context, upstream, ident string, opts *proto.DiffLawsOptions) (*proto.DiffLawsReply, error)
	GetLaw(ctx context.Context, upstream, ident string, opts *proto.GetLawOptions) (*proto.GetLawReply, error)
	ListLaws(context.Context, *proto.ListLawsOptions) (*proto.ListLawsReply, error)
	UpdateLaw(context.Context, *proto.LawSet, *proto.UpdateLawOptions) error
	// upstreams
	CreateUpstream(context.Context, *proto.Upstream) error
	GetUpstream(context.Context, string) (*proto.Upstream, error)
	ListUpstreams(context.Context) ([]*proto.Upstream, error)
	UpdateUpstream(context.Context, *proto.Upstream) error
	// users
	CreateUser(context.Context, *proto.User) (*proto.User, error)
	CreateUsers(context.Context, []*proto.User) ([]*proto.User, error)
	DeleteUser(context.Context, string) error
	GetUser(context.Context, *proto.GetUserOptions) (*proto.User, error)
	ListUsers(context.Context, *proto.ListUsersOptions) ([]*proto.User, int, error)
	UpdateUser(context.Context, string, *proto.UpdateUserOptions) (*proto.User, error)
	UploadAvatar(ctx context.Context, opts *proto.UploadAvatarOptions) error
}

//const (
//ServerLoading = iota
//ServerRunning
//ServerGoingOffline
//ServerOffline
//)

type ServerConfigOptions struct {
	LoginSuccessURL   string
	LoginFailedURL    string
	StaticDir         string
	AvatarDir         string
	StaticBucketName  string
	DefaultBucketName string
	UseSecureCookies  bool
	SigningKey        []byte
}

type service struct {
	db        database.Database
	fs        filesystem.Filesystem
	logger    log.Logger
	vc        vcs.VCS
	providers map[string]oauth.Provider
	sc        *securecookie.SecureCookie
	Opts      *ServerConfigOptions
}

func NewService(
	db database.Database,
	fs filesystem.Filesystem,
	logger log.Logger,
	vc vcs.VCS,
	providers map[string]oauth.Provider,
	sc *securecookie.SecureCookie,
	opts *ServerConfigOptions,
) Service {
	return &service{
		db:        db,
		logger:    logger,
		vc:        vc,
		providers: providers,
		sc:        sc,
		fs:        fs,
		Opts:      opts,
	}
}
