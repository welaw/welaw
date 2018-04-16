package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/welaw/welaw/proto"
	"github.com/welaw/welaw/services"
)

func (e Endpoints) GetServerStats(ctx context.Context) (*proto.ServerStats, error) {
	resp, err := e.GetServerStatsEndpoint(ctx, GetServerStatsRequest{})
	if err != nil {
		return nil, err
	}
	r := resp.(GetServerStatsResponse)
	return r.Stats, r.Err
}

func (e Endpoints) LoadRepos(ctx context.Context, opts *proto.LoadReposOptions) (*proto.LoadReposReply, error) {
	resp, err := e.LoadReposEndpoint(ctx, opts)
	if err != nil {
		return nil, err
	}
	r := resp.(LoadReposResponse)
	return r.Reply, r.Err
}

func (e Endpoints) SaveRepos(ctx context.Context, opts *proto.SaveReposOptions) (*proto.SaveReposReply, error) {
	resp, err := e.SaveReposEndpoint(ctx, opts)
	if err != nil {
		return nil, err
	}
	r := resp.(SaveReposResponse)
	return r.Reply, r.Err
}

func MakeGetServerStatsEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(GetServerStatsRequest)
		stats, err := svc.GetServerStats(ctx)
		return GetServerStatsResponse{Stats: stats, Err: err}, nil
	}
}

func MakeLoadReposEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoadReposRequest)
		rep, err := svc.LoadRepos(ctx, req.Opts)
		return LoadReposResponse{Reply: rep, Err: err}, nil
	}
}

func MakeSaveReposEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SaveReposRequest)
		rep, err := svc.SaveRepos(ctx, req.Opts)
		return SaveReposResponse{Reply: rep, Err: err}, nil
	}
}

type LoadReposRequest struct {
	Opts *proto.LoadReposOptions `json:"opts"`
}

type LoadReposResponse struct {
	Reply *proto.LoadReposReply `json:"reply"`
	Err   error                 `json:"-"`
}

func (r LoadReposResponse) Failed() error { return r.Err }

type SaveReposRequest struct {
	Opts *proto.SaveReposOptions `json:"opts"`
}

type SaveReposResponse struct {
	Reply *proto.SaveReposReply `json:"reply"`
	Err   error                 `json:"-"`
}

func (r SaveReposResponse) Failed() error { return r.Err }

type GetServerStatsRequest struct{}

type GetServerStatsResponse struct {
	Stats *proto.ServerStats `json:"stats"`
	Err   error              `json:"-"`
}

func (r GetServerStatsResponse) Failed() error { return r.Err }
