package services

import (
	"context"
	"fmt"
	"strings"

	apiv1 "github.com/welaw/welaw/api/v1"
	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/pkg/permissions"
)

func (svc service) CreateUser(ctx context.Context, user *apiv1.User) (u *apiv1.User, err error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}
	perm, err := svc.hasPermission(uid, permissions.OpUserView)
	if err != nil {
		return nil, err
	}
	if !perm {
		return nil, errs.ErrUnauthorized
	}
	if err = svc.verifyUsername(user); err != nil {
		return
	}
	if user.ProviderId == "" {
		return svc.createNewUser(user)
	}
	//existing, err := svc.db.GetUserByProviderId(user.ProviderId, false)
	switch {
	//case err == nil:
	//return svc.db.CreateUserWithRecord(user)
	case err == errs.ErrNotFound:
		return svc.createNewUser(user)
	}

	return nil, err
}

func (svc service) CreateUsers(ctx context.Context, users []*apiv1.User) (u []*apiv1.User, err error) {
	username, ok := ctx.Value("username").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}
	password, ok := ctx.Value("password").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}
	pass, err := svc.db.AuthorizeUser(username, password, permissions.OpUserCreate)
	switch {
	case err == errs.ErrNotFound:
		return nil, errs.ErrUnauthorized
	case err != nil:
		return nil, err
	case pass == false:
		return nil, errs.ErrUnauthorized
	}
	var done []*apiv1.User
	for _, u := range users {
		if err = svc.verifyUsername(u); err != nil {
			return done, err
		}
		_, err = svc.createNewUser(u)
		if err != nil {
			return done, err
		}
		done = append(done, u)
	}
	return done, nil
}

func (svc service) createUserWithId(user *apiv1.User) (u *apiv1.User, err error) {
	return svc.db.CreateUserWithId(user)
}

func (svc service) createUpstreamUser(user *apiv1.User) (u *apiv1.User, err error) {
	user.Provider = "welaw"
	u, err = svc.db.CreateUser(user)
	if err != nil {
		return nil, err
	}
	err = svc.db.CreateUserRoles(u.Username, []string{"upstream-user"})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (svc service) createNewUser(user *apiv1.User) (u *apiv1.User, err error) {
	u, err = svc.db.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (svc service) DeleteUser(ctx context.Context, username string) error {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return errs.ErrUnauthorized
	}
	admin, err := svc.db.HasPermission(uid, permissions.OpUserDelete)
	if err != nil {
		return err
	}
	u, err := svc.db.GetUserById(uid, true)
	if err != nil {
		return err
	}
	if !admin {
		if u.Username != username {
			return errs.ErrUnauthorized
		}
	}
	err = svc.db.DeleteUser(username)
	return err
}

func (svc service) GetUser(ctx context.Context, opts *apiv1.GetUserOptions) (user *apiv1.User, err error) {
	uid, ok := ctx.Value("user_id").(string)
	//var auth bool
	admin, err := svc.hasPermission(uid, permissions.OpUserView)
	if err != nil {
		return nil, err
	}

	switch {
	case opts == nil:
		return nil, errs.BadRequest("opts must be set")
	case opts.ReqType == apiv1.GetUserOptions_BY_UID:
		user, err = svc.db.GetUserById(opts.Uid, admin)
		if err != nil {
			return nil, err
		}
	case opts.ReqType == apiv1.GetUserOptions_BY_USERNAME:
		user, err = svc.db.GetUserByUsername(opts.Username, uid, admin)
		if err != nil {
			return nil, err
		}
	case opts.ReqType == apiv1.GetUserOptions_BY_CONTEXT:
		if !ok {
			return nil, errs.BadRequest("user_id not found")
		}
		user, err = svc.db.GetUserById(uid, admin)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errs.BadRequest("req_type not found: %+v", opts)
	}

	roles, err := svc.db.ListUserRoles(user.Username)
	if err != nil {
		return nil, err
	}

	if admin {
		var roleNames []string
		for _, r := range roles {
			roleNames = append(roleNames, r.Name)
		}
		user.Roles = roleNames
	}

	return user, nil
}

func (svc service) ListUsers(ctx context.Context, opts *apiv1.ListUsersOptions) (users []*apiv1.User, total int, err error) {
	uid, ok := ctx.Value("user_id").(string)

	switch {
	case opts == nil:
		return nil, 0, errs.ErrBadRequest
		// case opts.ReqType == apitv1.ListUserOptions_SEARCH
	case opts.All == true && opts.Search != "":
		if !ok {
			return nil, 0, errs.ErrUnauthorized
		}
		var permission bool
		permission, err = svc.hasPermission(uid, permissions.OpUserView)
		if err != nil {
			return
		}
		if !permission {
			return nil, 0, errs.ErrUnauthorized
		}
		users, total, err = svc.db.FilterAllUsers(opts.PageSize, opts.PageNum, true, opts.Search)
	case opts.All == true:
		if !ok {
			return nil, 0, errs.ErrUnauthorized
		}
		var permission bool
		permission, err = svc.db.HasPermission(uid, permissions.OpUserView)
		if err != nil {
			return
		}
		if !permission {
			return nil, 0, errs.ErrUnauthorized
		}

		users, total, err = svc.db.ListAllUsers(opts.PageSize, opts.PageNum, true)
	case opts.Upstream != "":
		users, total, err = svc.db.ListUpstreamUsers(opts.Upstream, opts.PageSize, opts.PageNum)
	default:
		if !ok {
			return nil, 0, errs.ErrUnauthorized
		}
		var permission bool
		permission, err = svc.db.HasPermission(uid, permissions.OpUserView)
		if err != nil {
			return
		}
		if !permission {
			return nil, 0, errs.ErrUnauthorized
		}
		users, total, err = svc.db.ListPublicUsers(opts.PageSize, opts.PageNum, true)
	}
	if err != nil {
		return
	}
	return users, total, nil
}

func (svc service) UpdateUser(ctx context.Context, username string, opts *apiv1.UpdateUserOptions) (u *apiv1.User, err error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	if opts == nil {
		return nil, fmt.Errorf("options not found")
	}

	var admin bool
	admin, err = svc.db.HasPermission(uid, permissions.OpUserUpdate)
	if err != nil {
		return
	}

	loggedInUser, err := svc.db.GetUserById(uid, true)
	if err != nil {
		return nil, err
	}

	if username != loggedInUser.Username && !admin {
		return nil, errs.ErrUnauthorized
	}

	user, err := svc.db.GetUserByUsername(username, loggedInUser.Uid, true)
	if err != nil {
		return nil, err
	}

	// TODO
	//u := &apiv1.User{Uid: _uid}
	if opts.Username != "" {
		user.Username = opts.Username
	}
	if opts.EmailPrivate {
		user.EmailPrivate = !user.EmailPrivate
	}
	if opts.FullNamePrivate {
		user.FullNamePrivate = !user.FullNamePrivate
	}
	if opts.PictureUrl != "" {
		user.PictureUrl = opts.PictureUrl
	}

	if opts.Email != "" || opts.Biography != "" || opts.PictureUrl != "" || opts.Password != "" {
		if !admin {
			return nil, errs.ErrUnauthorized
		}
	}

	if opts.Password != "" {
		err = svc.db.SetPassword(user.Uid, opts.Password)
	} else {
		u, err = svc.db.UpdateUser(user.Uid, user)
	}

	return u, err
}

func (svc service) UploadAvatar(ctx context.Context, opts *apiv1.UploadAvatarOptions) (err error) {
	uid, userIDFound := ctx.Value("user_id").(string)
	username, usernameFound := ctx.Value("username").(string)
	password, passwordFound := ctx.Value("password").(string)

	if opts == nil || opts.Image == nil || opts.Username == "" {
		return errs.ErrBadRequest
	}

	var userID string
	switch {
	case userIDFound:
		userID = uid
	case usernameFound && passwordFound:
		pass, err := svc.db.AuthorizeUser(username, password, permissions.OpUserCreate)
		switch {
		case err == errs.ErrNotFound:
			return errs.ErrUnauthorized
		case err != nil:
			return err
		case pass == false:
			return errs.ErrUnauthorized
		}
		user, err := svc.db.GetUserByUsername(opts.Username, "", false)
		if err != nil {
			return err
		}
		userID = user.Uid
	default:
		return errs.ErrUnauthorized
	}

	filename := fmt.Sprintf("%s/%s.jpg", svc.Opts.AvatarDir, userID)

	// save image
	err = svc.fs.Put(filename, "image/jpeg", opts.Image)
	if err != nil {
		return err
	}

	return nil
}

func (svc service) verifyUsername(user *apiv1.User) error {
	// TODO
	restricted := []string{
		"admin",
		"contact",
		"help",
		"lawmakers",
		"laws",
		"master",
	}
	switch {
	case strings.HasPrefix(user.Username, "bot_"):
		return fmt.Errorf("username cannot begin with 'bot_'")
	}
	for _, r := range restricted {
		if r == user.Username {
			return fmt.Errorf("that username is restricted")
		}
	}
	return nil
}
