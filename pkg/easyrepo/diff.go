package easyrepo

import (
	//import git "gopkg.in/libgit2/git2go.v26"
	git "gopkg.in/libgit2/git2go.v25"
)

type DiffResponse struct {
	Body string
}

func DiffCommits(repo *git.Repository, ours, theirs string) (*DiffResponse, error) {

	ourCommit, err := GetCommit(repo, ours)
	if err != nil {
		return nil, err
	}
	ourTree, err := ourCommit.Tree()
	if err != nil {
		return nil, err
	}

	theirCommit, err := GetCommit(repo, theirs)
	if err != nil {
		return nil, err
	}
	theirTree, err := theirCommit.Tree()
	if err != nil {
		return nil, err
	}

	//tree, err := repo.LookupTree(treeId)
	//if err != nil {
	//return nil, err
	//}
	//defer tree.Free()

	callbackInvoked := false
	opts := git.DiffOptions{
		NotifyCallback: func(diffSoFar *git.Diff, delta git.DiffDelta, matchedPathSpec string) error {
			callbackInvoked = true
			return nil
		},
		OldPrefix: "x1/",
		NewPrefix: "y1/",
	}
	diff, err := repo.DiffTreeToTree(ourTree, theirTree, &opts)
	if err != nil {
		return nil, err
	}
	//if !callbackInvoked {
	//return nil, fmt.Errorf("callback not invoked")
	//}

	if diff == nil {
		return nil, nil
	}

	files := make([]string, 0)
	hunks := make([]git.DiffHunk, 0)
	lines := make([]git.DiffLine, 0)
	patches := make([]string, 0)
	err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		patch, err := diff.Patch(len(patches))
		if err != nil {
			return nil, err
		}
		defer patch.Free()
		patchStr, err := patch.String()
		if err != nil {
			return nil, err
		}
		patches = append(patches, patchStr)

		files = append(files, file.OldFile.Path)
		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			hunks = append(hunks, hunk)
			return func(line git.DiffLine) error {
				lines = append(lines, line)
				return nil
			}, nil
		}, nil
	}, git.DiffDetailLines)
	if err != nil {
		return nil, err
	}

	if len(patches) == 0 {
		return &DiffResponse{Body: ""}, nil
	}
	return &DiffResponse{Body: patches[0]}, nil
}
