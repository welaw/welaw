package instrumentation

import (
	"context"
	"fmt"
	"time"

	apiv1 "github.com/welaw/welaw/api/v1"
)

func (mw instrumentatingMiddleware) CreateAnnotation(ctx context.Context, ann *apiv1.Annotation) (id string, err error) {
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

func (mw instrumentatingMiddleware) ListAnnotations(ctx context.Context, opts *apiv1.ListAnnotationsOptions) (rows []*apiv1.Annotation, total int, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "list_annotations", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rows, total, err = mw.next.ListAnnotations(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) CreateComment(ctx context.Context, comment *apiv1.Comment) (c *apiv1.Comment, err error) {
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

func (mw instrumentatingMiddleware) ListComments(ctx context.Context, opts *apiv1.ListCommentsOptions) (rows []*apiv1.Comment, total int, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "list_comments", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rows, total, err = mw.next.ListComments(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) GetComment(ctx context.Context, opts *apiv1.GetCommentOptions) (c *apiv1.Comment, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "get_comment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	c, err = mw.next.GetComment(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) UpdateComment(ctx context.Context, comment *apiv1.Comment) (c *apiv1.Comment, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "update_comment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	c, err = mw.next.UpdateComment(ctx, comment)
	return
}

func (mw instrumentatingMiddleware) LikeComment(ctx context.Context, opts *apiv1.LikeCommentOptions) (err error) {
	defer func(begin time.Time) {
		s := []string{"method", "like_comment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.LikeComment(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) CreateLaw(ctx context.Context, set *apiv1.LawSet, opts *apiv1.CreateLawOptions) (rep *apiv1.CreateLawReply, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "create_law", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rep, err = mw.next.CreateLaw(ctx, set, opts)
	return
}

func (mw instrumentatingMiddleware) CreateLaws(ctx context.Context, sets []*apiv1.LawSet, opts *apiv1.CreateLawsOptions) (l []*apiv1.LawSet, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "create_laws", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	l, err = mw.next.CreateLaws(ctx, sets, opts)
	return
}

func (mw instrumentatingMiddleware) DeleteLaw(ctx context.Context, u, i string, opts *apiv1.DeleteLawOptions) (err error) {
	defer func(begin time.Time) {
		s := []string{"method", "delete_law", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.DeleteLaw(ctx, u, i, opts)
	return
}

func (mw instrumentatingMiddleware) DiffLaws(ctx context.Context, upstream, ident string, opts *apiv1.DiffLawsOptions) (r *apiv1.DiffLawsReply, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "diff_laws", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	r, err = mw.next.DiffLaws(ctx, upstream, ident, opts)
	return
}

func (mw instrumentatingMiddleware) GetLaw(ctx context.Context, u, i string, opts *apiv1.GetLawOptions) (rep *apiv1.GetLawReply, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "get_law", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	rep, err = mw.next.GetLaw(ctx, u, i, opts)
	return
}

func (mw instrumentatingMiddleware) ListLaws(ctx context.Context, opts *apiv1.ListLawsOptions) (resp *apiv1.ListLawsReply, err error) {
	defer func(begin time.Time) {
		s := []string{"method", "list_laws", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	resp, err = mw.next.ListLaws(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) UpdateLaw(ctx context.Context, law *apiv1.LawSet, opts *apiv1.UpdateLawOptions) (err error) {
	defer func(begin time.Time) {
		s := []string{"method", "update_law", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(s...).Add(1)
		mw.requestLatency.With(s...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.UpdateLaw(ctx, law, opts)
	return err
}
