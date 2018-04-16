package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/welaw/welaw/proto"
	"github.com/welaw/welaw/services"
)

func (e Endpoints) CreateUser(ctx context.Context, newUser *proto.User) (*proto.User, error) {
	req := CreateUserRequest{User: newUser}
	resp, err := e.CreateUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(CreateUserResponse)
	return r.User, r.Err
}

func (e Endpoints) CreateUsers(ctx context.Context, users []*proto.User) ([]*proto.User, error) {
	req := CreateUsersRequest{Users: users}
	resp, err := e.CreateUsersEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(CreateUsersResponse)
	return r.Users, r.Err
}

func (e Endpoints) DeleteUser(ctx context.Context, username string) error {
	request := DeleteUserRequest{Username: username}
	r, err := e.DeleteUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	return r.(DeleteUserResponse).Err
}

func (e Endpoints) GetUser(ctx context.Context, opts *proto.GetUserOptions) (*proto.User, error) {
	req := GetUserRequest{opts}
	resp, err := e.GetUserEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(GetUserResponse)
	return r.User, r.Err
}

func (e Endpoints) ListUsers(ctx context.Context, opts *proto.ListUsersOptions) ([]*proto.User, int, error) {
	req := ListUsersRequest{Opts: opts}
	resp, err := e.ListUsersEndpoint(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	r := resp.(ListUsersResponse)
	return r.Users, r.Total, r.Err
}

func (e Endpoints) UpdateUser(ctx context.Context, username string, opts *proto.UpdateUserOptions) (*proto.User, error) {
	request := UpdateUserRequest{Username: username, Opts: opts}
	response, err := e.UpdateUserEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	r := response.(UpdateUserResponse)
	return r.User, r.Err
}

func (e Endpoints) UploadAvatar(ctx context.Context, opts *proto.UploadAvatarOptions) error {
	request := UploadAvatarRequest{Opts: opts}
	response, err := e.UploadAvatarEndpoint(ctx, request)
	if err != nil {
		return err
	}
	r := response.(UploadAvatarResponse)
	return r.Err
}

func MakeCreateUserEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		user, err := svc.CreateUser(ctx, req.User)
		return CreateUserResponse{User: user, Err: err}, nil
	}
}

func MakeCreateUsersEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUsersRequest)
		users, err := svc.CreateUsers(ctx, req.Users)
		return CreateUsersResponse{Users: users, Err: err}, nil
	}
}

func MakeDeleteUserEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteUserRequest)
		err := svc.DeleteUser(ctx, req.Username)
		return DeleteUserResponse{Err: err}, nil
	}
}

func MakeGetUserEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := svc.GetUser(ctx, req.Opts)
		if err != nil {
			if req.Opts != nil && req.Opts.Quiet {
				return GetUserResponse{}, nil
			}
			return GetUserResponse{Err: err}, nil
		}
		return GetUserResponse{User: user}, nil
	}
}

func MakeListUsersEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListUsersRequest)
		users, total, err := svc.ListUsers(ctx, req.Opts)
		return ListUsersResponse{Users: users, Total: total, Err: err}, nil
	}
}

func MakeUpdateUserEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserRequest)
		u, err := svc.UpdateUser(ctx, req.Username, req.Opts)
		return UpdateUserResponse{User: u, Err: err}, nil
	}
}

func MakeUploadAvatarEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UploadAvatarRequest)
		err := svc.UploadAvatar(ctx, req.Opts)
		return UploadAvatarResponse{Err: err}, nil
	}
}

type CreateUserRequest struct {
	User *proto.User `json:"user"`
}

type CreateUserResponse struct {
	User *proto.User `json:"user"`
	Err  error       `json:"-"`
}

func (r CreateUserResponse) Failed() error { return r.Err }

type CreateUsersRequest struct {
	Users []*proto.User `json:"suser"`
}

type CreateUsersResponse struct {
	Users []*proto.User `json:"users"`
	Err   error         `json:"-"`
}

func (r CreateUsersResponse) Failed() error { return r.Err }

type DeleteUserRequest struct {
	Username string `json:"username"`
}

type DeleteUserResponse struct {
	Err error `json:"-"`
}

func (r DeleteUserResponse) Failed() error { return r.Err }

type GetUserRequest struct {
	Opts *proto.GetUserOptions `json:"opts"`
}

type GetUserResponse struct {
	//GetUserReply *proto.GetUserReply `json:"reply"`
	User *proto.User `json:"user"`
	Err  error       `json:"-"`
}

func (r GetUserResponse) Failed() error { return r.Err }

type ListUsersRequest struct {
	Opts *proto.ListUsersOptions `json:"opts"`
}

type ListUsersResponse struct {
	Users []*proto.User `json:"users"`
	Total int           `json:"total"`
	Err   error         `json:"-"`
}

func (r ListUsersResponse) Failed() error { return r.Err }

type UpdateUserRequest struct {
	//User *proto.User              `json:"user"`
	Username string                   `json:"username"`
	Opts     *proto.UpdateUserOptions `json:"opts"`
}

type UpdateUserResponse struct {
	User *proto.User `json:"user"`
	Err  error       `json:"-"`
}

func (r UpdateUserResponse) Failed() error { return r.Err }

type UploadAvatarRequest struct {
	Opts *proto.UploadAvatarOptions `json:"opts"`
}

type UploadAvatarResponse struct {
	Err error `json:"-"`
}

func (r UploadAvatarResponse) Failed() error { return r.Err }
