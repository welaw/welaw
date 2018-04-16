package instrumentation

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/welaw/welaw/proto"
)

func (mw instrumentatingMiddleware) LoggedInCheck(ctx context.Context) (user *proto.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "logged_in_check", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.next.LoggedInCheck(ctx)
	return
}

func (mw instrumentatingMiddleware) Login(ctx context.Context, returnUrl, provider string) (c *http.Cookie, url string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "login", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	c, url, err = mw.next.Login(ctx, returnUrl, provider)
	return
}

func (mw instrumentatingMiddleware) LoginCallback(ctx context.Context, state, code string) (c, lc *http.Cookie, url string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "login_callback", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	c, lc, url, err = mw.next.LoginCallback(ctx, state, code)
	return
}

func (mw instrumentatingMiddleware) LoginAs(ctx context.Context, user *proto.User) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "login_as", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.LoginAs(ctx, user)
	return
}

func (mw instrumentatingMiddleware) Logout(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "logout", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.Logout(ctx)
	return
}
