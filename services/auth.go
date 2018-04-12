package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/pkg/generate"
	"github.com/welaw/welaw/pkg/oauth"
	"github.com/welaw/welaw/pkg/permissions"
)

func (svc service) LoggedInCheck(ctx context.Context) (*apiv1.User, error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, nil
	}
	user, err := svc.db.GetUserById(uid, false)
	switch {
	case err == errs.ErrNotFound:
		return nil, nil
	case err != nil:
		return nil, err
	}

	admin, err := svc.hasPermission(uid, permissions.OpUserView)
	if err != nil {
		return nil, err
	}
	if !admin {
		return user, nil
	}
	roles, err := svc.db.ListUserRoles(user.Username)
	if err != nil {
		return nil, err
	}
	var roleNames []string
	for _, r := range roles {
		roleNames = append(roleNames, r.Name)
	}
	user.Roles = roleNames

	return user, nil
}

func (svc service) Login(ctx context.Context, returnURL, provider string) (*http.Cookie, string, error) {
	p, ok := svc.providers[provider]
	if !ok {
		return nil, "", errs.ErrBadRequest
	}
	r := generate.RandString(10)
	state := fmt.Sprintf("%s,%s,%s", r, returnURL, provider)
	//enc := base64URLEncode(
	url := p.GetLoginURL(state)
	svc.logger.Log("method", "login", "auth_code_url", url)

	value, err := svc.sc.Encode("welaw_auth_login", state)
	if err != nil {
		return nil, "", err
	}
	cookie := http.Cookie{
		Name:     "welaw_auth_login",
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(time.Second * 10),
		MaxAge:   5000,
		Secure:   svc.Opts.UseSecureCookies,
		HttpOnly: true,
	}
	return &cookie, url, nil
}

func (svc service) LoginCallback(ctx context.Context, state, code string) (*http.Cookie, *http.Cookie, string, error) {
	svc.logger.Log("method", "login_callback", "state", state, "code", code)
	prev, ok := ctx.Value("state").(string)
	if !ok {
		return nil, nil, "", fmt.Errorf("state not found")
	}
	if prev != state {
		svc.logger.Log("unauthorized", "invalid state", "prev", prev, "state", state)
		return nil, nil, "", errs.ErrUnauthorized
	}
	f := strings.Split(state, ",")
	if len(f) < 3 {
		return nil, nil, "", errs.ErrBadRequest
	}
	returnURL := svc.Opts.LoginSuccessURL + "?url=" + f[1]
	provider := f[2]
	auth, err := svc.providers[provider].LoginUser(state, code)
	if err != nil {
		return nil, nil, "", err
	}
	var user *apiv1.User
	u, err := svc.db.GetUserByProviderId(auth.ProviderId, true)
	switch {
	case err == errs.ErrNotFound:
		// register new user
		user, err = svc.createNewUser(&apiv1.User{
			Provider:        provider,
			ProviderId:      auth.ProviderId,
			FullName:        auth.Name,
			FullNamePrivate: true,
			Email:           auth.Email,
			EmailPrivate:    true,
			PictureUrl:      auth.PictureUrl,
		})
		if err != nil {
			return nil, nil, "", err
		}
	case err != nil:
		return nil, nil, "", err
	default:
		if u.DeletedAt != "" {
			u.FullNamePrivate = true
			u.EmailPrivate = true
			user, err = svc.createUserWithId(u)
			if err != nil {
				return nil, nil, "", err
			}
		} else {
			user = u
		}
	}

	token := oauth.MakeClaimsToken(user.Uid)
	tkn, err := token.SignedString(svc.Opts.SigningKey)
	if err != nil {
		return nil, nil, "", err
	}
	loginCookie := http.Cookie{
		Name:     "welaw_auth_login",
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		MaxAge:   0,
		Secure:   svc.Opts.UseSecureCookies,
		HttpOnly: true,
	}
	value, err := svc.sc.Encode("welaw_auth", tkn)
	if err != nil {
		return nil, nil, "", err
	}
	cookie := http.Cookie{
		Name:     "welaw_auth",
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 1),
		MaxAge:   50000,
		Secure:   svc.Opts.UseSecureCookies,
		HttpOnly: true,
	}

	err = svc.db.UpdateLastLogin(user.Uid)
	if err != nil {
		return nil, nil, "", err
	}

	return &cookie, &loginCookie, returnURL, nil
}

func (svc service) addUserRoles(ctx context.Context, user *apiv1.User, role string) error {
	return svc.db.CreateUserRoles(user.Username, []string{role})
}

func (svc service) deleteUserRoles(ctx context.Context, user *apiv1.User, role string) error {
	return svc.db.DeleteUserRoles(user.Username, []string{role})
}

func (svc service) authorizeOperation(ctx context.Context, user_id, username, operation string) error {
	return nil
}

func (svc service) authorizeUser(ctx context.Context, user_id, role string) error {
	has, err := svc.db.HasRole(user_id, role)
	if err != nil {
		return err
	}
	if has == false {
		return errs.ErrUnauthorized
	}
	return nil
}

func (svc service) hasPermission(user_id, operation string) (bool, error) {
	if user_id == "" {
		return false, nil
	}
	//switch operation {
	//case permissions.OpUserView, permissions.OpLawView:
	//return true, nil
	//default:
	//return false, nil
	//}
	//}
	p, err := svc.db.HasPermission(user_id, operation)
	if err != nil {
		return false, err
	}
	return p, nil
}

func (svc service) LoginAs(ctx context.Context, authorization *apiv1.User) error {
	return nil
}

func (svc service) Logout(ctx context.Context) error {
	return nil
}
