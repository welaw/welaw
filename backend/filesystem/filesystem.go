package filesystem

type File interface {
	Size() int64
}

type Filesystem interface {
	Delete(string) error
	//Get(string) (io.ReadCloser, error)
	//Put(filename, filetype string, content io.Reader) error
	Get(string) ([]byte, error)
	Put(filename, filetype string, content []byte) error
	Exists(string) bool
}
