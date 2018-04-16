package endpoints

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/proto"
	"github.com/welaw/welaw/services"
)

func (e Endpoints) LoggedInCheck(ctx context.Context) (*proto.User, error) {
	resp, err := e.LoggedInCheckEndpoint(ctx, LoggedInCheckRequest{})
	if err != nil {
		return nil, err
	}
	r := resp.(LoggedInCheckResponse)
	return r.User, r.Err
}

func (e Endpoints) Login(ctx context.Context, returnUrl, provider string) (*http.Cookie, string, error) {
	resp, err := e.LoginEndpoint(ctx, LoginRequest{ReturnUrl: returnUrl, Provider: provider})
	if err != nil {
		return nil, "", err
	}
	r := resp.(LoginResponse)
	return r.Cookie, r.Url, r.Err
}

func (e Endpoints) LoginCallback(ctx context.Context, state, code string) (*http.Cookie, *http.Cookie, string, error) {
	resp, err := e.LoginCallbackEndpoint(ctx, LoginCallbackRequest{
		State: state,
		Code:  code,
	})
	if err != nil {
		return nil, nil, "", err
	}
	r := resp.(LoginCallbackResponse)
	return r.Cookie, r.LoginCookie, r.Url, r.Err
}

func (e Endpoints) LoginAs(ctx context.Context, user *proto.User) error {
	resp, err := e.LoginAsEndpoint(ctx, LoginAsRequest{User: user})
	if err != nil {
		return err
	}
	r := resp.(LoginAsResponse)
	return r.Err
}

func (e Endpoints) Logout(ctx context.Context) error {
	resp, err := e.LogoutEndpoint(ctx, LogoutRequest{})
	if err != nil {
		return err
	}
	r := resp.(LogoutResponse)
	return r.Err
}

func MakeLoggedInCheckEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(LoggedInCheckRequest)
		user, err := svc.LoggedInCheck(ctx)
		if err == errs.ErrUnauthorized {
			return LoggedInCheckResponse{}, nil
		}
		return LoggedInCheckResponse{User: user, Err: err}, nil
	}
}

func MakeLoginEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		cookie, url, err := svc.Login(ctx, req.ReturnUrl, req.Provider)
		return LoginResponse{Cookie: cookie, Url: url, Err: err}, nil
	}
}

func MakeLoginCallbackEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginCallbackRequest)
		c, lc, url, err := svc.LoginCallback(ctx, req.State, req.Code)
		return LoginCallbackResponse{Cookie: c, LoginCookie: lc, Url: url, Err: err}, nil
	}
}

func MakeLoginAsEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginAsRequest)
		err := svc.LoginAs(ctx, req.User)
		return LoginAsResponse{Err: err}, nil
	}
}

func MakeLogoutEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(LogoutRequest)
		err := svc.Logout(ctx)
		return LogoutResponse{Err: err}, nil
	}
}

type LoggedInCheckRequest struct{}

type LoggedInCheckResponse struct {
	User *proto.User `json:"user"`
	Err  error       `json:"-"`
}

func (r LoggedInCheckResponse) Failed() error { return r.Err }

type LoginRequest struct {
	ReturnUrl string `json:"return_url"`
	Provider  string `json:"provider"`
}

type LoginResponse struct {
	Cookie *http.Cookie `json:"-"`
	Url    string       `json:"url"`
	Err    error        `json:"-"`
}

func (r LoginResponse) Failed() error { return r.Err }

type LoginCallbackRequest struct {
	State string `json:"state"`
	Code  string `json:"code"`
}

type LoginCallbackResponse struct {
	LoginCookie *http.Cookie `json:"-"`
	Cookie      *http.Cookie `json:"-"`
	Url         string       `json:"url"`
	Err         error        `json:"-"`
}

func (r LoginCallbackResponse) Failed() error { return r.Err }

type LoginAsRequest struct {
	User *proto.User `json:"user"`
}

type LoginAsResponse struct {
	Err error `json:"-"`
}

func (r LoginAsResponse) Failed() error { return r.Err }

type LogoutRequest struct {
}

type LogoutResponse struct {
	Err error `json:"-"`
}

func (r LogoutResponse) Failed() error { return r.Err }
