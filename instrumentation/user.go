package instrumentation

import (
	"context"
	"fmt"
	"time"

	"github.com/welaw/welaw/proto"
)

func (mw instrumentatingMiddleware) CreateUser(ctx context.Context, user *proto.User) (u *proto.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "create_user", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	u, err = mw.next.CreateUser(ctx, user)
	return
}

func (mw instrumentatingMiddleware) CreateUsers(ctx context.Context, users []*proto.User) (u []*proto.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "create_users", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	u, err = mw.next.CreateUsers(ctx, users)
	return
}

func (mw instrumentatingMiddleware) DeleteUser(ctx context.Context, username string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "delete_user", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.DeleteUser(ctx, username)
	return
}

func (mw instrumentatingMiddleware) GetUser(ctx context.Context, opts *proto.GetUserOptions) (user *proto.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "get_user", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.next.GetUser(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) ListUsers(ctx context.Context, opts *proto.ListUsersOptions) (users []*proto.User, total int, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "list_users", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	users, total, err = mw.next.ListUsers(ctx, opts)
	return
}

func (mw instrumentatingMiddleware) UpdateUser(ctx context.Context, username string, opts *proto.UpdateUserOptions) (u *proto.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "update_user", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	u, err = mw.next.UpdateUser(ctx, username, opts)
	return
}

func (mw instrumentatingMiddleware) UploadAvatar(ctx context.Context, opts *proto.UploadAvatarOptions) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "upload_avatar", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.UploadAvatar(ctx, opts)
	return
}
