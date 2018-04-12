package easyrepo

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	//git "github.com/libgit2/git2go"
	"github.com/stretchr/testify/require"
	git "gopkg.in/libgit2/git2go.v25"
	//git "gopkg.in/libgit2/git2go.v26"
	//git "github.com/libgit2/git2go"
)

func TestCreateCommit(t *testing.T) {
	dir, err := ioutil.TempDir("", "testrepo")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

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
}

func TestCreateFirstCommit(t *testing.T) {
	dir, err := ioutil.TempDir("", "testrepo")
	require.NoError(t, err)
	//defer os.RemoveAll(dir)

	repo, err := CreateRepo(dir)
	require.NoError(t, err)
	require.NotNil(t, repo)

	c, err := CreateFirstCommit(repo, &Commit{
		Filename: "test.txt",
		Body:     []byte("test body"),
		Msg:      "commit message",
		Sig: &git.Signature{
			Name:  "testuser",
			Email: "testuser@welaw.org",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestDiffCommits(t *testing.T) {
}

func TestGetCommit(t *testing.T) {
}

func TestGetEasyCommit(t *testing.T) {
}
