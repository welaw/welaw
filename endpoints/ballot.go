package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/services"
)

func (e Endpoints) CreateVote(ctx context.Context, vote *apiv1.Vote, opts *apiv1.CreateVoteOptions) (*apiv1.Vote, error) {
	req := CreateVoteRequest{
		Vote: vote,
		Opts: opts,
	}
	resp, err := e.CreateVoteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(CreateVoteResponse)
	return r.Vote, r.Err
}

func (e Endpoints) CreateVotes(ctx context.Context, votes []*apiv1.Vote, opts *apiv1.CreateVotesOptions) ([]*apiv1.Vote, error) {
	req := CreateVotesRequest{
		Votes: votes,
		Opts:  opts,
	}
	resp, err := e.CreateVotesEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(CreateVotesResponse)
	return r.Votes, r.Err
}

func (e Endpoints) GetVote(ctx context.Context, upstream, ident string, opts *apiv1.GetVoteOptions) (*apiv1.Vote, error) {
	req := GetVoteRequest{
		Upstream: upstream,
		Ident:    ident,
		Opts:     opts,
	}
	resp, err := e.GetVoteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(GetVoteResponse)
	return r.Vote, r.Err
}

func (e Endpoints) DeleteVote(ctx context.Context, upstream, ident string, opts *apiv1.DeleteVoteOptions) error {
	req := DeleteVoteRequest{
		Upstream: upstream,
		Ident:    ident,
		Opts:     opts,
	}
	resp, err := e.DeleteVoteEndpoint(ctx, req)
	if err != nil {
		return err
	}
	r := resp.(DeleteVoteResponse)
	return r.Err
}

func (e Endpoints) ListVotes(ctx context.Context, opts *apiv1.ListVotesOptions) (*apiv1.ListVotesReply, error) {
	req := ListVotesRequest{Opts: opts}
	resp, err := e.ListVotesEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(ListVotesResponse)
	return &apiv1.ListVotesReply{
		Votes: r.Votes,
		Total: r.Total,
	}, r.Err
}

func (e Endpoints) UpdateVote(ctx context.Context, v *apiv1.Vote, opts *apiv1.UpdateVoteOptions) (*apiv1.Vote, error) {
	req := UpdateVoteRequest{
		Vote: v,
		Opts: opts,
	}
	resp, err := e.UpdateVoteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(UpdateVoteResponse)
	if r.Err != nil {
		return nil, r.Err
	}
	return r.Vote, nil
}

func MakeCreateVoteEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateVoteRequest)
		vote, err := svc.CreateVote(ctx, req.Vote, req.Opts)
		return CreateVoteResponse{Vote: vote, Err: err}, nil
	}
}

func MakeCreateVotesEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateVotesRequest)
		votes, err := svc.CreateVotes(ctx, req.Votes, req.Opts)
		return CreateVotesResponse{Votes: votes, Err: err}, nil
	}
}

func MakeDeleteVoteEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteVoteRequest)
		err := svc.DeleteVote(ctx, req.Upstream, req.Ident, req.Opts)
		return DeleteVoteResponse{Err: err}, nil
	}
}

func MakeGetVoteEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetVoteRequest)
		v, err := svc.GetVote(ctx, req.Upstream, req.Ident, req.Opts)
		if err != nil {
			if req.Opts != nil && req.Opts.GetQuiet() {
				return GetVoteResponse{}, nil
			}
		}
		return GetVoteResponse{Vote: v, Err: err}, nil
	}
}

func MakeListVotesEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListVotesRequest)
		rep, err := svc.ListVotes(ctx, req.Opts)
		if err != nil {
			return ListVotesResponse{Err: err}, nil
		}
		return ListVotesResponse{
			Votes: rep.Votes,
			Total: rep.Total,
		}, nil
	}
}

func MakeUpdateVoteEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateVoteRequest)
		vote, err := svc.UpdateVote(ctx, req.Vote, req.Opts)
		if err != nil {
			return UpdateVoteResponse{Err: err}, nil
		}
		return UpdateVoteResponse{Vote: vote}, nil
	}
}

type CreateVoteRequest struct {
	Vote *apiv1.Vote              `json:"vote"`
	Opts *apiv1.CreateVoteOptions `json:"opts"`
}

type CreateVoteResponse struct {
	Vote *apiv1.Vote `json:"vote"`
	Err  error       `json:"-"`
}

func (r CreateVoteResponse) Failed() error { return r.Err }

type CreateVotesRequest struct {
	Votes []*apiv1.Vote             `json:"votes"`
	Opts  *apiv1.CreateVotesOptions `json:"opts"`
}

type CreateVotesResponse struct {
	Votes []*apiv1.Vote `json:"votes"`
	Err   error         `json:"-"`
}

func (r CreateVotesResponse) Failed() error { return r.Err }

type UpdateVoteRequest struct {
	Vote *apiv1.Vote              `json:"vote"`
	Opts *apiv1.UpdateVoteOptions `json:"opts"`
}

type UpdateVoteResponse struct {
	Vote *apiv1.Vote `json:"vote"`
	Err  error       `json"-"`
}

func (r UpdateVoteResponse) Failed() error { return r.Err }

type GetVoteRequest struct {
	Upstream string                `json:"upstream"`
	Ident    string                `json:"ident"`
	Opts     *apiv1.GetVoteOptions `json:"opts"`
}

type GetVoteResponse struct {
	Vote *apiv1.Vote `json:"vote"`
	Err  error       `json:"-"`
}

func (r GetVoteResponse) Failed() error { return r.Err }

type DeleteVoteRequest struct {
	Upstream string                   `json:"upstream"`
	Ident    string                   `json:"ident"`
	Opts     *apiv1.DeleteVoteOptions `json:"opts"`
}

type DeleteVoteResponse struct {
	Err error `json:"-"`
}

func (r DeleteVoteResponse) Failed() error { return r.Err }

type ListVotesRequest struct {
	Opts *apiv1.ListVotesOptions `json:"opts"`
}

type ListVotesResponse struct {
	Votes []*apiv1.Vote `json:"votes"`
	Total int32         `json:"total"`
	Err   error         `json:"-"`
}

func (r ListVotesResponse) Failed() error { return r.Err }
