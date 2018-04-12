package instrumentation

import (
	"context"
	"fmt"
	"time"

	apiv1 "github.com/welaw/welaw/api/v1"
)

func (mw instrumentatingMiddleware) CreateUpstream(ctx context.Context, u *apiv1.Upstream) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "create_upstream", "error", ""}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.CreateUpstream(ctx, u)
	return
}

func (mw instrumentatingMiddleware) GetUpstream(ctx context.Context, ident string) (u *apiv1.Upstream, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "get_upstream", "error", ""}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	u, err = mw.next.GetUpstream(ctx, ident)
	return
}

func (mw instrumentatingMiddleware) ListUpstreams(ctx context.Context) (res []*apiv1.Upstream, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "list_upstreams", "error", ""}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	res, err = mw.next.ListUpstreams(ctx)
	return
}

func (mw instrumentatingMiddleware) UpdateUpstream(ctx context.Context, u *apiv1.Upstream) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "update_upstream", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.UpdateUpstream(ctx, u)
	return
}
