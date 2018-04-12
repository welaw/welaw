package filesystem

import (
	"os"
	"path/filepath"
)

type localFS struct {
	Root string
}

func create(path string) (f *os.File, err error) {
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0700)

	f, err = os.Create(path)
	if err != nil {
		return
	}
	return
}

func (fs *localFS) Put(path string, body []byte) (err error) {
	var f *os.File
	exists, err := fs.Exists(path)
	if err != nil {
		return
	}
	if exists {
		f, err = os.Open(path)
	} else {
		f, err = create(path)
	}
	_, err = f.Write(body)
	return
}

func (fs *localFS) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}
