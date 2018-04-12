package logging

import (
	"context"
	"fmt"
	"time"

	apiv1 "github.com/welaw/welaw/api/v1"
)

func (mw loggingMiddleware) CreateVote(ctx context.Context, vote *apiv1.Vote, opts *apiv1.CreateVoteOptions) (v *apiv1.Vote, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "create_vote",
			"vote", vote,
			"opts", opts,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	v, err = mw.next.CreateVote(ctx, vote, opts)
	return
}

func (mw loggingMiddleware) CreateVotes(ctx context.Context, votes []*apiv1.Vote, opts *apiv1.CreateVotesOptions) (v []*apiv1.Vote, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "create_votes",
			"votes", len(votes),
			"opts", opts,
			"return", len(v),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	v, err = mw.next.CreateVotes(ctx, votes, opts)
	return
}

func (mw loggingMiddleware) GetVote(ctx context.Context, upstream, ident string, opts *apiv1.GetVoteOptions) (v *apiv1.Vote, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "get_vote",
			"upstream", upstream,
			"ident", ident,
			"opts", opts,
			"vote", fmt.Sprintf("%+v", v),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	v, err = mw.next.GetVote(ctx, upstream, ident, opts)
	return
}

func (mw loggingMiddleware) DeleteVote(ctx context.Context, upstream, ident string, opts *apiv1.DeleteVoteOptions) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "delete_vote",
			"upstream", upstream,
			"ident", ident,
			"opts", opts,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.DeleteVote(ctx, upstream, ident, opts)
	return
}

func (mw loggingMiddleware) ListVotes(ctx context.Context, opts *apiv1.ListVotesOptions) (rep *apiv1.ListVotesReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "list_votes",
			"opts", opts,
			"reply", fmt.Sprintf("%+v", rep),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	rep, err = mw.next.ListVotes(ctx, opts)
	return
}

func (mw loggingMiddleware) UpdateVote(ctx context.Context, vote *apiv1.Vote, opts *apiv1.UpdateVoteOptions) (v *apiv1.Vote, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "update_vote",
			"vote", vote,
			"opts", opts,
			"response", v,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	v, err = mw.next.UpdateVote(ctx, vote, opts)
	return
}
