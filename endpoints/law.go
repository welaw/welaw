package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/services"
)

func (e Endpoints) CreateAnnotation(ctx context.Context, ann *apiv1.Annotation) (string, error) {
	req := apiv1.Annotation{}
	resp, err := e.CreateAnnotationEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	r := resp.(CreateAnnotationResponse)
	return r.Id, r.Err
}

func (e Endpoints) DeleteAnnotation(ctx context.Context, id string) error {
	req := DeleteAnnotationRequest{Id: id}
	resp, err := e.DeleteAnnotationEndpoint(ctx, req)
	if err != nil {
		return err
	}
	r := resp.(DeleteAnnotationResponse)
	return r.Err
}

func (e Endpoints) ListAnnotations(ctx context.Context, opts *apiv1.ListAnnotationsOptions) ([]*apiv1.Annotation, int, error) {
	req := ListAnnotationsRequest{}
	resp, err := e.ListAnnotationsEndpoint(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	r := resp.(ListAnnotationsResponse)
	if r.Err != nil {
		return nil, 0, err
	}
	return r.Rows, r.Total, nil
}

func (e Endpoints) DeleteComment(ctx context.Context, uid string) error {
	req := DeleteCommentRequest{Uid: uid}
	resp, err := e.DeleteCommentEndpoint(ctx, req)
	if err != nil {
		return err
	}
	r := resp.(DeleteCommentResponse)
	return r.Err
}

func (e Endpoints) GetComment(ctx context.Context, opts *apiv1.GetCommentOptions) (*apiv1.Comment, error) {
	req := GetCommentRequest{Opts: opts}
	resp, err := e.GetCommentEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(GetCommentResponse)
	if r.Err != nil {
		return nil, err
	}
	return r.Comment, nil
}

func (e Endpoints) ListComments(ctx context.Context, opts *apiv1.ListCommentsOptions) ([]*apiv1.Comment, int, error) {
	req := ListCommentsRequest{}
	resp, err := e.ListCommentsEndpoint(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	r := resp.(ListCommentsResponse)
	if r.Err != nil {
		return nil, 0, err
	}
	return r.Rows, r.Total, nil
}

func (e Endpoints) CreateComment(ctx context.Context, comment *apiv1.Comment) (*apiv1.Comment, error) {
	req := CreateCommentRequest{Comment: comment}
	resp, err := e.CreateCommentEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(CreateCommentResponse)
	if r.Err != nil {
		return nil, err
	}
	return r.Comment, nil
}

func (e Endpoints) UpdateComment(ctx context.Context, comment *apiv1.Comment) (*apiv1.Comment, error) {
	req := UpdateCommentRequest{Comment: comment}
	resp, err := e.UpdateCommentEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(UpdateCommentResponse)
	if r.Err != nil {
		return nil, err
	}
	return r.Comment, nil
}

func (e Endpoints) LikeComment(ctx context.Context, opts *apiv1.LikeCommentOptions) error {
	req := LikeCommentRequest{Opts: opts}
	resp, err := e.LikeCommentEndpoint(ctx, req)
	if err != nil {
		return err
	}
	r := resp.(LikeCommentResponse)
	return r.Err
}

func (e Endpoints) CreateLaw(ctx context.Context, set *apiv1.LawSet, opts *apiv1.CreateLawOptions) (*apiv1.CreateLawReply, error) {
	req := CreateLawRequest{LawSet: set, Opts: opts}
	resp, err := e.CreateLawEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(CreateLawResponse)
	return &apiv1.CreateLawReply{
		LawSet: r.LawSet,
	}, r.Err
}

func (e Endpoints) CreateLaws(ctx context.Context, sets []*apiv1.LawSet, opts *apiv1.CreateLawsOptions) ([]*apiv1.LawSet, error) {
	req := CreateLawsRequest{Laws: sets, Opts: opts}
	resp, err := e.CreateLawsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(CreateLawsResponse)
	return r.Laws, r.Err
}

func (e Endpoints) DeleteLaw(ctx context.Context, upstream, ident string, opts *apiv1.DeleteLawOptions) error {
	req := DeleteLawRequest{
		Upstream: upstream,
		Ident:    ident,
		Opts:     opts,
	}
	resp, err := e.DeleteLawEndpoint(ctx, req)
	if err != nil {
		return err
	}
	r := resp.(DeleteLawResponse)
	return r.Err
}

func (e Endpoints) DiffLaws(ctx context.Context, upstream, ident string, opts *apiv1.DiffLawsOptions) (*apiv1.DiffLawsReply, error) {
	req := DiffLawsRequest{
		Upstream: upstream,
		Ident:    ident,
		Opts:     opts,
	}
	resp, err := e.DiffLawsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(DiffLawsResponse)
	return &apiv1.DiffLawsReply{
		Diff:   r.Diff,
		Theirs: r.Theirs,
	}, r.Err
}

func (e Endpoints) GetLaw(ctx context.Context, upstream, ident string, opts *apiv1.GetLawOptions) (*apiv1.GetLawReply, error) {
	req := GetLawRequest{
		Upstream: upstream,
		Ident:    ident,
		Opts:     opts,
	}
	resp, err := e.GetLawEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(GetLawResponse)
	return &apiv1.GetLawReply{
		LawSet: r.LawSet,
	}, r.Err
}

func (e Endpoints) ListLaws(ctx context.Context, opts *apiv1.ListLawsOptions) (*apiv1.ListLawsReply, error) {
	req := ListLawsRequest{Opts: opts}
	resp, err := e.ListLawsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(ListLawsResponse)
	return &apiv1.ListLawsReply{LawSets: r.Laws, Total: int32(r.Total)}, r.Err
}

func (e Endpoints) UpdateLaw(ctx context.Context, set *apiv1.LawSet, opts *apiv1.UpdateLawOptions) error {
	req := UpdateLawRequest{Law: set, Opts: opts}
	resp, err := e.UpdateLawEndpoint(ctx, req)
	if err != nil {
		return err
	}
	r := resp.(UpdateLawResponse)
	return r.Err
}

func MakeCreateAnnotationEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(apiv1.Annotation)
		id, err := svc.CreateAnnotation(ctx, &req)
		return CreateAnnotationResponse{Id: id, Err: err}, nil
	}
}

func MakeDeleteAnnotationEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteAnnotationRequest)
		err := svc.DeleteAnnotation(ctx, req.Id)
		return DeleteAnnotationResponse{Err: err}, nil
	}
}

func MakeListAnnotationsEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListAnnotationsRequest)
		rows, total, err := svc.ListAnnotations(ctx, req.Opts)
		if err != nil {
			return ListAnnotationsResponse{Err: err}, nil
		}
		return ListAnnotationsResponse{Rows: rows, Total: total}, nil
	}
}

func MakeGetCommentEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetCommentRequest)
		rep, err := svc.GetComment(ctx, req.Opts)
		if err != nil {
			if req.Opts != nil && req.Opts.GetQuiet() {
				return GetCommentResponse{}, nil
			}
			return GetCommentResponse{Err: err}, nil
		}
		return GetCommentResponse{Comment: rep}, nil
	}
}

func MakeListCommentsEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListCommentsRequest)
		rows, total, err := svc.ListComments(ctx, req.Opts)
		if err != nil {
			return ListCommentsResponse{Err: err}, nil
		}
		return ListCommentsResponse{Rows: rows, Total: total}, nil
	}
}

func MakeCreateCommentEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateCommentRequest)
		rep, err := svc.CreateComment(ctx, req.Comment)
		if err != nil {
			return CreateCommentResponse{Err: err}, nil
		}
		return CreateCommentResponse{Comment: rep}, nil
	}
}

func MakeDeleteCommentEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteCommentRequest)
		err := svc.DeleteComment(ctx, req.Uid)
		return DeleteCommentResponse{Err: err}, nil
	}
}

func MakeUpdateCommentEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateCommentRequest)
		rep, err := svc.UpdateComment(ctx, req.Comment)
		if err != nil {
			return UpdateCommentResponse{Err: err}, nil
		}
		return UpdateCommentResponse{Comment: rep}, nil
	}
}

func MakeLikeCommentEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LikeCommentRequest)
		err := svc.LikeComment(ctx, req.Opts)
		return LikeCommentResponse{Err: err}, nil
	}
}

func MakeCreateLawEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateLawRequest)
		rep, err := svc.CreateLaw(ctx, req.LawSet, req.Opts)
		if err != nil {
			return CreateLawResponse{Err: err}, nil
		}
		return CreateLawResponse{LawSet: rep.LawSet}, nil
	}
}

func MakeCreateLawsEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateLawsRequest)
		laws, err := svc.CreateLaws(ctx, req.Laws, req.Opts)
		if err != nil {
			return CreateLawsResponse{Err: err}, nil
		}
		return CreateLawsResponse{Laws: laws}, nil
	}
}

func MakeDeleteLawEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteLawRequest)
		err := svc.DeleteLaw(ctx, req.Upstream, req.Ident, req.Opts)
		return DeleteLawResponse{Err: err}, nil
	}
}

func MakeDiffLawsEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DiffLawsRequest)
		r, err := svc.DiffLaws(ctx, req.Upstream, req.Ident, req.Opts)
		if err != nil {
			return DiffLawsResponse{Err: err}, nil
		}
		return DiffLawsResponse{Diff: r.Diff, Theirs: r.Theirs}, nil
	}
}

func MakeGetLawEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetLawRequest)
		rep, err := svc.GetLaw(ctx, req.Upstream, req.Ident, req.Opts)
		if err != nil {
			return GetLawResponse{Err: err}, nil
		}
		return GetLawResponse{LawSet: rep.LawSet, Err: err}, nil
	}
}

func MakeListLawsEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListLawsRequest)
		resp, err := svc.ListLaws(ctx, req.Opts)
		if err != nil {
			return ListLawsResponse{Err: err}, nil
		}
		return ListLawsResponse{Laws: resp.LawSets, Total: int(resp.Total)}, nil
	}
}

func MakeUpdateLawEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateLawRequest)
		err := svc.UpdateLaw(ctx, req.Law, req.Opts)
		return UpdateLawResponse{Err: err}, nil
	}
}

type CreateAnnotationResponse struct {
	Id  string `json:"id"`
	Err error  `json:"-"`
}

func (r CreateAnnotationResponse) Failed() error { return r.Err }

type DeleteAnnotationRequest struct {
	Id string `json:"id"`
}

type DeleteAnnotationResponse struct {
	Err error `json:"-"`
}

func (r DeleteAnnotationResponse) Failed() error { return r.Err }

type ListAnnotationsRequest struct {
	Opts *apiv1.ListAnnotationsOptions `json:"opts"`
}

type ListAnnotationsResponse struct {
	Rows  []*apiv1.Annotation `json:"rows"`
	Total int                 `json:"total"`
	Err   error               `json:"-"`
}

func (r ListAnnotationsResponse) Failed() error { return r.Err }

type GetCommentRequest struct {
	Opts *apiv1.GetCommentOptions `json:"opts"`
}

type GetCommentResponse struct {
	Comment *apiv1.Comment `json:"comment"`
	Err     error          `json:"-"`
}

func (r GetCommentResponse) Failed() error { return r.Err }

type ListCommentsRequest struct {
	Opts *apiv1.ListCommentsOptions `json:"opts"`
}

type ListCommentsResponse struct {
	Rows  []*apiv1.Comment `json:"rows"`
	Total int              `json:"total"`
	Err   error            `json:"-"`
}

func (r ListCommentsResponse) Failed() error { return r.Err }

type DeleteCommentRequest struct {
	Uid string `json:"uid"`
}

type DeleteCommentResponse struct {
	Err error `json:"-"`
}

func (r DeleteCommentResponse) Failed() error { return r.Err }

type UpdateCommentRequest struct {
	Comment *apiv1.Comment `json:"comment"`
}

type UpdateCommentResponse struct {
	Comment *apiv1.Comment `json:"comment"`
	Err     error          `json:"-"`
}

func (r UpdateCommentResponse) Failed() error { return r.Err }

type LikeCommentRequest struct {
	Opts *apiv1.LikeCommentOptions `json:"opts"`
}

type LikeCommentResponse struct {
	Err error `json:"-"`
}

func (r LikeCommentResponse) Failed() error { return r.Err }

type CreateCommentRequest struct {
	Comment *apiv1.Comment `json:"comment"`
}

type CreateCommentResponse struct {
	Comment *apiv1.Comment `json:"comment"`
	Err     error          `json:"-"`
}

func (r CreateCommentResponse) Failed() error { return r.Err }

type CreateLawRequest struct {
	LawSet *apiv1.LawSet           `json:"law_set"`
	Opts   *apiv1.CreateLawOptions `json:"opts"`
}

type CreateLawResponse struct {
	LawSet *apiv1.LawSet `json:"law_set"`
	Err    error         `json:"-"`
}

func (r CreateLawResponse) Failed() error { return r.Err }

type CreateLawsRequest struct {
	Laws []*apiv1.LawSet          `json:"laws"`
	Opts *apiv1.CreateLawsOptions `json:"opts"`
}

type CreateLawsResponse struct {
	Laws []*apiv1.LawSet `json:"laws"`
	Err  error           `json:"-"`
}

func (r CreateLawsResponse) Failed() error { return r.Err }

type DeleteLawRequest struct {
	Upstream string                  `json:"upstream"`
	Ident    string                  `json:"ident"`
	Opts     *apiv1.DeleteLawOptions `json:"opts"`
}

type DeleteLawResponse struct {
	Err error `json"-"`
}

func (r DeleteLawResponse) Failed() error { return r.Err }

type DiffLawsRequest struct {
	Upstream string                 `json:"upstream"`
	Ident    string                 `json:"ident"`
	Opts     *apiv1.DiffLawsOptions `json:"opts"`
}

type DiffLawsResponse struct {
	Diff   string        `json:"diff"`
	Theirs *apiv1.LawSet `json:"theirs"`
	Err    error         `json:"-"`
}

func (r DiffLawsResponse) Failed() error { return r.Err }

type GetLawRequest struct {
	Upstream string               `json:"upstream"`
	Ident    string               `json:"ident"`
	Opts     *apiv1.GetLawOptions `json:"opts"`
}

type GetLawResponse struct {
	LawSet *apiv1.LawSet `json:"law_set"`
	Err    error         `json:"-"`
}

func (r GetLawResponse) Failed() error { return r.Err }

type ListLawsRequest struct {
	Opts *apiv1.ListLawsOptions `json:"opts"`
}

type ListLawsResponse struct {
	Laws  []*apiv1.LawSet `json:"laws"`
	Total int             `json:"total"`
	Err   error           `json:"-"`
}

func (r ListLawsResponse) Failed() error { return r.Err }

type UpdateLawRequest struct {
	Law  *apiv1.LawSet           `json:"law"`
	Opts *apiv1.UpdateLawOptions `json:"opts"`
}

type UpdateLawResponse struct {
	Err error `json:"-"`
}

func (r UpdateLawResponse) Failed() error { return r.Err }
