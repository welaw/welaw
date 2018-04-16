package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/pkg/permissions"
	"github.com/welaw/welaw/proto"
)

func (svc service) CreateUser(ctx context.Context, user *proto.User) (u *proto.User, err error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}
	if perm, err := svc.hasPermission(uid, permissions.OpUserCreate, user); err != nil {
		return nil, err
	} else if !perm {
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

func (svc service) CreateUsers(ctx context.Context, users []*proto.User) (u []*proto.User, err error) {
	username, ok := ctx.Value("username").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}
	password, ok := ctx.Value("password").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	if perm, err := svc.db.AuthorizeUser(username, password, permissions.OpUserCreate); err != nil {
		return nil, err
	} else if !perm {
		return nil, errs.ErrUnauthorized
	}

	var done []*proto.User
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

func (svc service) createUserWithId(user *proto.User) (u *proto.User, err error) {
	return svc.db.CreateUserWithId(user)
}

func (svc service) createUpstreamUser(user *proto.User) (u *proto.User, err error) {
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

func (svc service) createNewUser(user *proto.User) (u *proto.User, err error) {
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

	u, err := svc.db.GetUserByUsername(username, false)
	if err != nil {
		return err
	}

	if perm, err := svc.hasPermission(uid, permissions.OpUserDelete, u); err != nil {
		return err
	} else if !perm {
		return errs.ErrUnauthorized
	}

	// get user by username
	// send to has perm

	//u, err := svc.db.GetUserById(uid, true)
	//if err != nil {
	//return err
	//}
	//if !admin {
	//if u.Username != username {
	//return errs.ErrUnauthorized
	//}
	//}
	// get user by username

	err = svc.db.DeleteUser(username)
	return err
}

func (svc service) GetUser(ctx context.Context, opts *proto.GetUserOptions) (user *proto.User, err error) {
	uid, ok := ctx.Value("user_id").(string)

	// TODO
	switch {
	case opts == nil:
		return nil, errs.BadRequest("opts must be set")
	case opts.ReqType == proto.GetUserOptions_BY_UID:
		user, err = svc.db.GetUserById(opts.Uid, false)
		if err != nil {
			return nil, err
		}
	case opts.ReqType == proto.GetUserOptions_BY_USERNAME:
		user, err = svc.db.GetUserByUsername(opts.Username, false)
		if err != nil {
			return nil, err
		}
	case opts.ReqType == proto.GetUserOptions_BY_CONTEXT:
		if !ok {
			return nil, errs.BadRequest("user_id not found")
		}
		user, err = svc.db.GetUserById(uid, false)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errs.BadRequest("req_type not found: %+v", opts)
	}

	perm, err := svc.hasPermission(uid, permissions.OpUserView, user)
	if err != nil {
		return nil, err
	}
	// TODO
	if perm {
		user, err = svc.db.GetUserById(user.Uid, true)
	}

	perm, err = svc.hasPermission(uid, permissions.OpRolesList, user)
	if err != nil {
		return nil, err
	}
	if perm {
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

func (svc service) ListUsers(ctx context.Context, opts *proto.ListUsersOptions) (users []*proto.User, total int, err error) {
	uid, ok := ctx.Value("user_id").(string)

	switch {
	case opts == nil:
		return nil, 0, errs.ErrBadRequest
		// case opts.ReqType == apitv1.ListUserOptions_SEARCH
	case opts.All == true && opts.Search != "":
		if !ok {
			return nil, 0, errs.ErrUnauthorized
		}
		if perm, err := svc.hasPermission(uid, permissions.OpUserList, nil); err != nil {
			return nil, 0, err
		} else if !perm {
			return nil, 0, errs.ErrUnauthorized
		}
		users, total, err = svc.db.FilterAllUsers(opts.PageSize, opts.PageNum, true, opts.Search)
	case opts.All == true:
		if !ok {
			return nil, 0, errs.ErrUnauthorized
		}
		if perm, err := svc.hasPermission(uid, permissions.OpUserList, nil); err != nil {
			return nil, 0, err
		} else if !perm {
			return nil, 0, errs.ErrUnauthorized
		}

		users, total, err = svc.db.ListAllUsers(opts.PageSize, opts.PageNum, true)
	case opts.Upstream != "":
		users, total, err = svc.db.ListUpstreamUsers(opts.Upstream, opts.PageSize, opts.PageNum)
	default:
		if !ok {
			return nil, 0, errs.ErrUnauthorized
		}
		if perm, err := svc.hasPermission(uid, permissions.OpUserList, nil); err != nil {
			return nil, 0, err
		} else if !perm {
			return nil, 0, errs.ErrUnauthorized
		}
		users, total, err = svc.db.ListPublicUsers(opts.PageSize, opts.PageNum, true)
	}
	if err != nil {
		return
	}
	return users, total, nil
}

func (svc service) UpdateUser(ctx context.Context, username string, opts *proto.UpdateUserOptions) (u *proto.User, err error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	if opts == nil {
		return nil, fmt.Errorf("options not found")
	}

	user, err := svc.db.GetUserByUsername(username, true)
	if err != nil {
		return nil, err
	}

	if perm, err := svc.hasPermission(uid, permissions.OpUserUpdate, user); err != nil {
		return nil, err
	} else if !perm {
		return nil, errs.ErrUnauthorized
	}

	//loggedInUser, err := svc.db.GetUserById(uid, perm)
	//if err != nil {
	//return nil, err
	//}

	//if username != loggedInUser.Username && !admin {
	//return nil, errs.ErrUnauthorized
	//}

	//user, err := svc.db.GetUserByUsername(username, loggedInUser.Uid, true)
	//if err != nil {
	//return nil, err
	//}

	// TODO
	//u := &proto.User{Uid: _uid}
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

	//if opts.Email != "" || opts.Biography != "" || opts.PictureUrl != "" || opts.Password != "" {
	//if !admin {
	//return nil, errs.ErrUnauthorized
	//}
	//}

	// TODO
	if opts.Password != "" {
		err = svc.db.SetPassword(user.Uid, opts.Password)
	} else {
		u, err = svc.db.UpdateUser(user.Uid, user)
	}

	return u, err
}

func (svc service) UploadAvatar(ctx context.Context, opts *proto.UploadAvatarOptions) (err error) {
	// check for bad input
	if opts == nil || opts.Image == nil || opts.Username == "" {
		return errs.ErrBadRequest
	}

	// authorize
	username, _ := ctx.Value("username").(string)
	password, _ := ctx.Value("password").(string)
	uid, _ := ctx.Value("user_id").(string)
	if username != "" && password != "" {
		perm, err := svc.db.AuthorizeUser(username, password, permissions.OpUserUpdate)
		switch {
		case err == errs.ErrNotFound:
			return errs.ErrUnauthorized
		case err != nil:
			return err
		case !perm:
			return errs.ErrUnauthorized
		}
		user, err := svc.db.GetUserByUsername(username, false)
		if err != nil {
			return err
		}
		uid = user.Uid
	}
	if uid == "" {
		return errs.ErrUnauthorized
	}

	// get the target user
	user, err := svc.db.GetUserByUsername(opts.Username, false)
	if perm, err := svc.hasPermission(uid, permissions.OpUserUpdate, user); err != nil {
		return err
	} else if !perm {
		return errs.ErrUnauthorized
	}

	// save image
	filename := fmt.Sprintf("%s/%s.jpg", svc.Opts.AvatarDir, user.Uid)
	err = svc.fs.Put(filename, "image/jpeg", opts.Image)
	return err
}

func (svc service) verifyUsername(user *proto.User) error {
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
