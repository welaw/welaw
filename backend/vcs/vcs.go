package vcs

import (
	"path/filepath"

	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes"
	"github.com/welaw/welaw/pkg/easyrepo"
	"github.com/welaw/welaw/proto"
	//git "gopkg.in/libgit2/git2go.v26"
	git "gopkg.in/libgit2/git2go.v25"
	//git "github.com/libgit2/git2go"
)

const (
	defaultFile = "law.md"
)

type VCS interface {
	CreateLaw(*proto.LawSet) (hash string, err error)
	CreateBranch(string, *proto.LawSet) (string, error)
	CreateVersion(*proto.LawSet, string) (string, error)
	DeleteBranch(upstream, ident, branch string) error
	DeleteLaw(upstream, ident string) error
	DiffVersions(upstream, ident, ours, theirs string) (string, error)
	GetBranchHead(*proto.LawSet) (*easyrepo.Commit, error)
	GetVersion(upstream, ident, hash string) (*easyrepo.Commit, error)
	//AcceptChanges(upstream, ident, branch, target string) (string, error)
	RepoPath(upstream, ident string) string
	GetPath() string
}

type vcs struct {
	path   string
	logger log.Logger
}

// TODO backend
func NewVcs(path string, logger log.Logger) (VCS, error) {
	return &vcs{
		path:   path,
		logger: logger,
	}, nil
}

func (vc *vcs) GetPath() string {
	return vc.path
}

func (vc *vcs) RepoPath(upstream, ident string) string {
	return filepath.Join(vc.path, upstream, ident)
}

func (vc *vcs) CreateLaw(set *proto.LawSet) (string, error) {
	path := vc.RepoPath(set.Law.Upstream, set.Law.Ident)
	repo, err := easyrepo.CreateRepo(path)
	if err != nil {
		return "", err
	}
	when, _ := ptypes.Timestamp(set.Version.PublishedAt)
	// create first version
	c1, err := easyrepo.CreateFirstCommit(repo, &easyrepo.Commit{
		Filename: defaultFile,
		Body:     []byte(set.Version.Body),
		Msg:      set.Version.Msg,
		Sig: &git.Signature{
			Name:  set.Author.Username,
			Email: set.Author.Email,
			When:  when,
		},
	})
	if err != nil {
		return "", err
	}

	// create user branch
	_, err = easyrepo.CreateBranch(repo, c1.String(), set.Branch.Name)
	if err != nil {
		return "", err
	}

	return c1.String(), nil
}

func (vc *vcs) CreateBranch(parent string, set *proto.LawSet) (string, error) {
	path := vc.RepoPath(set.Law.Upstream, set.Law.Ident)
	repo, err := easyrepo.OpenRepo(path)
	if err != nil {
		return "", err
	}
	_, err = easyrepo.CreateBranch(repo, parent, set.Branch.Name)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (vc *vcs) CreateVersion(set *proto.LawSet, parent string) (string, error) {
	vc.logger.Log("method", "create_version", "ident", set.Law.Ident, "version", set.Version.Version, "parent", parent)

	path := vc.RepoPath(set.Law.Upstream, set.Law.Ident)
	repo, err := easyrepo.OpenRepo(path)
	if err != nil {
		return "", err
	}

	when, _ := ptypes.Timestamp(set.Version.PublishedAt)

	c1, err := easyrepo.CreateCommit(repo, set.Branch.Name, &easyrepo.Commit{
		Filename: defaultFile,
		Body:     []byte(set.Version.Body),
		Msg:      set.Version.Msg,
		Parent:   parent,
		Sig: &git.Signature{
			Name:  set.Author.Username,
			Email: set.Author.Email,
			When:  when,
		},
	})
	if err != nil {
		return "", err
	}

	return c1.String(), nil
}

func (vc *vcs) DeleteBranch(upstream, ident, branch string) error {
	path := vc.RepoPath(upstream, ident)
	repo, err := easyrepo.OpenRepo(path)
	if err != nil {
		return err
	}
	return easyrepo.DeleteBranch(repo, branch)
}

func (vc *vcs) DeleteLaw(upstream, ident string) error {
	path := vc.RepoPath(upstream, ident)
	return easyrepo.DeleteRepo(path)
}

func (vc *vcs) DiffVersions(upstream, ident, ours, theirs string) (string, error) {
	vc.logger.Log("method", "diff_versions", "upstream", upstream, "ident", ident, "ours", ours, "theirs", theirs)

	path := vc.RepoPath(upstream, ident)
	repo, err := easyrepo.OpenRepo(path)
	if err != nil {
		return "", err
	}
	resp, err := easyrepo.DiffCommits(repo, ours, theirs)
	if err != nil {
		return "", err
	}

	return resp.Body, nil
}

func (vc *vcs) GetBranchHead(lb *proto.LawSet) (*easyrepo.Commit, error) {
	return nil, nil
}

func (vc *vcs) GetVersion(upstream, ident, hash string) (*easyrepo.Commit, error) {
	vc.logger.Log("method", "get_version",
		"upstream", upstream,
		"ident", ident,
		"hash", hash,
	)

	path := vc.RepoPath(upstream, ident)
	repo, err := easyrepo.OpenRepo(path)
	if err != nil {
		return nil, err
	}

	return easyrepo.GetEasyCommitByHash(repo, hash, defaultFile)
}

//func (vc *vcs) AcceptChanges(upstream, ident, branch, target string) (string, error) {
//vc.logger.Log("method", "accept_changes",
//"upstream", upstream,
//"ident", ident,
//"branch", branch,
//"target", target)

//path := vc.RepoPath(upstream, ident)
//repo, err := easyrepo.OpenRepo(path)
//if err != nil {
//return "", err
//}
//hash, err := easyrepo.MergeBranches(repo, branch, target)
//if err != nil {
//return "", err
//}
//return hash, nil
//}
