package oauth

import (
	"context"
	"net/http"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/securecookie"
)

const (
	cookieName      = "welaw_auth"
	cookieLoginName = "welaw_auth_login"
)

func CookieToHTTPRequest() httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		return ctx
	}
}

func LoginCookieToHTTPContext(logger log.Logger, keyFunc stdjwt.Keyfunc, sc *securecookie.SecureCookie) httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		cookie, err := r.Cookie(cookieLoginName)
		if err != nil {
			logger.Log("method", "secure_login_cookie_to_http_context", "error getting cookie", err)
			return ctx
		}
		var decoded string
		err = sc.Decode(cookieLoginName, cookie.Value, &decoded)
		if err != nil {
			logger.Log("method", "secure_login_cookie_to_http_context", "error decoding cookie", err)
			return ctx
		}
		return context.WithValue(ctx, "state", decoded)
	}
}

func CookieToHTTPContext(logger log.Logger, keyFunc stdjwt.Keyfunc, sc *securecookie.SecureCookie) httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			logger.Log("method", "secure_cookie_to_http_context", "error getting cookie", err)
			return ctx
		}
		var decoded string
		err = sc.Decode(cookieName, cookie.Value, &decoded)
		if err != nil {
			logger.Log("method", "secure_cookie_to_http_context", "error decoding cookie", err.Error())
			return ctx
		}
		token, err := stdjwt.ParseWithClaims(decoded, stdjwt.MapClaims{}, func(token *stdjwt.Token) (interface{}, error) {
			if token.Method != method {
				logger.Log("method", "secure_cookie_to_http_context", "token methods mismatch", err)
				return nil, jwt.ErrUnexpectedSigningMethod
			}
			return keyFunc(token)
		})
		if err != nil {
			logger.Log("method", "secure_cookie_to_http_context", "error parsing claims", err)
			return ctx
		}
		c := token.Claims.(stdjwt.MapClaims)
		su, ok := c["sub"].(string)
		if !ok {
			logger.Log("method", "secure_cookie_to_http_context", "sub not found", err)
			return ctx
		}
		return context.WithValue(ctx, "user_id", su)
	}
}

func AuthGuard(next http.Handler) http.HandlerFunc {
	return BasicAuthGuard(JwtAuthGuard(next))
}

func BasicAuthGuard(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, _, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		// TODO
		if username == "" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func JwtAuthGuard(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			http.Error(w, "authorization failed with error", http.StatusUnauthorized)
			return
		}
		if cookie == nil {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func BasicAuthToHTTPRequest(username, password string) httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		r.SetBasicAuth(username, password)
		return ctx
	}
}

func BasicAuthToHTTPContext() httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		username, password, ok := r.BasicAuth()
		if !ok {
			return ctx
		}

		ctx = context.WithValue(ctx, "password", password)
		return context.WithValue(ctx, "username", username)
	}
}
