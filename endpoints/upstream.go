package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/welaw/welaw/proto"
	"github.com/welaw/welaw/services"
)

func (e Endpoints) CreateUpstream(ctx context.Context, u *proto.Upstream) error {
	req := CreateUpstreamRequest{Upstream: u}
	resp, err := e.CreateUpstreamEndpoint(ctx, req)
	if err != nil {
		return err
	}
	r := resp.(CreateUpstreamResponse)
	return r.Err
}

func (e Endpoints) GetUpstream(ctx context.Context, ident string) (*proto.Upstream, error) {
	req := GetUpstreamRequest{Ident: ident}
	resp, err := e.GetUpstreamEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(GetUpstreamResponse)
	return r.Upstream, r.Err
}

func (e Endpoints) ListUpstreams(ctx context.Context) ([]*proto.Upstream, error) {
	req := ListUpstreamsRequest{}
	resp, err := e.ListUpstreamsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(ListUpstreamsResponse)
	return r.Upstreams, r.Err
}

func (e Endpoints) UpdateUpstream(ctx context.Context, upstream *proto.Upstream) error {
	req := UpdateUpstreamRequest{Upstream: upstream}
	resp, err := e.UpdateUpstreamEndpoint(ctx, req)
	if err != nil {
		return err
	}
	r := resp.(UpdateUpstreamResponse)
	return r.Err
}

func MakeCreateUpstreamEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUpstreamRequest)
		err := svc.CreateUpstream(ctx, req.Upstream)
		return CreateUpstreamResponse{Err: err}, nil
	}
}

func MakeGetUpstreamEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUpstreamRequest)
		u, err := svc.GetUpstream(ctx, req.Ident)
		if err != nil {
			return GetUpstreamResponse{Err: err}, nil
		}
		return GetUpstreamResponse{Upstream: u}, nil
	}
}

func MakeListUpstreamsEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ListUpstreamsRequest)
		res, err := svc.ListUpstreams(ctx)
		return ListUpstreamsResponse{Upstreams: res, Err: err}, nil
	}
}

func MakeUpdateUpstreamEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUpstreamRequest)
		err := svc.UpdateUpstream(ctx, req.Upstream)
		return UpdateUpstreamResponse{Err: err}, nil
	}
}

type CreateUpstreamRequest struct {
	Upstream *proto.Upstream `json:"upstream"`
}

type CreateUpstreamResponse struct {
	Err error `json:"-"`
}

func (r CreateUpstreamResponse) Failed() error { return r.Err }

type GetUpstreamRequest struct {
	Ident string `json:"ident"`
}

type GetUpstreamResponse struct {
	Upstream *proto.Upstream `json:"upstream"`
	Laws     []*proto.LawSet `json:"laws"`
	Users    []*proto.User   `json:"user"`
	Err      error           `json:"-"`
}

func (r GetUpstreamResponse) Failed() error { return r.Err }

type ListUpstreamsRequest struct{}

type ListUpstreamsResponse struct {
	Upstreams []*proto.Upstream `json:"upstreams"`
	Err       error             `json:"-"`
}

func (r ListUpstreamsResponse) Failed() error { return r.Err }

type UpdateUpstreamRequest struct {
	Upstream *proto.Upstream `json:"upstream"`
}

type UpdateUpstreamResponse struct {
	Err error `json:"-"`
}

func (r UpdateUpstreamResponse) Failed() error { return r.Err }
