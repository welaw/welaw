package services

import (
	"context"

	apiv1 "github.com/welaw/welaw/api/v1"
)

func (svc service) CreateUpstream(_ context.Context, u *apiv1.Upstream) (err error) {
	err = svc.db.CreateUpstream(u)
	return
}

func (svc service) GetUpstream(ctx context.Context, ident string) (u *apiv1.Upstream, err error) {
	u, err = svc.db.GetUpstream(ident)
	if err != nil {
		return
	}
	tags, err := svc.db.ListUpstreamTags(ident)
	if err != nil {
		return
	}
	u.Tags = tags
	return
}

func (svc service) ListUpstreams(_ context.Context) (upstreams []*apiv1.Upstream, err error) {
	upstreams, err = svc.db.ListUpstreams()
	return
}

func (svc service) UpdateUpstream(_ context.Context, u *apiv1.Upstream) (err error) {
	err = svc.db.UpdateUpstream(&apiv1.Upstream{
		Name:        u.Name,
		Description: u.Description,
	})
	return
}
