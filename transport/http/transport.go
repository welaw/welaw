package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	stdhttp "net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/welaw/welaw/endpoints"
	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/pkg/oauth"

	stdjwt "github.com/dgrijalva/jwt-go"
	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type key int

const (
	BaseURL           = "/api/v1/"
	contextReqKey key = iota
)

func MakeHTTPHandler(set endpoints.Endpoints, tracer stdopentracing.Tracer, logger log.Logger, keyfunc stdjwt.Keyfunc, sc *securecookie.SecureCookie) stdhttp.Handler {
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errs.ErrorEncoder),
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerBefore(oauth.BasicAuthToHTTPContext()),
		httptransport.ServerBefore(oauth.CookieToHTTPContext(logger, keyfunc, sc)),
	}

	r := mux.NewRouter()

	r.Handle(BaseURL+"admin/stats", httptransport.NewServer(
		set.GetServerStatsEndpoint,
		decodeHTTPGetServerStatsRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "get_server_stats", logger)))...,
	)).Methods("GET")

	r.Handle(BaseURL+"admin/repos", httptransport.NewServer(
		set.SaveReposEndpoint,
		decodeHTTPSaveReposRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "repos", logger)))...,
	)).Methods("POST")

	// Auth

	r.Handle(BaseURL+"auth/login", httptransport.NewServer(
		set.LoginEndpoint,
		decodeHTTPLoginRequest,
		encodeHTTPLoginResponse,
		append(
			opts,
			httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "login", logger)),
			httptransport.ServerBefore(storeHTTPRequest))...,
	)).Methods("GET")

	r.Handle(BaseURL+"auth/login-callback", httptransport.NewServer(
		set.LoginCallbackEndpoint,
		decodeHTTPLoginCallbackRequest,
		encodeHTTPLoginCallbackResponse,
		append(
			opts,
			httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "login-callback", logger)),
			httptransport.ServerBefore(oauth.LoginCookieToHTTPContext(logger, keyfunc, sc)),
			httptransport.ServerBefore(storeHTTPRequest))...,
	)).Methods("GET")

	r.Handle(BaseURL+"auth/logged-in-check", httptransport.NewServer(
		set.LoggedInCheckEndpoint,
		decodeHTTPLoggedInCheckRequest,
		encodeHTTPResponse,
		append(
			opts,
			httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "logged_in_check", logger)),
			httptransport.ServerBefore(storeHTTPRequest))...,
	)).Methods("GET")

	r.Handle(BaseURL+"auth/login-as", httptransport.NewServer(
		set.LoginAsEndpoint,
		decodeHTTPLoginAsRequest,
		encodeHTTPResponse,
		append(
			opts,
			httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "login_as", logger)),
			httptransport.ServerBefore(storeHTTPRequest))...,
	)).Methods("GET")

	r.Handle(BaseURL+"auth/logout", httptransport.NewServer(
		set.LogoutEndpoint,
		decodeHTTPLogoutRequest,
		encodeHTTPLogoutResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "logout", logger)))...,
	)).Methods("GET")

	/*
		Ballot
	*/

	r.Handle(BaseURL+"vote/create", httptransport.NewServer(
		set.CreateVoteEndpoint,
		decodeHTTPCreateVoteRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_vote", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"vote/create-batch", httptransport.NewServer(
		set.CreateVotesEndpoint,
		decodeHTTPCreateVotesRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_votes", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"vote/update", httptransport.NewServer(
		set.UpdateVoteEndpoint,
		decodeHTTPUpdateVoteRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "update_vote", logger)))...,
	)).Methods("POST")

	//r.Handle(BaseURL+"vote/{law_id}/{user_id}", httptransport.NewServer(
	r.Handle(BaseURL+"vote", httptransport.NewServer(
		set.GetVoteEndpoint,
		decodeHTTPGetVoteRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "get_vote", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"vote/delete", httptransport.NewServer(
		set.DeleteVoteEndpoint,
		decodeHTTPDeleteVoteRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "delete_vote", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"vote/list", httptransport.NewServer(
		set.ListVotesEndpoint,
		decodeHTTPListVotesRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "list_votes", logger)))...,
	)).Methods("POST")

	/*
		Law
	*/

	r.Handle(BaseURL+"law/comment/annotation/create", httptransport.NewServer(
		set.CreateAnnotationEndpoint,
		decodeHTTPCreateAnnotationRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_annotation", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/comment/annotation/delete/{id}", httptransport.NewServer(
		set.DeleteAnnotationEndpoint,
		decodeHTTPDeleteAnnotationRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "delete_annotation", logger)))...,
	)).Methods("DELETE")

	r.Handle(BaseURL+"law/comment/annotation/list", httptransport.NewServer(
		set.ListAnnotationsEndpoint,
		decodeHTTPListAnnotationsRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "list_annotations", logger)))...,
	)).Methods("GET")

	r.Handle(BaseURL+"law/comment/create", httptransport.NewServer(
		set.CreateCommentEndpoint,
		decodeHTTPCreateCommentRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_comment", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/comment/delete", httptransport.NewServer(
		set.DeleteCommentEndpoint,
		decodeHTTPDeleteCommentRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "delete_comment", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/comment", httptransport.NewServer(
		set.GetCommentEndpoint,
		decodeHTTPGetCommentRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "get_comment", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/comment/update", httptransport.NewServer(
		set.UpdateCommentEndpoint,
		decodeHTTPUpdateCommentRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "update_comment", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/comment/like", httptransport.NewServer(
		set.LikeCommentEndpoint,
		decodeHTTPLikeCommentRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "like_comment", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/comment/list", httptransport.NewServer(
		set.ListCommentsEndpoint,
		decodeHTTPListCommentsRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "list_comments", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/create", httptransport.NewServer(
		set.CreateLawEndpoint,
		decodeHTTPCreateLawRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_law", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/create-batch", httptransport.NewServer(
		set.CreateLawsEndpoint,
		decodeHTTPCreateLawsRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_laws", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/list", httptransport.NewServer(
		set.ListLawsEndpoint,
		decodeHTTPListLawsRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "list_laws", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/delete/{upstream}/{ident}", httptransport.NewServer(
		set.DeleteLawEndpoint,
		decodeHTTPDeleteLawRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "delete_law", logger)))...,
	)).Methods("GET")

	r.Handle(BaseURL+"law/diff", httptransport.NewServer(
		set.DiffLawsEndpoint,
		decodeHTTPDiffLawsRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "diff_laws", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law", httptransport.NewServer(
		set.GetLawEndpoint,
		decodeHTTPGetLawRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "get_law", logger)))...,
	)).Methods("POST")

	r.Handle(BaseURL+"law/update", httptransport.NewServer(
		set.UpdateLawEndpoint,
		decodeHTTPUpdateLawRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "update_law", logger)))...,
	)).Methods("POST")

	// Upstreams

	r.Handle(BaseURL+"upstream/list", httptransport.NewServer(
		set.ListUpstreamsEndpoint,
		decodeHTTPListUpstreamsRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "list_upstreams", logger)))...,
	)).Methods("GET")

	// update upstream
	r.Handle(BaseURL+"upstream/update", httptransport.NewServer(
		set.UpdateUpstreamEndpoint,
		decodeHTTPUpdateUpstreamRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "update_upstream", logger)))...,
	)).Methods("POST")

	// get upstream
	r.Handle(BaseURL+"upstream/{ident}", httptransport.NewServer(
		set.GetUpstreamEndpoint,
		decodeHTTPGetUpstreamRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "get_upstream", logger)))...,
	)).Methods("GET")

	// Users

	// get user
	r.Handle(BaseURL+"user", httptransport.NewServer(
		set.GetUserEndpoint,
		decodeHTTPGetUserRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "get_user", logger)))...,
	)).Methods("POST")

	// create user
	r.Handle(BaseURL+"user/create", httptransport.NewServer(
		set.CreateUserEndpoint,
		decodeHTTPCreateUserRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_user", logger)))...,
	)).Methods("POST")

	// create users
	r.Handle(BaseURL+"user/create-batch", httptransport.NewServer(
		set.CreateUsersEndpoint,
		decodeHTTPCreateUsersRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "create_users", logger)))...,
	)).Methods("POST")

	// delete user
	r.Handle(BaseURL+"user/{username}/delete", httptransport.NewServer(
		set.DeleteUserEndpoint,
		decodeHTTPDeleteUserRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "delete_user", logger)))...,
	)).Methods("GET")

	// list users
	r.Handle(BaseURL+"user/list", httptransport.NewServer(
		set.ListUsersEndpoint,
		decodeHTTPListUsersRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "list_users", logger)))...,
	)).Methods("POST")

	// update user
	r.Handle(BaseURL+"user/update", httptransport.NewServer(
		set.UpdateUserEndpoint,
		decodeHTTPUpdateUserRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "update_user", logger)))...,
	)).Methods("POST")

	// update user
	r.Handle(BaseURL+"user/update/avatar", httptransport.NewServer(
		set.UploadAvatarEndpoint,
		decodeHTTPUploadAvatarRequest,
		encodeHTTPResponse,
		append(opts, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "upload_avatar", logger)))...,
	)).Methods("POST")

	// metrics
	r.Handle(BaseURL+"metrics", stdprometheus.Handler())

	return r
}

func encodeHTTPRequest(_ context.Context, r *stdhttp.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	//r.Header.Set("Content-Type", "application/json")
	//r.ContentLength = int64(len(buf.Bytes()))
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeHTTPResponse(ctx context.Context, w stdhttp.ResponseWriter, response interface{}) error {
	if failer, ok := response.(endpoints.Failer); ok && failer.Failed() != nil {
		errs.ErrorEncoder(ctx, failer.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func storeHTTPRequest(ctx context.Context, req *stdhttp.Request) context.Context {
	ctx, _ = context.WithTimeout(ctx, time.Duration(30*time.Second))
	return context.WithValue(ctx, contextReqKey, req)
}
