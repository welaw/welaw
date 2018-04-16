package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/pkg/generate"
	"github.com/welaw/welaw/pkg/oauth"
	"github.com/welaw/welaw/pkg/permissions"
	"github.com/welaw/welaw/proto"
)

func (svc service) LoggedInCheck(ctx context.Context) (*proto.User, error) {
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

	if perm, err := svc.hasPermission(uid, permissions.OpUserView, user); err != nil {
		return nil, err
	} else if perm {
		roles, err := svc.db.ListUserRoles(user.Username)
		if err != nil {
			return nil, err
		}
		var roleNames []string
		for _, r := range roles {
			roleNames = append(roleNames, r.Name)
		}
		user.Roles = roleNames
	}

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
	var user *proto.User
	u, err := svc.db.GetUserByProviderId(auth.ProviderId, true)
	switch {
	case err == errs.ErrNotFound:
		// register new user
		user, err = svc.createNewUser(&proto.User{
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

func (svc service) addUserRoles(ctx context.Context, user *proto.User, role string) error {
	return svc.db.CreateUserRoles(user.Username, []string{role})
}

func (svc service) deleteUserRoles(ctx context.Context, user *proto.User, role string) error {
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

func (svc service) hasPermission(user_id, operation string, target interface{}) (perm bool, err error) {

	// owner, group, other

	//switch operation {
	//case permissions.OpUserView, permissions.OpLawView:
	//return true, nil
	//default:
	//return false, nil
	//}
	//}

	if user_id == "" {
		return false, nil
	}

	// create
	// if user -> permission and group
	// if law -> permission and group
	// if branch -> anybody
	// if version -> anybody
	// if upstream -> permission
	// if vote -> anybody

	// list
	// if user -> permission and group unless upstream target
	// if law -> anybody
	// if branch -> anybody
	// if version -> anybody
	// if upstream -> anybody
	// if vote -> anybody

	// update, delete
	// if user -> if owner or permission and group
	// if law -> if owner or permission and group
	// if branch -> if owner or permission and group
	// if version -> if owner or permission and group
	// if upstream -> if owner or permission and group
	// if vote -> if owner or permission and group

	// view
	// if user -> permission and group unless upstream target
	// if law -> anybody
	// if branch -> anybody
	// if version -> anybody
	// if upstream -> anybody
	// if vote -> anybody

	// get owner and upstream group of target object
	var (
		//user     *proto.User
		owner_id string
		upstream string
	)
	if u, ok := target.(*proto.User); ok {
		//user = u.Uid
		owner_id = u.Uid
		upstream = u.GetUpstream()
	} else if u, ok := target.(*proto.Upstream); ok {
		owner_id = u.GetUserId()
		upstream = u.GetIdent()
	} else if law, ok := target.(*proto.Law); ok {
		owner_id = law.GetUserId()
		upstream = law.GetUpstream()
	} else if branch, ok := target.(*proto.Branch); ok {
		owner_id = branch.GetUserId()
		upstream = branch.GetUpstream()
	} else if version, ok := target.(*proto.Version); ok {
		owner_id = version.GetUserId()
		upstream = version.GetUpstreamGroup()
	} else if vote, ok := target.(*proto.Vote); ok {
		user, err := svc.db.GetUserById(vote.GetUserId(), false)
		if err != nil {
			return false, err
		}
		owner_id = user.GetUid()
		upstream = vote.GetUpstream()
	} else if c, ok := target.(*proto.Comment); ok {
		owner_id = c.GetUserId()
		upstream = c.GetUpstream()
	} else {
		return false, errs.BadRequest("unknown target: %+v", target)
	}

	// check if owner
	if user_id == owner_id {
		return true, nil
	}

	if perm, err := svc.db.HasPermission(user_id, operation); err != nil {
		return false, err
	} else if !perm {
		return false, nil
	}

	// scope
	scope, err := svc.db.GetUserAuthScope(user_id, operation)
	if err != nil {
		return false, err
	} else if len(scope) == 0 {
		return false, nil
	}
	for _, s := range scope {
		if s == upstream {
			return true, nil
		}
	}

	return false, nil
}

func (svc service) LoginAs(ctx context.Context, authorization *proto.User) error {
	return nil
}

func (svc service) Logout(ctx context.Context) error {
	return nil
}
