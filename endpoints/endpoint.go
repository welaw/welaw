package endpoints

import (
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	// Admin
	GetServerStatsEndpoint endpoint.Endpoint
	LoadReposEndpoint      endpoint.Endpoint
	SaveReposEndpoint      endpoint.Endpoint
	// Auth
	LoggedInCheckEndpoint endpoint.Endpoint
	LoginEndpoint         endpoint.Endpoint
	LoginCallbackEndpoint endpoint.Endpoint
	LoginAsEndpoint       endpoint.Endpoint
	LogoutEndpoint        endpoint.Endpoint
	// Ballot
	CreateVoteEndpoint  endpoint.Endpoint
	CreateVotesEndpoint endpoint.Endpoint
	UpdateVoteEndpoint  endpoint.Endpoint
	GetVoteEndpoint     endpoint.Endpoint
	DeleteVoteEndpoint  endpoint.Endpoint
	ListVotesEndpoint   endpoint.Endpoint
	// Law
	CreateAnnotationEndpoint endpoint.Endpoint
	DeleteAnnotationEndpoint endpoint.Endpoint
	ListAnnotationsEndpoint  endpoint.Endpoint
	CreateCommentEndpoint    endpoint.Endpoint
	DeleteCommentEndpoint    endpoint.Endpoint
	GetCommentEndpoint       endpoint.Endpoint
	ListCommentsEndpoint     endpoint.Endpoint
	UpdateCommentEndpoint    endpoint.Endpoint
	LikeCommentEndpoint      endpoint.Endpoint
	CreateLawEndpoint        endpoint.Endpoint
	CreateLawsEndpoint       endpoint.Endpoint
	DeleteLawEndpoint        endpoint.Endpoint
	GetLawEndpoint           endpoint.Endpoint
	DiffLawsEndpoint         endpoint.Endpoint
	ListLawsEndpoint         endpoint.Endpoint
	UpdateLawEndpoint        endpoint.Endpoint
	// Upstreams
	CreateUpstreamEndpoint endpoint.Endpoint
	GetUpstreamEndpoint    endpoint.Endpoint
	ListUpstreamsEndpoint  endpoint.Endpoint
	UpdateUpstreamEndpoint endpoint.Endpoint
	// Users
	CreateUserEndpoint   endpoint.Endpoint
	CreateUsersEndpoint  endpoint.Endpoint
	DeleteUserEndpoint   endpoint.Endpoint
	GetUserEndpoint      endpoint.Endpoint
	ListUsersEndpoint    endpoint.Endpoint
	UpdateUserEndpoint   endpoint.Endpoint
	UploadAvatarEndpoint endpoint.Endpoint
}

type Failer interface {
	Failed() error
}
