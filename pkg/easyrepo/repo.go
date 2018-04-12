package easyrepo

import (
	"os"

	git "gopkg.in/libgit2/git2go.v25"
	//git "gopkg.in/libgit2/git2go.v26"
)

func CreateRepo(path string) (*git.Repository, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}
	if err != nil {
		return nil, err
	}
	return git.InitRepository(path, false)
}

func DeleteRepo(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(path)
}

func OpenRepo(path string) (*git.Repository, error) {
	return git.OpenRepository(path)
}
