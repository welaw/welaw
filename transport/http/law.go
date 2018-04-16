package http

import (
	"context"
	"encoding/json"
	"fmt"
	stdhttp "net/http"

	"github.com/gorilla/mux"
	"github.com/welaw/welaw/endpoints"
	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/proto"
)

func decodeHTTPCreateAnnotationRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req proto.Annotation
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPDeleteAnnotationRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, fmt.Errorf("id not found")
	}
	return endpoints.DeleteAnnotationRequest{Id: id}, nil
}

func decodeHTTPListAnnotationsRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.ListAnnotationsRequest
	var opts proto.ListAnnotationsOptions
	vars := r.URL.Query()
	v, ok := vars["req_type"]
	if ok && len(v) > 0 {
		switch {
		case v[0] == proto.ListAnnotationsOptions_BY_COMMENT.String():
			opts.ReqType = proto.ListAnnotationsOptions_BY_COMMENT
			v, ok = vars["comment_id"]
			if ok && len(v) > 0 {
				opts.CommentId = v[0]
			}
		}
	}
	req.Opts = &opts
	return req, nil
}

func decodeHTTPDeleteCommentRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.DeleteCommentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPGetCommentRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.GetCommentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPListCommentsRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.ListCommentsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPUpdateCommentRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.UpdateCommentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPLikeCommentRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.LikeCommentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPCreateCommentRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.CreateCommentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPCreateLawRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.CreateLawRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPCreateLawsRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.CreateLawsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPDeleteLawRequest(_ context.Context, req *stdhttp.Request) (interface{}, error) {
	vars := mux.Vars(req)
	upstream, ok := vars["upstream"]
	if !ok {
		return nil, fmt.Errorf("upstream not found")
	}
	ident, ok := vars["ident"]
	if !ok {
		return nil, fmt.Errorf("ident not found")
	}
	return endpoints.DeleteLawRequest{Upstream: upstream, Ident: ident}, nil
}

func decodeHTTPDiffLawsRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.DiffLawsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPGetLawRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.GetLawRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPGetLawResponse(_ context.Context, resp *stdhttp.Response) (interface{}, error) {
	if resp.StatusCode != stdhttp.StatusOK {
		return nil, errs.ErrorDecoder(resp)
	}
	var response endpoints.GetLawResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func decodeHTTPListLawsRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var resp endpoints.ListLawsRequest
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPUpdateLawRequest(_ context.Context, r *stdhttp.Request) (interface{}, error) {
	var req endpoints.UpdateLawRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
