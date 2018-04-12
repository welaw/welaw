package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/welaw/welaw/endpoints"
	"github.com/welaw/welaw/pkg/errs"
)

func decodeHTTPCreateUserResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.CreateUserResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPCreateUsersResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.CreateUsersResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPGetUserResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.GetUserResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPUploadAvatarResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.UploadAvatarResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func encodeHTTPUploadAvatarRequest(_ context.Context, req *http.Request, request interface{}) error {
	r := request.(endpoints.UploadAvatarRequest)
	if r.Opts == nil {
		return errs.BadRequest("opts not found")
	}

	filename := r.Opts.Filename
	image := r.Opts.Image
	username := r.Opts.Username

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	defer w.Close()
	fw, err := w.CreateFormFile("image", filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fw, bytes.NewReader(image)); err != nil {
		return err
	}
	if err = w.WriteField("username", username); err != nil {
		return err
	}
	w.Close()
	req.Body = ioutil.NopCloser(&buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return nil
}
