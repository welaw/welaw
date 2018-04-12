package easyrepo

import (
	"errors"
	"fmt"

	git "gopkg.in/libgit2/git2go.v25"
	//git "gopkg.in/libgit2/git2go.v26"
)

func MergeBranches(repo *git.Repository, ours, theirs string) (string, error) {
	ourBranch, err := repo.LookupBranch(ours, git.BranchLocal)
	if err != nil {
		return "", err
	}
	defer ourBranch.Free()
	theirBranch, err := repo.LookupBranch(theirs, git.BranchLocal)
	if err != nil {
		return "", err
	}
	defer theirBranch.Free()
	mergeOpts, _ := git.DefaultMergeOptions()
	//mergeOpts.FileFavor = git.MergeFileFavorNormal
	mergeOpts.TreeFlags = git.MergeTreeFailOnConflict
	ourCommit, err := repo.LookupCommit(ourBranch.Target())
	if err != nil {
		return "", err
	}
	defer ourCommit.Free()
	theirCommit, err := repo.LookupCommit(theirBranch.Target())
	if err != nil {
		return "", err
	}
	defer theirCommit.Free()
	idx, err := repo.MergeCommits(ourCommit, theirCommit, &mergeOpts)
	if err != nil {
		return "", err
	}
	defer idx.Free()
	if idx.HasConflicts() {
		return "", errors.New("merge conflicts")
	}

	treeId, err := idx.WriteTreeTo(repo)
	tree, err := repo.LookupTree(treeId)
	if err != nil {
		fmt.Println("repo.LookupTree(treeId) error")
		return "", err
	}
	defer tree.Free()
	sig := theirCommit.Author()
	c, err := repo.CreateCommit("refs/heads/master", sig, sig, fmt.Sprintf("Merged %s into %s", theirs, ours), tree, ourCommit)
	if err != nil {
		return "", err
	}
	return c.String(), nil
}
