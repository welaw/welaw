package instrumentation

import (
	"context"
	"fmt"
	"time"

	"github.com/welaw/welaw/proto"
)

func (mw instrumentatingMiddleware) CreateVote(ctx context.Context, vote *proto.Vote, opts *proto.CreateVoteOptions) (v *proto.Vote, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "create_vote", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	v, err = mw.next.CreateVote(ctx, vote, opts)
	return
}

func (mw instrumentatingMiddleware) CreateVotes(ctx context.Context, votes []*proto.Vote, opts *proto.CreateVotesOptions) (v []*proto.Vote, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "create_votes", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	v, err = mw.next.CreateVotes(ctx, votes, opts)
	return
}

func (mw instrumentatingMiddleware) GetVote(ctx context.Context, upstream, ident string, opts *proto.GetVoteOptions) (v *proto.Vote, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "get_vote", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	v, err = mw.next.GetVote(ctx, upstream, ident, opts)
	return
}

func (mw instrumentatingMiddleware) DeleteVote(ctx context.Context, upstream, ident string, opts *proto.DeleteVoteOptions) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "delete_vote", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.DeleteVote(ctx, upstream, ident, opts)
	return
}

func (mw instrumentatingMiddleware) ListVotes(ctx context.Context, opts *proto.ListVotesOptions) (rep *proto.ListVotesReply, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "list_votes", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rep, err = mw.next.ListVotes(ctx, opts)
	return rep, err
}

func (mw instrumentatingMiddleware) UpdateVote(ctx context.Context, vote *proto.Vote, opts *proto.UpdateVoteOptions) (v *proto.Vote, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "update_vote", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	v, err = mw.next.UpdateVote(ctx, vote, opts)
	return
}
