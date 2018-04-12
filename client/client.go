package client

import (
	"net/url"
	"os"
	"time"

	"golang.org/x/time/rate"

	stdlog "log"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"github.com/welaw/welaw/endpoints"
	"github.com/welaw/welaw/pkg/oauth"
	"github.com/welaw/welaw/services"
)

func NewClient(welawUrl string) services.Service {
	u, err := url.Parse(welawUrl)
	if err != nil {
		stdlog.Fatal(err)
	}

	tracer := stdopentracing.GlobalTracer()
	logger := log.NewLogfmtLogger(os.Stderr)

	return NewClientWithOptions(u, tracer, logger, []httptransport.ClientOption{})
}

func NewClientBasicAuth(welawUrl, username, password string) services.Service {
	u, err := url.Parse(welawUrl)
	if err != nil {
		stdlog.Fatal(err)
	}

	tracer := stdopentracing.GlobalTracer()
	logger := log.NewLogfmtLogger(os.Stderr)

	return NewClientWithOptions(u, tracer, logger, []httptransport.ClientOption{
		httptransport.ClientBefore(oauth.BasicAuthToHTTPRequest(username, password)),
	})
}

func NewClientWithOptions(addr *url.URL, tracer stdopentracing.Tracer, logger log.Logger, options []httptransport.ClientOption) services.Service {
	throttler := ratelimit.NewDelayingLimiter(rate.NewLimiter(rate.Every(time.Second), 5))
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// Ballot

	var createVoteEndpoint endpoint.Endpoint
	{
		createVoteEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/vote/create"),
			encodeHTTPRequest,
			decodeHTTPCreateVoteResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		createVoteEndpoint = opentracing.TraceClient(tracer, "create_vote")(createVoteEndpoint)
		createVoteEndpoint = throttler(createVoteEndpoint)
		createVoteEndpoint = limiter(createVoteEndpoint)
		createVoteEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "create_vote",
			Timeout: 30 * time.Second,
		}))(createVoteEndpoint)
	}

	var createVotesEndpoint endpoint.Endpoint
	{
		createVotesEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/vote/create-batch"),
			encodeHTTPRequest,
			decodeHTTPCreateVotesResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		createVotesEndpoint = opentracing.TraceClient(tracer, "create_votes")(createVotesEndpoint)
		createVotesEndpoint = throttler(createVotesEndpoint)
		createVotesEndpoint = limiter(createVotesEndpoint)
		createVotesEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "create_votes",
			Timeout: 30 * time.Second,
		}))(createVotesEndpoint)
	}

	var getVoteEndpoint endpoint.Endpoint
	{
		getVoteEndpoint = httptransport.NewClient(
			"GET",
			copyUrl(addr, "/api/v1/vote"),
			encodeHTTPRequest,
			decodeHTTPGetVoteResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		getVoteEndpoint = opentracing.TraceClient(tracer, "get_vote")(getVoteEndpoint)
		getVoteEndpoint = throttler(getVoteEndpoint)
		getVoteEndpoint = limiter(getVoteEndpoint)
		getVoteEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "get_vote",
			Timeout: 30 * time.Second,
		}))(getVoteEndpoint)
	}

	var deleteVoteEndpoint endpoint.Endpoint
	{
		deleteVoteEndpoint = httptransport.NewClient(
			"GET",
			copyUrl(addr, "/api/v1/vote/delete"),
			encodeHTTPRequest,
			decodeHTTPDeleteVoteResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		deleteVoteEndpoint = opentracing.TraceClient(tracer, "delete_vote")(deleteVoteEndpoint)
		deleteVoteEndpoint = throttler(deleteVoteEndpoint)
		deleteVoteEndpoint = limiter(deleteVoteEndpoint)
		deleteVoteEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "delete_vote",
			Timeout: 30 * time.Second,
		}))(deleteVoteEndpoint)
	}

	var listVotesEndpoint endpoint.Endpoint
	{
		listVotesEndpoint = httptransport.NewClient(
			"GET",
			copyUrl(addr, "/api/v1/vote/list"),
			encodeHTTPRequest,
			decodeHTTPListVotesResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		listVotesEndpoint = opentracing.TraceClient(tracer, "list_votes")(listVotesEndpoint)
		listVotesEndpoint = throttler(listVotesEndpoint)
		listVotesEndpoint = limiter(listVotesEndpoint)
		listVotesEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "list_votes",
			Timeout: 30 * time.Second,
		}))(listVotesEndpoint)
	}

	var updateVoteEndpoint endpoint.Endpoint
	{
		updateVoteEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/vote/update"),
			encodeHTTPRequest,
			decodeHTTPUpdateVoteResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		updateVoteEndpoint = opentracing.TraceClient(tracer, "update_vote")(updateVoteEndpoint)
		updateVoteEndpoint = throttler(updateVoteEndpoint)
		updateVoteEndpoint = limiter(updateVoteEndpoint)
		updateVoteEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "update_vote",
			Timeout: 30 * time.Second,
		}))(updateVoteEndpoint)
	}

	// Laws

	var createLawEndpoint endpoint.Endpoint
	{
		createLawEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/law/create"),
			encodeHTTPRequest,
			decodeHTTPCreateLawResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		createLawEndpoint = opentracing.TraceClient(tracer, "create_law")(createLawEndpoint)
		createLawEndpoint = throttler(createLawEndpoint)
		createLawEndpoint = limiter(createLawEndpoint)
		createLawEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "create_law",
			Timeout: 30 * time.Second,
		}))(createLawEndpoint)
	}

	var createLawsEndpoint endpoint.Endpoint
	{
		createLawsEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/law/create-batch"),
			encodeHTTPRequest,
			decodeHTTPCreateLawsResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		createLawsEndpoint = opentracing.TraceClient(tracer, "create_law")(createLawsEndpoint)
		createLawsEndpoint = throttler(createLawsEndpoint)
		createLawsEndpoint = limiter(createLawsEndpoint)
		createLawsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "create_laws",
			Timeout: 30 * time.Second,
		}))(createLawsEndpoint)
	}

	var getLawEndpoint endpoint.Endpoint
	{
		getLawEndpoint = httptransport.NewClient(
			"GET",
			copyUrl(addr, "/api/v1/law"),
			encodeHTTPRequest,
			decodeHTTPGetLawResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		getLawEndpoint = opentracing.TraceClient(tracer, "get_law")(getLawEndpoint)
		getLawEndpoint = throttler(getLawEndpoint)
		getLawEndpoint = limiter(getLawEndpoint)
		getLawEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "get_law",
			Timeout: 30 * time.Second,
		}))(getLawEndpoint)
	}

	// Upstream

	var createUpstreamEndpoint endpoint.Endpoint
	{
		createUpstreamEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/upstream/create"),
			encodeHTTPRequest,
			decodeHTTPCreateUpstreamResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		createUpstreamEndpoint = opentracing.TraceClient(tracer, "create_upstream")(createUpstreamEndpoint)
		createUpstreamEndpoint = throttler(createUpstreamEndpoint)
		createUpstreamEndpoint = limiter(createUpstreamEndpoint)
		createUpstreamEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "create_upstream",
			Timeout: 30 * time.Second,
		}))(createUpstreamEndpoint)
	}

	var getUpstreamEndpoint endpoint.Endpoint
	{
		getUpstreamEndpoint = httptransport.NewClient(
			"GET",
			copyUrl(addr, "/api/v1/upstream"),
			encodeHTTPRequest,
			decodeHTTPGetUpstreamResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		getUpstreamEndpoint = opentracing.TraceClient(tracer, "get_upstream")(getUpstreamEndpoint)
		getUpstreamEndpoint = throttler(getUpstreamEndpoint)
		getUpstreamEndpoint = limiter(getUpstreamEndpoint)
		getUpstreamEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "get_upstream",
			Timeout: 30 * time.Second,
		}))(getUpstreamEndpoint)
	}

	var listUpstreamsEndpoint endpoint.Endpoint
	{
		listUpstreamsEndpoint = httptransport.NewClient(
			"GET",
			copyUrl(addr, "/api/v1/upstream/list"),
			encodeHTTPRequest,
			decodeHTTPListUpstreamsResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		listUpstreamsEndpoint = opentracing.TraceClient(tracer, "list_upstreams")(listUpstreamsEndpoint)
		listUpstreamsEndpoint = throttler(listUpstreamsEndpoint)
		listUpstreamsEndpoint = limiter(listUpstreamsEndpoint)
		listUpstreamsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "list_upstreams",
			Timeout: 30 * time.Second,
		}))(listUpstreamsEndpoint)
	}

	// Users

	var createUserEndpoint endpoint.Endpoint
	{
		createUserEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/user/create"),
			encodeHTTPRequest,
			decodeHTTPCreateUserResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		createUserEndpoint = opentracing.TraceClient(tracer, "create_user")(createUserEndpoint)
		createUserEndpoint = throttler(createUserEndpoint)
		createUserEndpoint = limiter(createUserEndpoint)
		createUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "create_user",
			Timeout: 30 * time.Second,
		}))(createUserEndpoint)
	}

	var createUsersEndpoint endpoint.Endpoint
	{
		createUsersEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/user/create-batch"),
			encodeHTTPRequest,
			decodeHTTPCreateUsersResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		createUsersEndpoint = opentracing.TraceClient(tracer, "create_users")(createUsersEndpoint)
		createUsersEndpoint = throttler(createUsersEndpoint)
		createUsersEndpoint = limiter(createUsersEndpoint)
		createUsersEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "create_users",
			Timeout: 30 * time.Second,
		}))(createUsersEndpoint)
	}

	var getUserEndpoint endpoint.Endpoint
	{
		getUserEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/user"),
			encodeHTTPRequest,
			decodeHTTPGetUserResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		getUserEndpoint = opentracing.TraceClient(tracer, "get_user")(getUserEndpoint)
		getUserEndpoint = throttler(getUserEndpoint)
		getUserEndpoint = limiter(getUserEndpoint)
		getUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "get_user_by_username",
			Timeout: 30 * time.Second,
		}))(getUserEndpoint)
	}

	var uploadAvatarEndpoint endpoint.Endpoint
	{
		uploadAvatarEndpoint = httptransport.NewClient(
			"POST",
			copyUrl(addr, "/api/v1/user/update/avatar"),
			encodeHTTPUploadAvatarRequest,
			decodeHTTPUploadAvatarResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(tracer, logger)))...,
		).Endpoint()
		uploadAvatarEndpoint = opentracing.TraceClient(tracer, "upload_avatar")(uploadAvatarEndpoint)
		uploadAvatarEndpoint = throttler(uploadAvatarEndpoint)
		uploadAvatarEndpoint = limiter(uploadAvatarEndpoint)
		uploadAvatarEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "upload_avatar",
			Timeout: 30 * time.Second,
		}))(uploadAvatarEndpoint)
	}

	return endpoints.Endpoints{
		// Ballot
		CreateVoteEndpoint:  createVoteEndpoint,
		CreateVotesEndpoint: createVotesEndpoint,
		GetVoteEndpoint:     getVoteEndpoint,
		DeleteVoteEndpoint:  deleteVoteEndpoint,
		ListVotesEndpoint:   listVotesEndpoint,
		UpdateVoteEndpoint:  updateVoteEndpoint,
		// Laws
		CreateLawEndpoint:  createLawEndpoint,
		CreateLawsEndpoint: createLawsEndpoint,
		//GetLawEndpoint:    getLawEndpoint,

		// Upstreams
		CreateUpstreamEndpoint: createUpstreamEndpoint,
		GetUpstreamEndpoint:    getUpstreamEndpoint,
		ListUpstreamsEndpoint:  listUpstreamsEndpoint,
		//UpdateUpstreamEndpoint: updateUpstreamEndpoint,

		// Users
		CreateUserEndpoint:   createUserEndpoint,
		CreateUsersEndpoint:  createUsersEndpoint,
		GetUserEndpoint:      getUserEndpoint,
		UploadAvatarEndpoint: uploadAvatarEndpoint,
	}
}
