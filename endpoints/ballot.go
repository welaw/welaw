package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/welaw/welaw/proto"
	"github.com/welaw/welaw/services"
)

func (e Endpoints) CreateVote(ctx context.Context, vote *proto.Vote, opts *proto.CreateVoteOptions) (*proto.Vote, error) {
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

func (e Endpoints) CreateVotes(ctx context.Context, votes []*proto.Vote, opts *proto.CreateVotesOptions) ([]*proto.Vote, error) {
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

func (e Endpoints) GetVote(ctx context.Context, upstream, ident string, opts *proto.GetVoteOptions) (*proto.Vote, error) {
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

func (e Endpoints) DeleteVote(ctx context.Context, upstream, ident string, opts *proto.DeleteVoteOptions) error {
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

func (e Endpoints) ListVotes(ctx context.Context, opts *proto.ListVotesOptions) (*proto.ListVotesReply, error) {
	req := ListVotesRequest{Opts: opts}
	resp, err := e.ListVotesEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(ListVotesResponse)
	return &proto.ListVotesReply{
		Votes: r.Votes,
		Total: r.Total,
	}, r.Err
}

func (e Endpoints) UpdateVote(ctx context.Context, v *proto.Vote, opts *proto.UpdateVoteOptions) (*proto.Vote, error) {
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
	Vote *proto.Vote              `json:"vote"`
	Opts *proto.CreateVoteOptions `json:"opts"`
}

type CreateVoteResponse struct {
	Vote *proto.Vote `json:"vote"`
	Err  error       `json:"-"`
}

func (r CreateVoteResponse) Failed() error { return r.Err }

type CreateVotesRequest struct {
	Votes []*proto.Vote             `json:"votes"`
	Opts  *proto.CreateVotesOptions `json:"opts"`
}

type CreateVotesResponse struct {
	Votes []*proto.Vote `json:"votes"`
	Err   error         `json:"-"`
}

func (r CreateVotesResponse) Failed() error { return r.Err }

type UpdateVoteRequest struct {
	Vote *proto.Vote              `json:"vote"`
	Opts *proto.UpdateVoteOptions `json:"opts"`
}

type UpdateVoteResponse struct {
	Vote *proto.Vote `json:"vote"`
	Err  error       `json"-"`
}

func (r UpdateVoteResponse) Failed() error { return r.Err }

type GetVoteRequest struct {
	Upstream string                `json:"upstream"`
	Ident    string                `json:"ident"`
	Opts     *proto.GetVoteOptions `json:"opts"`
}

type GetVoteResponse struct {
	Vote *proto.Vote `json:"vote"`
	Err  error       `json:"-"`
}

func (r GetVoteResponse) Failed() error { return r.Err }

type DeleteVoteRequest struct {
	Upstream string                   `json:"upstream"`
	Ident    string                   `json:"ident"`
	Opts     *proto.DeleteVoteOptions `json:"opts"`
}

type DeleteVoteResponse struct {
	Err error `json:"-"`
}

func (r DeleteVoteResponse) Failed() error { return r.Err }

type ListVotesRequest struct {
	Opts *proto.ListVotesOptions `json:"opts"`
}

type ListVotesResponse struct {
	Votes []*proto.Vote `json:"votes"`
	Total int32         `json:"total"`
	Err   error         `json:"-"`
}

func (r ListVotesResponse) Failed() error { return r.Err }
