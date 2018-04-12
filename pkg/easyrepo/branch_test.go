package easyrepo

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	git "gopkg.in/libgit2/git2go.v25"
	//git "gopkg.in/libgit2/git2go.v26"
	//git "github.com/libgit2/git2go"
)

const (
	filename = "testing.txt"
)

func TestCreateBranch(t *testing.T) {
	dir, err := ioutil.TempDir("", "testrepo")
	require.NoError(t, err)
	//defer os.RemoveAll(dir)

	repo, err := CreateRepo(dir)
	require.NoError(t, err)
	require.NotNil(t, repo)

	c1, err := CreateFirstCommit(repo, &Commit{
		Filename: "test.txt",
		Body:     []byte("test body\n"),
		Msg:      "commit message",
		Sig: &git.Signature{
			Name:  "testuser",
			Email: "testuser@welaw.org",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, c1)

	br, err := CreateBranch(repo, c1.String(), "test-branch")
	require.NoError(t, err)
	require.NotNil(t, br)

	//err = CheckoutBranch(repo, "test-branch", nil)
	//require.NoError(t, err)

	c2, err := CreateCommit(repo, "test-branch", &Commit{
		Parent:   c1.String(),
		Filename: "test.txt",
		Body:     []byte("changed\n"),
		Msg:      "first commit on branch",
		Sig: &git.Signature{
			Name:  "testuser",
			Email: "testuser@welaw.org",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, c2)

	//err = CheckoutBranch(repo, "test-branch", nil)
	//require.NoError(t, err)
}

func TestDeleteBranch(t *testing.T) {
}

func TestGetHead(t *testing.T) {
}

func TestGetHeadBody(t *testing.T) {
}
