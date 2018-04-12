package logging

import (
	"context"
	"fmt"
	"time"

	apiv1 "github.com/welaw/welaw/api/v1"
)

func (mw loggingMiddleware) CreateUser(ctx context.Context, user *apiv1.User) (u *apiv1.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "create_user",
			"user", fmt.Sprintf("%+v", user),
			"response", fmt.Sprintf("%+v", u),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	u, err = mw.next.CreateUser(ctx, user)
	return
}

func (mw loggingMiddleware) CreateUsers(ctx context.Context, users []*apiv1.User) (u []*apiv1.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "create_users",
			"users", fmt.Sprintf("%+v", users),
			"response", fmt.Sprintf("%+v", u),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	u, err = mw.next.CreateUsers(ctx, users)
	return
}

func (mw loggingMiddleware) DeleteUser(ctx context.Context, username string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "delete_user",
			"username", username,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.DeleteUser(ctx, username)
	return
}

func (mw loggingMiddleware) GetUser(ctx context.Context, opts *apiv1.GetUserOptions) (user *apiv1.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "get_user",
			"opts", fmt.Sprintf("%+v", opts),
			"user", fmt.Sprintf("%v", user),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	user, err = mw.next.GetUser(ctx, opts)
	return
}

func (mw loggingMiddleware) ListUsers(ctx context.Context, opts *apiv1.ListUsersOptions) (users []*apiv1.User, total int, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "list_users",
			"opts", opts,
			"users", fmt.Sprintf("%v", users),
			"total", total,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	users, total, err = mw.next.ListUsers(ctx, opts)
	return
}

func (mw loggingMiddleware) UpdateUser(ctx context.Context, username string, opts *apiv1.UpdateUserOptions) (u *apiv1.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "update_user",
			"username", username,
			"opts", fmt.Sprintf("%+v", opts),
			"response", u,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	u, err = mw.next.UpdateUser(ctx, username, opts)
	return
}

func (mw loggingMiddleware) UploadAvatar(ctx context.Context, opts *apiv1.UploadAvatarOptions) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "upload_avatar",
			"filename", opts.Filename,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.UploadAvatar(ctx, opts)
	return
}
