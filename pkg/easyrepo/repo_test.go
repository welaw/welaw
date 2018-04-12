package easyrepo

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateRepo(t *testing.T) {
	dir, err := ioutil.TempDir("", "testrepo")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	repo, err := CreateRepo(dir)
	require.NoError(t, err)
	require.NotNil(t, repo)
}

func TestDeleteRepo(t *testing.T) {
	dir, err := ioutil.TempDir("", "testrepo")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	repo, err := CreateRepo(dir)
	require.NoError(t, err)
	require.NotNil(t, repo)

	err = DeleteRepo(dir)
	require.NoError(t, err)
}

func TestOpenRepo(t *testing.T) {
	dir, err := ioutil.TempDir("", "testrepo")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	repo, err := CreateRepo(dir)
	require.NoError(t, err)
	require.NotNil(t, repo)

	r, err := OpenRepo(dir)
	require.NoError(t, err)
	require.NotNil(t, r)
}
