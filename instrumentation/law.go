package instrumentation

import (
	"context"
	"fmt"
	"time"

	"github.com/welaw/welaw/proto"
)

func (mw instrumentatingMiddleware) CreateAnnotation(ctx context.Context, ann *proto.Annotation) (id string, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "create_annotation", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	id, err = mw.next.CreateAnnotation(ctx, ann)
	return
}

func (mw instrumentatingMiddleware) DeleteAnnotation(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		s := []string{"method", "delete_annotation", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.DeleteAnnotation(ctx, id)
	return
}

func (mw instrumentatingMiddleware) ListAnnotations(ctx context.Context, opts *proto.ListAnnotationsOptions) (rows []*proto.Annotation, total int, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "list_annotations", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rows, total, err = mw.next.ListAnnotations(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) CreateComment(ctx context.Context, comment *proto.Comment) (c *proto.Comment, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "create_comment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	c, err = mw.next.CreateComment(ctx, comment)
	return
}

func (mw instrumentatingMiddleware) DeleteComment(ctx context.Context, uid string) (err error) {
	defer func(begin time.Time) {
		s := []string{"method", "delete_comment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.DeleteComment(ctx, uid)
	return
}

func (mw instrumentatingMiddleware) ListComments(ctx context.Context, opts *proto.ListCommentsOptions) (rows []*proto.Comment, total int, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "list_comments", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rows, total, err = mw.next.ListComments(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) GetComment(ctx context.Context, opts *proto.GetCommentOptions) (c *proto.Comment, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "get_comment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	c, err = mw.next.GetComment(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) UpdateComment(ctx context.Context, comment *proto.Comment) (c *proto.Comment, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "update_comment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	c, err = mw.next.UpdateComment(ctx, comment)
	return
}

func (mw instrumentatingMiddleware) LikeComment(ctx context.Context, opts *proto.LikeCommentOptions) (err error) {
	defer func(begin time.Time) {
		s := []string{"method", "like_comment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.LikeComment(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) CreateLaw(ctx context.Context, set *proto.LawSet, opts *proto.CreateLawOptions) (rep *proto.CreateLawReply, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "create_law", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rep, err = mw.next.CreateLaw(ctx, set, opts)
	return
}

func (mw instrumentatingMiddleware) CreateLaws(ctx context.Context, sets []*proto.LawSet, opts *proto.CreateLawsOptions) (l []*proto.LawSet, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "create_laws", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	l, err = mw.next.CreateLaws(ctx, sets, opts)
	return
}

func (mw instrumentatingMiddleware) DeleteLaw(ctx context.Context, u, i string, opts *proto.DeleteLawOptions) (err error) {
	defer func(begin time.Time) {
		s := []string{"method", "delete_law", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.DeleteLaw(ctx, u, i, opts)
	return
}

func (mw instrumentatingMiddleware) DiffLaws(ctx context.Context, upstream, ident string, opts *proto.DiffLawsOptions) (r *proto.DiffLawsReply, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "diff_laws", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	r, err = mw.next.DiffLaws(ctx, upstream, ident, opts)
	return
}

func (mw instrumentatingMiddleware) GetLaw(ctx context.Context, u, i string, opts *proto.GetLawOptions) (rep *proto.GetLawReply, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "get_law", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rep, err = mw.next.GetLaw(ctx, u, i, opts)
	return
}

func (mw instrumentatingMiddleware) ListLaws(ctx context.Context, opts *proto.ListLawsOptions) (resp *proto.ListLawsReply, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "list_laws", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	resp, err = mw.next.ListLaws(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) UpdateLaw(ctx context.Context, law *proto.LawSet, opts *proto.UpdateLawOptions) (err error) {
	defer func(begin time.Time) {
		s := []string{"method", "update_law", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.UpdateLaw(ctx, law, opts)
	return err
}
