package services

import (
	"context"

	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/pkg/permissions"
	"github.com/welaw/welaw/proto"
)

func (svc service) CreateUpstream(ctx context.Context, u *proto.Upstream) (err error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return errs.ErrUnauthorized
	}
	if perm, err := svc.hasPermission(uid, permissions.OpUpstreamCreate, u); err != nil {
		return err
	} else if !perm {
		return errs.ErrUnauthorized
	}

	err = svc.db.CreateUpstream(u)
	return
}

func (svc service) GetUpstream(ctx context.Context, ident string) (u *proto.Upstream, err error) {
	u, err = svc.db.GetUpstream(ident)
	if err != nil {
		return
	}
	//tags, err := svc.db.ListUpstreamTags(ident)
	//if err != nil {
	//return
	//}
	//u.Tags = tags
	return
}

func (svc service) ListUpstreams(_ context.Context) (upstreams []*proto.Upstream, err error) {
	upstreams, err = svc.db.ListUpstreams()
	return
}

func (svc service) UpdateUpstream(ctx context.Context, u *proto.Upstream) (err error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return errs.ErrUnauthorized
	}
	if perm, err := svc.hasPermission(uid, permissions.OpUpstreamUpdate, u); err != nil {
		return err
	} else if !perm {
		return errs.ErrUnauthorized
	}

	err = svc.db.UpdateUpstream(&proto.Upstream{
		Name:        u.Name,
		Description: u.Description,
	})
	return
}
