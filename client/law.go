package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/welaw/welaw/endpoints"
	"github.com/welaw/welaw/pkg/errs"
)

func decodeHTTPCreateLawResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.CreateLawResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func decodeHTTPCreateLawsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.CreateLawsResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func decodeHTTPGetLawResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, errs.ErrorDecoder(resp)
	}
	var response endpoints.GetLawResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
