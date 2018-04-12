package easyrepo

import (
	git "gopkg.in/libgit2/git2go.v25"
	//git "gopkg.in/libgit2/git2go.v26"
)

func CreateBranchFromCommit(repo *git.Repository, parent *git.Commit, name string) (*git.Branch, error) {
	br, err := repo.CreateBranch(name, parent, false)
	if err != nil {
		return nil, err
	}
	br.Free()
	return br, nil
}

func CreateBranch(repo *git.Repository, parentHash, name string) (*git.Branch, error) {
	// get parent commit
	parent, err := GetCommit(repo, parentHash)
	if err != nil {
		return nil, err
	}
	defer parent.Free()

	br, err := repo.CreateBranch(name, parent, false)
	if err != nil {
		return nil, err
	}
	br.Free()
	return br, nil
}

func DeleteBranch(repo *git.Repository, branch string) error {
	br, err := repo.LookupBranch(branch, git.BranchLocal)
	if err != nil {
		return err
	}
	defer br.Free()
	return br.Delete()
}

func GetHead(repo *git.Repository, branch string) (*git.Commit, error) {
	br, err := repo.LookupBranch(branch, git.BranchLocal)
	if err != nil {
		return nil, err
	}
	defer br.Free()
	commitId := br.Target()
	commit, err := repo.LookupCommit(commitId)
	if err != nil {
		return nil, err
	}
	return commit, nil
}

func GetHeadBody(repo *git.Repository, branch, filename string) (*Commit, error) {
	head, err := GetHead(repo, branch)
	if err != nil {
		return nil, err
	}
	return GetEasyCommit(repo, head, filename)
}
