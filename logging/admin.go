package logging

import (
	"context"
	"time"

	"github.com/welaw/welaw/proto"
)

func (mw loggingMiddleware) GetServerStats(ctx context.Context) (stats *proto.ServerStats, err error) {
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

func (mw loggingMiddleware) LoadRepos(ctx context.Context, opts *proto.LoadReposOptions) (r *proto.LoadReposReply, err error) {
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

func (mw loggingMiddleware) SaveRepos(ctx context.Context, opts *proto.SaveReposOptions) (r *proto.SaveReposReply, err error) {
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
