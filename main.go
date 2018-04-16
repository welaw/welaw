package main

import (
	"flag"
	"fmt"
	stdhttp "net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"github.com/welaw/welaw/backend/database"
	"github.com/welaw/welaw/backend/filesystem"
	"github.com/welaw/welaw/backend/vcs"
	"github.com/welaw/welaw/endpoints"
	"github.com/welaw/welaw/instrumentation"
	"github.com/welaw/welaw/logging"
	"github.com/welaw/welaw/pkg/oauth"
	"github.com/welaw/welaw/services"
	"github.com/welaw/welaw/transport/http"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "Welaw service (http) listen address")
	)
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("warning: godotenv: .env not found")
	}

	// logger
	var (
		logger log.Logger
	)
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// transport
	var (
		fieldKeys = []string{"method", "error"}
		tracer    = stdopentracing.GlobalTracer()
		mux       = stdhttp.NewServeMux()
	)

	// filesystem
	var (
		repos = os.Getenv("REPOS_DIR")
	)
	// errors
	var (
		errc = make(chan error, 1)
	)

	// jwt
	var (
		signingKey = securecookie.GenerateRandomKey(16)
		keyfunc    = func(token *stdjwt.Token) (interface{}, error) {
			return signingKey, nil
		}
	)

	// oauth
	var (
		amznClientId = os.Getenv("AMZN_CLIENT_ID")
		amznSecret   = os.Getenv("AMZN_SECRET")
		googClientId = os.Getenv("GOOG_CLIENT_ID")
		googSecret   = os.Getenv("GOOG_SECRET")
		msftClientId = os.Getenv("MSFT_CLIENT_ID")
		msftSecret   = os.Getenv("MSFT_SECRET")
	)

	// securecookie
	var (
		hashKey  = securecookie.GenerateRandomKey(16)
		blockKey = securecookie.GenerateRandomKey(16)
		sc       = securecookie.New(hashKey, blockKey)
	)

	// database
	var (
		connStr = os.Getenv("POSTGRES_CONNECTION")
	)

	var (
		useSecureCookies    bool
		useSecureCookiesEnv = os.Getenv("USE_SECURE_COOKIES")
	)
	{
		useSecureCookies, err = strconv.ParseBool(useSecureCookiesEnv)
		if err != nil {
			fmt.Println("warning: error parsing USE_SECURE_COOKIES environment variable")
			useSecureCookies = false
		}
	}

	var (
		originsEnv       = os.Getenv("ALLOWED_ORIGINS")
		loginRedirectURL = os.Getenv("LOGIN_REDIRECT_URL")
		//avatarURL        = os.Getenv("AVATAR_URL")
		opts = &services.ServerConfigOptions{
			LoginFailedURL:   os.Getenv("LOGIN_FAILED_URL"),
			LoginSuccessURL:  os.Getenv("LOGIN_SUCCESS_URL"),
			UseSecureCookies: useSecureCookies,
			SigningKey:       signingKey,
			AvatarDir:        os.Getenv("STATIC_DIR"),
		}
	)

	if os.Getenv("HTTP_ADDR") != "" {
		*httpAddr = os.Getenv("HTTP_ADDR")
	}

	var service services.Service
	{
		providers := map[string]oauth.Provider{
			oauth.ProviderAmazon:    oauth.NewProviderAmazon(loginRedirectURL, amznClientId, amznSecret),
			oauth.ProviderGoogle:    oauth.NewProviderGoogle(loginRedirectURL, googClientId, googSecret),
			oauth.ProviderMicrosoft: oauth.NewProviderMicrosoft(loginRedirectURL, msftClientId, msftSecret),
		}
		db := database.NewDatabase(connStr, logger, &database.DatabaseConfigOptions{
			//AvatarURL: fmt.Sprintf("%s/%s", os.Getenv("STATIC_DIR"), opts.AvatarDir),
			AvatarURL: os.Getenv("AVATAR_URL"),
		})
		fs := filesystem.NewS3FS(
			os.Getenv("AWS_REGION"),
			os.Getenv("AWS_BUCKET"),
		)
		vcs, err := vcs.NewVcs(repos, logger)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		service = services.NewService(
			db,
			fs,
			logger,
			vcs,
			providers,
			sc,
			opts,
		)
		service = logging.NewLoggingMiddleware(logger, service)
		service = instrumentation.NewInstrumentatingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "welaw",
				Subsystem: "law_service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "welaw",
				Subsystem: "law_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			service,
		)

		httpEndpoints := endpoints.Endpoints{
			// Admin
			GetServerStatsEndpoint: endpoints.MakeGetServerStatsEndpoint(service),
			SaveReposEndpoint:      endpoints.MakeSaveReposEndpoint(service),
			// Auth
			LoggedInCheckEndpoint: endpoints.MakeLoggedInCheckEndpoint(service),
			LoginEndpoint:         endpoints.MakeLoginEndpoint(service),
			LoginCallbackEndpoint: endpoints.MakeLoginCallbackEndpoint(service),
			LogoutEndpoint:        endpoints.MakeLogoutEndpoint(service),
			// Ballot
			CreateVoteEndpoint:  endpoints.MakeCreateVoteEndpoint(service),
			CreateVotesEndpoint: endpoints.MakeCreateVotesEndpoint(service),
			DeleteVoteEndpoint:  endpoints.MakeDeleteVoteEndpoint(service),
			GetVoteEndpoint:     endpoints.MakeGetVoteEndpoint(service),
			ListVotesEndpoint:   endpoints.MakeListVotesEndpoint(service),
			UpdateVoteEndpoint:  endpoints.MakeUpdateVoteEndpoint(service),
			// Law
			CreateAnnotationEndpoint: endpoints.MakeCreateAnnotationEndpoint(service),
			DeleteAnnotationEndpoint: endpoints.MakeDeleteAnnotationEndpoint(service),
			ListAnnotationsEndpoint:  endpoints.MakeListAnnotationsEndpoint(service),
			ListCommentsEndpoint:     endpoints.MakeListCommentsEndpoint(service),
			CreateCommentEndpoint:    endpoints.MakeCreateCommentEndpoint(service),
			DeleteCommentEndpoint:    endpoints.MakeDeleteCommentEndpoint(service),
			GetCommentEndpoint:       endpoints.MakeGetCommentEndpoint(service),
			UpdateCommentEndpoint:    endpoints.MakeUpdateCommentEndpoint(service),
			LikeCommentEndpoint:      endpoints.MakeLikeCommentEndpoint(service),
			CreateLawEndpoint:        endpoints.MakeCreateLawEndpoint(service),
			CreateLawsEndpoint:       endpoints.MakeCreateLawsEndpoint(service),
			DiffLawsEndpoint:         endpoints.MakeDiffLawsEndpoint(service),
			GetLawEndpoint:           endpoints.MakeGetLawEndpoint(service),
			ListLawsEndpoint:         endpoints.MakeListLawsEndpoint(service),
			// Upstream
			GetUpstreamEndpoint:    endpoints.MakeGetUpstreamEndpoint(service),
			ListUpstreamsEndpoint:  endpoints.MakeListUpstreamsEndpoint(service),
			UpdateUpstreamEndpoint: endpoints.MakeUpdateUpstreamEndpoint(service),
			// User
			CreateUserEndpoint:   endpoints.MakeCreateUserEndpoint(service),
			CreateUsersEndpoint:  endpoints.MakeCreateUsersEndpoint(service),
			DeleteUserEndpoint:   endpoints.MakeDeleteUserEndpoint(service),
			GetUserEndpoint:      endpoints.MakeGetUserEndpoint(service),
			ListUsersEndpoint:    endpoints.MakeListUsersEndpoint(service),
			UpdateUserEndpoint:   endpoints.MakeUpdateUserEndpoint(service),
			UploadAvatarEndpoint: endpoints.MakeUploadAvatarEndpoint(service),
		}

		mux.Handle(http.BaseURL, http.MakeHTTPHandler(httpEndpoints, tracer, logger, keyfunc, sc))

		// Serve HTTP
		go func() {
			logger.Log("msg", "HTTP", "addr", httpAddr)

			// TODO
			origins := strings.Split(originsEnv, ",")
			crs := cors.New(cors.Options{
				AllowedOrigins:   origins,
				AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
				AllowedHeaders:   []string{"*"},
				AllowCredentials: true,
			})

			logger.Log("err", stdhttp.ListenAndServe(*httpAddr, crs.Handler(mux)))
		}()

		go func() {
			c := make(chan os.Signal)
			signal.Notify(c, syscall.SIGINT)
			errc <- fmt.Errorf("%s", <-c)
		}()
		logger.Log("terminated", <-errc)
		db.Close()
	}

}
