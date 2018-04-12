package services

import (
	"context"

	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/pkg/permissions"
)

func (svc service) CreateAnnotation(ctx context.Context, ann *apiv1.Annotation) (string, error) {
	_, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", errs.Unauthorized("user_id not found")
	}
	//perm, err := svc.hasPermission(uid, permissions.OpCommentCreate)
	//if err != nil {
	//return "", err
	//}
	//if !perm {
	//return "", errs.Unauthorized("insufficient permissions")
	//}
	id, err := svc.db.CreateAnnotation(ann)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (svc service) DeleteAnnotation(ctx context.Context, id string) error {
	_, ok := ctx.Value("user_id").(string)
	if !ok {
		return errs.Unauthorized("user_id not found")
	}
	// TODO
	return svc.db.DeleteAnnotationById(id)
}

func (svc service) ListAnnotations(ctx context.Context, opts *apiv1.ListAnnotationsOptions) ([]*apiv1.Annotation, int, error) {
	switch {
	case opts == nil:
		return nil, 0, errs.ErrBadRequest
	case opts.ReqType == apiv1.ListAnnotationsOptions_BY_COMMENT:
		return svc.db.ListAnnotations(opts.CommentId)
	default:
		return nil, 0, errs.ErrBadRequest
	}
}

func (svc service) CreateComment(ctx context.Context, comment *apiv1.Comment) (*apiv1.Comment, error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}
	_, err := svc.hasPermission(uid, permissions.OpCommentCreate)
	if err != nil {
		return nil, err
	}
	//if !perm {
	//return nil, errs.ErrUnauthorized
	//}
	return svc.db.CreateComment(uid, comment)
}

func (svc service) DeleteComment(ctx context.Context, uid string) error {
	return svc.db.DeleteComment(uid)
}

func (svc service) GetComment(ctx context.Context, opts *apiv1.GetCommentOptions) (*apiv1.Comment, error) {
	switch {
	case opts == nil:
		return nil, errs.ErrBadRequest
	case opts.ReqType == apiv1.GetCommentOptions_BY_USER_VERSION:
		return svc.db.GetCommentByUserVersion(opts.Username, opts.Upstream, opts.Ident, opts.Branch, opts.Version)
	case opts.ReqType == apiv1.GetCommentOptions_BY_UID:
		return svc.db.GetCommentByUid(opts.Uid)
	}
	return nil, errs.ErrBadRequest
}

func (svc service) ListComments(ctx context.Context, opts *apiv1.ListCommentsOptions) ([]*apiv1.Comment, int, error) {
	uid, _ := ctx.Value("user_id").(string)
	switch {
	case opts == nil:
		return nil, 0, errs.ErrBadRequest
	case opts.ReqType == apiv1.ListCommentsOptions_BY_USERNAME:
		return svc.db.ListCommentsByUsername(uid, opts.Username)
	case opts.ReqType == apiv1.ListCommentsOptions_BY_VERSION:
		return svc.db.ListCommentsByVersion(
			uid,
			opts.Upstream,
			opts.Ident,
			opts.Branch,
			opts.OrderBy,
			opts.Version,
			opts.PageSize,
			opts.PageNum,
			opts.Desc,
		)
	}
	return nil, 0, errs.ErrBadRequest
}

func (svc service) UpdateComment(ctx context.Context, comment *apiv1.Comment) (c *apiv1.Comment, err error) {
	c, err = svc.db.UpdateComment(comment)
	if err != nil {
		return
	}
	c, err = svc.db.GetCommentByUid(c.Uid)
	if err != nil {
		return
	}
	return
}

func (svc service) LikeComment(ctx context.Context, opts *apiv1.LikeCommentOptions) error {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return errs.ErrUnauthorized
	}
	return svc.db.LikeComment(opts.CommentId, uid)
}
