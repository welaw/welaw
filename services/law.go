package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/pkg/permissions"
	"github.com/welaw/welaw/proto"
)

const (
	latestLabel = "latest"
	masterLabel = "master"
)

func (svc service) CreateLaw(ctx context.Context, set *proto.LawSet, opts *proto.CreateLawOptions) (*proto.CreateLawReply, error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok || uid == "" {
		return nil, errs.ErrUnauthorized
	}
	//if perm, err := svc.hasPermission(uid, permissions.OpVersionCreate, set); err != nil {
	//return nil, err
	//} else if !perm {
	//return nil, errs.ErrUnauthorized
	//}

	user, err := svc.db.GetUserById(uid, true)
	if err != nil {
		return nil, err
	}
	var law *proto.LawSet
	switch {
	case opts == nil:
		return nil, errs.BadRequest("opts not found")
	// case opts.ReqType == proto.CreateLawOptions_VERSION
	case opts.Branch != "" && opts.Version != 0:
		law, err = svc.createVersion(ctx, user, opts.Branch, uint32(opts.Version), set)
	default:
		return nil, errs.ErrBadRequest
		//return nil, errs.BadRequest("req_type not found: %s", opts.ReqType)
	}
	if err != nil {
		return nil, err
	}
	return &proto.CreateLawReply{LawSet: law}, nil
}

func (svc service) CreateLaws(ctx context.Context, sets []*proto.LawSet, opts *proto.CreateLawsOptions) ([]*proto.LawSet, error) {
	username, usernameFound := ctx.Value("username").(string)
	password, passwordFound := ctx.Value("password").(string)
	if !usernameFound || !passwordFound {
		return nil, errs.ErrUnauthorized
	}
	if username == "" || password == "" {
		return nil, errs.ErrUnauthorized
	}

	pass, err := svc.db.AuthorizeUser(username, password, permissions.OpLawCreate)
	switch {
	case err == errs.ErrNotFound:
		return nil, errs.ErrUnauthorized
	case err != nil:
		return nil, err
	case !pass:
		return nil, errs.ErrUnauthorized
	}

	user, err := svc.db.GetUserByUsername(username, false)
	if err != nil {
		return nil, err
	}
	var done []*proto.LawSet
	for _, l := range sets {
		law, err := svc.createUpstreamLaw(ctx, user, l)
		if err != nil {
			return done, err
		}
		done = append(done, law)
	}
	return done, nil
}

func (svc service) createUpstreamLaw(ctx context.Context, user *proto.User, set *proto.LawSet) (*proto.LawSet, error) {
	if user.GetUpstream() != set.GetLaw().GetUpstream() {
		svc.logger.Log("error", "user upstream does not match law upstream",
			"user", user.GetUpstream(),
			"law", set.GetLaw().GetUpstream())
	}

	if set.Author == nil || set.Author.Username == "" {
		return nil, errs.BadRequest("law author not set: %+v", set.Author)
	}

	author, err := svc.db.GetUserByUsername(set.Author.Username, true)
	if err != nil {
		return nil, err
	}
	set.Branch.Name = author.Username
	set.Author.Email = author.Email
	set.Author.FullName = author.FullName

	master, err := svc.db.GetVersionByLatest(user.Uid, set.Law.Upstream, set.Law.Ident, masterLabel)
	// TODO check master version published_at is older than new version
	switch {
	case err == errs.ErrNotFound:
		newLaw, err := svc.createFirstVersion(ctx, set)
		if err != nil {
			return nil, err
		}
		return newLaw, nil
	case err != nil:
		return nil, err
	}

	branchHead, err := svc.db.GetVersionByLatest(user.Uid, set.Law.Upstream, set.Law.Ident, set.Branch.Name)
	switch {
	case err == errs.ErrNotFound:
		// TODO
		branchHead, err = svc.createBranch(ctx, master.Branch.Name, master.Version.Version, set)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	userLaw, err := svc.createLaw(ctx, branchHead.Version.Hash, set)
	if err != nil {
		return nil, err
	}
	// create commit on master branch
	set.Branch.Name = masterLabel
	masterHead, err := svc.db.GetVersionByLatest(user.Uid, set.Law.Upstream, set.Law.Ident, set.Branch.Name)
	if err != nil {
		return nil, err
	}
	_, err = svc.createLaw(ctx, masterHead.Version.Hash, set)
	if err != nil {
		return nil, err
	}
	return userLaw, nil
}

func (svc service) createLaw(ctx context.Context, parent string, set *proto.LawSet) (*proto.LawSet, error) {
	tx, err := svc.db.Begin()
	if err != nil {
		return nil, err
	}
	newLaw, err := svc.db.CreateVersion(tx, set)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newLaw.Version.Hash, err = svc.vc.CreateVersion(newLaw, parent)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = svc.db.UpdateVersion(tx, newLaw)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return newLaw, nil
}

func (svc service) createBranch(_ context.Context, branch string, version uint32, set *proto.LawSet) (*proto.LawSet, error) {
	tx, err := svc.db.Begin()
	if err != nil {
		return nil, err
	}
	found, err := svc.db.GetVersionByNumber("", set.Law.Upstream, set.Law.Ident, branch, version)
	if err != nil {
		return nil, err
	}
	_, err = svc.db.CreateBranch(
		tx,
		set.Law.Upstream,
		set.Law.Ident,
		set.Branch.Name,
		set.Author.Username,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = svc.vc.CreateBranch(found.Version.Hash, set)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return set, nil
}

func (svc service) createFirstVersion(ctx context.Context, set *proto.LawSet) (law *proto.LawSet, err error) {
	if set.Author == nil {
		return nil, fmt.Errorf("author is nil")
	}
	//if err = svc.verifyLaw(set.Law); err != nil {
	//return nil, err
	//}

	set.Version.Hash, err = svc.vc.CreateLaw(set)
	if err != nil {
		return nil, err
	}

	err = svc.db.CreateFirstVersion(set)
	if err != nil {
		vcErr := svc.vc.DeleteLaw(set.Law.Upstream, set.Law.Ident)
		if vcErr != nil {
			svc.logger.Log("method", "create_first_version", "error deleting law in vcs", vcErr.Error())
		}
		return nil, err
	}
	return set, nil
}

func (svc service) createVersion(ctx context.Context, user *proto.User, branch string, version uint32, set *proto.LawSet) (*proto.LawSet, error) {
	svc.logger.Log("method", "create_version", "user", user, "branch", branch, "version", version, "law", set)

	set.Author.Uid = user.Uid
	set.Author.Username = user.Username
	set.Author.FullName = user.FullName
	set.Author.Email = user.Email
	set.Author.PictureUrl = user.PictureUrl

	// add check for law first
	// get parent
	parent, err := svc.db.GetVersionByNumber(user.Uid, set.Law.Upstream, set.Law.Ident, branch, version)
	if err != nil {
		return nil, err
	}
	_, err = svc.db.GetVersionByLatest(user.Uid, set.Law.Upstream, set.Law.Ident, set.Branch.Name)
	switch {
	case err == errs.ErrNotFound:
		_, err := svc.createBranch(ctx, branch, version, set)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	tx, err := svc.db.Begin()
	if err != nil {
		return nil, err
	}
	newVersion, err := svc.db.CreateVersion(tx, set)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	hash, err := svc.vc.CreateVersion(newVersion, parent.Version.Hash)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newVersion.Version.Hash = hash
	_, err = svc.db.UpdateVersion(tx, newVersion)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return newVersion, nil
}

func (svc service) DeleteLaw(ctx context.Context, upstream, ident string, opts *proto.DeleteLawOptions) error {

	switch {
	// case opts == proto.DeleteLawOptions_BRANCH:
	// return svc.deleteBranch(upstream, ident, opts.Branch)
	}
	return nil
}

func (svc service) deleteBranch(ctx context.Context, upstream, ident, branch string) error {
	return nil
}

func (svc service) deleteVersion(ctx context.Context, upstream, ident, branch string, version uint32) error {
	return nil
}

func (svc service) DiffLaws(ctx context.Context, upstream, ident string, opts *proto.DiffLawsOptions) (r *proto.DiffLawsReply, err error) {
	switch {
	//case opts.ReqType == proto.DiffLawsOptions_VERSION:
	case opts.OurBranch != "" && opts.OurVersion != "" && opts.TheirBranch != "" && opts.TheirVersion != "":
		return svc.diffVersion(ctx, upstream, ident, opts.OurBranch, opts.TheirBranch, opts.OurVersion, opts.TheirVersion)
	default:
		return nil, errs.ErrBadRequest
	}
}

func (svc service) diffBranch(_ context.Context, upstream, ident, branch1, branch2 string) (string, error) {
	return "", nil
}

func (svc service) diffVersion(_ context.Context, upstream, ident, ourBranch, theirBranch, ourVersion, theirVersion string) (r *proto.DiffLawsReply, err error) {
	svc.logger.Log("method", "diff_version", "upstream", upstream)
	ourV, err := svc.db.GetVersion("", upstream, ident, ourBranch, ourVersion)
	if err != nil {
		return
	}
	theirV, err := svc.db.GetVersion("", upstream, ident, theirBranch, theirVersion)
	if err != nil {
		return
	}
	diff, err := svc.vc.DiffVersions(upstream, ident, ourV.Version.Hash, theirV.Version.Hash)
	if err != nil {
		return
	}
	return &proto.DiffLawsReply{Diff: diff, Theirs: theirV}, nil
}

func (svc service) GetLaw(ctx context.Context, upstream, ident string, opts *proto.GetLawOptions) (*proto.GetLawReply, error) {
	uid, _ := ctx.Value("user_id").(string)
	var branch, version string
	if opts == nil {
		branch = masterLabel
		version = latestLabel
	} else {
		branch = opts.Branch
		version = opts.Version
	}
	set, err := svc.getVersion(uid, upstream, ident, branch, version)
	if err != nil {
		return nil, err
	}
	return &proto.GetLawReply{LawSet: set}, nil
}

func (svc service) getVersion(user_id, upstream, ident, branch, version string) (v *proto.LawSet, err error) {
	v, err = svc.db.GetVersion(user_id, upstream, ident, branch, version)
	if err != nil {
		return
	}
	ec, err := svc.vc.GetVersion(upstream, ident, v.Version.Hash)
	if err != nil {
		return
	}
	v.Version.Body = string(ec.Body)
	return
}

func (svc service) ListLaws(_ context.Context, opts *proto.ListLawsOptions) (resp *proto.ListLawsReply, err error) {
	var sets []*proto.LawSet
	var total int32
	var suggestions []string

	switch {
	case opts.ReqType == proto.ListLawsOptions_BRANCH_VERSIONS:
		versions, err := svc.db.ListBranchVersions(opts.Upstream, opts.Ident, opts.Branch)
		if err != nil {
			return nil, err
		}
		for _, v := range versions {
			sets = append(sets, &proto.LawSet{Version: v})
		}
	case opts.ReqType == proto.ListLawsOptions_LAW_BRANCHES:
		sets, err = svc.db.ListLawBranches(opts.Upstream, opts.Ident)
		if err != nil {
			return nil, err
		}
	case opts.ReqType == proto.ListLawsOptions_SEARCH:
		sets, total, err = svc.db.FilterUpstreamLaws(opts.Upstream, opts.OrderBy, opts.Desc, opts.PageSize, opts.PageNum, opts.Search)
		if err != nil {
			return nil, err
		}
		if len(sets) == 0 {
			//suggestions, err = svc.db.ListSuggestions(opts.Search)
			//if err != nil {
			//return nil, err
			//}
		}
	case opts.ReqType == proto.ListLawsOptions_UPSTREAM_LAWS:
		sets, total, err = svc.db.ListUpstreamLaws(opts.Upstream, opts.OrderBy, opts.Desc, opts.PageSize, opts.PageNum)
		if err != nil {
			return nil, err
		}
	case opts.ReqType == proto.ListLawsOptions_USER_LAWS:
		sets, total, err = svc.db.ListUserLaws(opts.Username, opts.OrderBy, opts.Desc, opts.PageSize, opts.PageNum)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errs.ErrBadRequest
	}
	return &proto.ListLawsReply{LawSets: sets, Total: total, Suggestions: suggestions}, nil
}

func (svc service) UpdateLaw(ctx context.Context, l *proto.LawSet, opts *proto.UpdateLawOptions) (err error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return errs.ErrUnauthorized
	}
	if perm, err := svc.hasPermission(uid, permissions.OpLawUpdate, l); err != nil {
		return err
	} else if !perm {
		return errs.ErrUnauthorized
	}

	// TODO
	m := &proto.Law{
		Upstream:    l.Law.Upstream,
		Ident:       l.Law.Ident,
		Description: l.Law.Description,
	}
	err = svc.db.UpdateLaw(m)
	if err == sql.ErrNoRows {
		return errs.ErrNotFound
	}
	return
}

//func (svc service) verifyLaw(l *proto.Law) error {
//// check for for existing
//_, err := svc.db.GetVersionByLatest("", l.Upstream, l.Ident, "master")
//switch {
//case err == errs.ErrNotFound:
//svc.logger.Log("master branch of law not found: upstream=%v, ident=%v", l.Upstream, l.Ident)
//break
//case err != nil:
//return err
//default:
//return errs.ErrExists
//}
//// check ident
//// check repo path
//return nil
//}
