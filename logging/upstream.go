package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/welaw/welaw/proto"
)

func (mw loggingMiddleware) CreateUpstream(ctx context.Context, u *proto.Upstream) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "create_upstream",
			"upstream", fmt.Sprintf("%+v", u),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.CreateUpstream(ctx, u)
	return
}

func (mw loggingMiddleware) GetUpstream(ctx context.Context, ident string) (u *proto.Upstream, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "get_upstream",
			"ident", ident,
			"upstream", fmt.Sprintf("%+v", u),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	u, err = mw.next.GetUpstream(ctx, ident)
	return
}

func (mw loggingMiddleware) ListUpstreams(ctx context.Context) (upstreams []*proto.Upstream, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "list_upstreams",
			"upstreams", fmt.Sprintf("%v", upstreams),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	upstreams, err = mw.next.ListUpstreams(ctx)
	return
}

func (mw loggingMiddleware) UpdateUpstream(ctx context.Context, u *proto.Upstream) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "update",
			"upstream", fmt.Sprintf("%+v", u),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.UpdateUpstream(ctx, u)
	return
}
