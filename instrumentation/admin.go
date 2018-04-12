package instrumentation

import (
	"context"
	"fmt"
	"time"

	apiv1 "github.com/welaw/welaw/api/v1"
)

func (mw instrumentatingMiddleware) GetServerStats(ctx context.Context) (stats *apiv1.ServerStats, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "get_server_stats", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	stats, err = mw.next.GetServerStats(ctx)
	return
}

func (mw instrumentatingMiddleware) LoadRepos(ctx context.Context, opts *apiv1.LoadReposOptions) (r *apiv1.LoadReposReply, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "load_repos", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	r, err = mw.next.LoadRepos(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) SaveRepos(ctx context.Context, opts *apiv1.SaveReposOptions) (r *apiv1.SaveReposReply, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "save_repos", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	r, err = mw.next.SaveRepos(ctx, opts)
	return
}
