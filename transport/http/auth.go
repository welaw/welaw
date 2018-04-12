package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/welaw/welaw/endpoints"
	"github.com/welaw/welaw/pkg/errs"
)

type GoogleResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func decodeHTTPLoginRequest(_ context.Context, req *http.Request) (interface{}, error) {
	vars := req.URL.Query()
	var url, provider string
	urls, ok := vars["url"]
	if ok && len(urls) > 0 {
		url = urls[0]
	}
	providers, ok := vars["provider"]
	if ok && len(providers) > 0 {
		provider = providers[0]
	}
	return endpoints.LoginRequest{ReturnUrl: url, Provider: provider}, nil
}

func encodeHTTPLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	req := ctx.Value(contextReqKey).(*http.Request)
	r := response.(endpoints.LoginResponse)
	http.SetCookie(w, r.Cookie)
	http.Redirect(w, req, r.Url, http.StatusTemporaryRedirect)
	return nil
}

func decodeHTTPLoginCallbackRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	state := req.FormValue("state")
	code := req.FormValue("code")
	return endpoints.LoginCallbackRequest{State: state, Code: code}, nil
}

func encodeHTTPLoginCallbackResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	req := ctx.Value(contextReqKey).(*http.Request)
	r := response.(endpoints.LoginCallbackResponse)
	if r.LoginCookie != nil {
		http.SetCookie(w, r.LoginCookie)
	}
	if r.Cookie != nil {
		http.SetCookie(w, r.Cookie)
	}
	http.Redirect(w, req, r.Url, http.StatusTemporaryRedirect)
	return nil
}

func decodeHTTPLoginAsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.LoginAsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPLoggedInCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoints.LoggedInCheckRequest{}, nil
}

func encodeHTTPLoginAsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(endpoints.LoginAsResponse)
	if r.Err != nil {
		errs.ErrorEncoder(ctx, r.Err, w)
		return nil
	}
	return json.NewEncoder(w).Encode(r)
}

func decodeHTTPLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoints.LogoutRequest{}, nil
}

func encodeHTTPLogoutResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(endpoints.LogoutResponse)
	if r.Err != nil {
		errs.ErrorEncoder(ctx, r.Err, w)
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "welaw_auth",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		MaxAge:   0,
		Path:     "/",
	})

	return json.NewEncoder(w).Encode(r)
}
