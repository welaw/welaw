package services

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/securecookie"
	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/backend/database"
	"github.com/welaw/welaw/backend/filesystem"
	"github.com/welaw/welaw/backend/vcs"
	"github.com/welaw/welaw/pkg/oauth"
)

type Service interface {
	// admin
	GetServerStats(context.Context) (*apiv1.ServerStats, error)
	LoadRepos(context.Context, *apiv1.LoadReposOptions) (*apiv1.LoadReposReply, error)
	SaveRepos(context.Context, *apiv1.SaveReposOptions) (*apiv1.SaveReposReply, error)
	// auth
	LoggedInCheck(context.Context) (*apiv1.User, error)
	Login(ctx context.Context, returnURL, provider string) (*http.Cookie, string, error)
	LoginCallback(ctx context.Context, code, state string) (*http.Cookie, *http.Cookie, string, error)
	LoginAs(context.Context, *apiv1.User) error
	Logout(context.Context) error
	// ballot
	CreateVote(context.Context, *apiv1.Vote, *apiv1.CreateVoteOptions) (*apiv1.Vote, error)
	CreateVotes(context.Context, []*apiv1.Vote, *apiv1.CreateVotesOptions) ([]*apiv1.Vote, error)
	DeleteVote(ctx context.Context, upstream, ident string, opts *apiv1.DeleteVoteOptions) error
	GetVote(ctx context.Context, upstream, ident string, opts *apiv1.GetVoteOptions) (*apiv1.Vote, error)
	ListVotes(context.Context, *apiv1.ListVotesOptions) (*apiv1.ListVotesReply, error)
	UpdateVote(context.Context, *apiv1.Vote, *apiv1.UpdateVoteOptions) (*apiv1.Vote, error)
	// law
	CreateAnnotation(context.Context, *apiv1.Annotation) (string, error)
	DeleteAnnotation(context.Context, string) error
	ListAnnotations(context.Context, *apiv1.ListAnnotationsOptions) ([]*apiv1.Annotation, int, error)
	ListComments(context.Context, *apiv1.ListCommentsOptions) ([]*apiv1.Comment, int, error)
	CreateComment(context.Context, *apiv1.Comment) (*apiv1.Comment, error)
	DeleteComment(context.Context, string) error
	GetComment(ctx context.Context, opts *apiv1.GetCommentOptions) (*apiv1.Comment, error)
	UpdateComment(context.Context, *apiv1.Comment) (*apiv1.Comment, error)
	LikeComment(context.Context, *apiv1.LikeCommentOptions) error
	CreateLaw(context.Context, *apiv1.LawSet, *apiv1.CreateLawOptions) (*apiv1.CreateLawReply, error)
	CreateLaws(context.Context, []*apiv1.LawSet, *apiv1.CreateLawsOptions) ([]*apiv1.LawSet, error)
	DeleteLaw(ctx context.Context, upstream, ident string, opts *apiv1.DeleteLawOptions) error
	DiffLaws(ctx context.Context, upstream, ident string, opts *apiv1.DiffLawsOptions) (*apiv1.DiffLawsReply, error)
	GetLaw(ctx context.Context, upstream, ident string, opts *apiv1.GetLawOptions) (*apiv1.GetLawReply, error)
	ListLaws(context.Context, *apiv1.ListLawsOptions) (*apiv1.ListLawsReply, error)
	UpdateLaw(context.Context, *apiv1.LawSet, *apiv1.UpdateLawOptions) error
	// upstreams
	CreateUpstream(context.Context, *apiv1.Upstream) error
	GetUpstream(context.Context, string) (*apiv1.Upstream, error)
	ListUpstreams(context.Context) ([]*apiv1.Upstream, error)
	UpdateUpstream(context.Context, *apiv1.Upstream) error
	// users
	CreateUser(context.Context, *apiv1.User) (*apiv1.User, error)
	CreateUsers(context.Context, []*apiv1.User) ([]*apiv1.User, error)
	DeleteUser(context.Context, string) error
	GetUser(context.Context, *apiv1.GetUserOptions) (*apiv1.User, error)
	ListUsers(context.Context, *apiv1.ListUsersOptions) ([]*apiv1.User, int, error)
	UpdateUser(context.Context, string, *apiv1.UpdateUserOptions) (*apiv1.User, error)
	UploadAvatar(ctx context.Context, opts *apiv1.UploadAvatarOptions) error
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
