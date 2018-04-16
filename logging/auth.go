package logging

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/welaw/welaw/proto"
)

func (mw loggingMiddleware) LoggedInCheck(ctx context.Context) (user *proto.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "logged_in_check",
			"user", fmt.Sprintf("%+v", user),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	user, err = mw.next.LoggedInCheck(ctx)
	return
}

func (mw loggingMiddleware) Login(ctx context.Context, returnUrl, provider string) (c *http.Cookie, url string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "login",
			"return_url", returnUrl,
			"provider", provider,
			"err", err,
			"url", url,
			"took", time.Since(begin),
		)
	}(time.Now())
	c, url, err = mw.next.Login(ctx, returnUrl, provider)
	return
}

func (mw loggingMiddleware) LoginCallback(ctx context.Context, state, code string) (c, lc *http.Cookie, url string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "login_callback",
			"state", state,
			"code", code,
			"url", url,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	c, lc, url, err = mw.next.LoginCallback(ctx, state, code)
	return
}

func (mw loggingMiddleware) LoginAs(ctx context.Context, user *proto.User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "login_as",
			"user", fmt.Sprintf("%+v", user),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.LoginAs(ctx, user)
	return
}

func (mw loggingMiddleware) Logout(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "logout",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.Logout(ctx)
	return
}
