package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/welaw/welaw/endpoints"
	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/proto"
)

func decodeHTTPCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPCreateUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateUsersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPDeleteUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		return nil, errs.ErrBadRequest
	}
	return endpoints.DeleteUserRequest{Username: username}, nil
}

func decodeHTTPGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPListUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ListUsersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPUpdateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPUploadAvatarRequest(_ context.Context, r *http.Request) (interface{}, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return nil, err
	}
	var req endpoints.UploadAvatarRequest
	var opts proto.UploadAvatarOptions
	var buf bytes.Buffer
	io.Copy(&buf, file)
	opts.Image = buf.Bytes()
	opts.Filename = header.Filename
	username := r.FormValue("username")
	opts.Username = username
	req.Opts = &opts
	return req, err
}
