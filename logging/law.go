package logging

import (
	"context"
	"fmt"
	"time"

	apiv1 "github.com/welaw/welaw/api/v1"
)

func (mw loggingMiddleware) CreateAnnotation(ctx context.Context, ann *apiv1.Annotation) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "create_annotation",
			"ann", fmt.Sprintf("%+v", ann),
			"id", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	id, err = mw.next.CreateAnnotation(ctx, ann)
	return
}

func (mw loggingMiddleware) DeleteAnnotation(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "delete_annotation",
			"id", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.DeleteAnnotation(ctx, id)
	return
}

func (mw loggingMiddleware) ListAnnotations(ctx context.Context, opts *apiv1.ListAnnotationsOptions) (rows []*apiv1.Annotation, total int, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "list_annotations",
			"opts", fmt.Sprintf("%+v", opts),
			"rows", fmt.Sprintf("%+v", rows),
			"total", total,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	rows, total, err = mw.next.ListAnnotations(ctx, opts)
	return
}

func (mw loggingMiddleware) ListComments(ctx context.Context, opts *apiv1.ListCommentsOptions) (rows []*apiv1.Comment, total int, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "list_comments",
			"opts", fmt.Sprintf("%+v", opts),
			"rows", fmt.Sprintf("%+v", rows),
			"total", total,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	rows, total, err = mw.next.ListComments(ctx, opts)
	return
}

func (mw loggingMiddleware) CreateComment(ctx context.Context, comment *apiv1.Comment) (c *apiv1.Comment, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "create_comment",
			"comment", fmt.Sprintf("%+v", comment),
			"c", fmt.Sprintf("%+v", c),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	c, err = mw.next.CreateComment(ctx, comment)
	return
}

func (mw loggingMiddleware) DeleteComment(ctx context.Context, uid string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "delete_comment",
			"uid", uid,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.DeleteComment(ctx, uid)
	return
}

func (mw loggingMiddleware) GetComment(ctx context.Context, opts *apiv1.GetCommentOptions) (c *apiv1.Comment, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "get_comment",
			"opts", fmt.Sprintf("%+v", opts),
			"response", fmt.Sprintf("%+v", c),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	c, err = mw.next.GetComment(ctx, opts)
	return
}

func (mw loggingMiddleware) UpdateComment(ctx context.Context, comment *apiv1.Comment) (c *apiv1.Comment, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "update_comment",
			"comment", fmt.Sprintf("%+v", comment),
			"c", fmt.Sprintf("%+v", c),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	c, err = mw.next.UpdateComment(ctx, comment)
	return
}

func (mw loggingMiddleware) LikeComment(ctx context.Context, opts *apiv1.LikeCommentOptions) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "like_comment",
			"opts", fmt.Sprintf("%+v", opts),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.LikeComment(ctx, opts)
	return
}

func (mw loggingMiddleware) CreateLaw(ctx context.Context, set *apiv1.LawSet, opts *apiv1.CreateLawOptions) (rep *apiv1.CreateLawReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "create_law",
			"upstream", set.Law.Upstream,
			"ident", set.Law.Ident,
			"title", set.Law.Title,
			"author_email", set.Author.Email,
			"tags", set.Version.Tags,
			"opts", opts,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	rep, err = mw.next.CreateLaw(ctx, set, opts)
	return
}

func (mw loggingMiddleware) CreateLaws(ctx context.Context, sets []*apiv1.LawSet, opts *apiv1.CreateLawsOptions) (l []*apiv1.LawSet, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "create_laws",
			"laws", len(sets),
			"opts", opts,
			"return", len(l),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	l, err = mw.next.CreateLaws(ctx, sets, opts)
	return
}

func (mw loggingMiddleware) DeleteLaw(ctx context.Context, upstream, ident string, opts *apiv1.DeleteLawOptions) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "delete_law",
			"upstream", upstream,
			"ident", ident,
			"opts", opts,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.DeleteLaw(ctx, upstream, ident, opts)
	return
}

func (mw loggingMiddleware) DiffLaws(ctx context.Context, upstream, ident string, opts *apiv1.DiffLawsOptions) (r *apiv1.DiffLawsReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "diff_laws",
			"upstream", upstream,
			"ident", ident,
			"opts", opts,
			"reply", r,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	r, err = mw.next.DiffLaws(ctx, upstream, ident, opts)
	return
}

func (mw loggingMiddleware) GetLaw(ctx context.Context, upstream, ident string, opts *apiv1.GetLawOptions) (rep *apiv1.GetLawReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "get_law",
			"upstream", upstream,
			"ident", ident,
			"opts", opts,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	rep, err = mw.next.GetLaw(ctx, upstream, ident, opts)
	return
}

func (mw loggingMiddleware) ListLaws(ctx context.Context, opts *apiv1.ListLawsOptions) (resp *apiv1.ListLawsReply, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "list_laws",
			"opts", fmt.Sprintf("%+v", opts),
			"response", resp,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	resp, err = mw.next.ListLaws(ctx, opts)
	return
}

func (mw loggingMiddleware) UpdateLaw(ctx context.Context, set *apiv1.LawSet, opts *apiv1.UpdateLawOptions) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "update_law",
			"law", fmt.Sprintf("%+v", set),
			"opts", opts,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.UpdateLaw(ctx, set, opts)
	return
}
