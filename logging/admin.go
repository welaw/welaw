package logging

import (
	"context"
	"time"

	apiv1 "github.com/welaw/welaw/api/v1"
)

func (mw loggingMiddleware) GetServerStats(ctx context.Context) (stats *apiv1.ServerStats, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "get_server_stats",
			"stats", stats,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	stats, err = mw.next.GetServerStats(ctx)
	return
}

func (mw loggingMiddleware) LoadRepos(ctx context.Context, opts *apiv1.LoadReposOptions) (r *apiv1.LoadReposReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "load_repos",
			"opts", opts,
			"response", r,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	r, err = mw.next.LoadRepos(ctx, opts)
	return
}

func (mw loggingMiddleware) SaveRepos(ctx context.Context, opts *apiv1.SaveReposOptions) (r *apiv1.SaveReposReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "save_repos",
			"opts", opts,
			"response", r,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	r, err = mw.next.SaveRepos(ctx, opts)
	return
}
