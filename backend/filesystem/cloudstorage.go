package filesystem

import (
	"context"
	"io/ioutil"

	"cloud.google.com/go/storage"

	"github.com/go-kit/kit/log"
)

type CloudStorage struct {
	logger log.Logger
	bucket string
	client *storage.Client
}

func NewCloudStorage(bucket string, logger log.Logger, ctx context.Context) *CloudStorage {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return &CloudStorage{
		bucket: bucket,
		client: client,
		logger: logger,
	}

}

func (fs *CloudStorage) Exists(filename string) bool {
	return false
}

func (fs *CloudStorage) Get(filename string) ([]byte, error) {
	rc, err := fs.client.Bucket(fs.bucket).Object(filename).NewReader(context.Background())
	if err != nil {
		fs.logger.Log("error reading data from bucket", err)
		return nil, err
	}
	defer rc.Close()
	d, err := ioutil.ReadAll(rc)
	if err != nil {
		fs.logger.Log("error reading data from bucket", err)
		return nil, err
	}
	return d, nil
}

func (fs *CloudStorage) Put(filename, filetype string, content []byte) error {
	w := fs.client.Bucket(fs.bucket).Object(filename).NewWriter(context.Background())
	w.ObjectAttrs = storage.ObjectAttrs{
		Bucket:      fs.bucket,
		ContentType: filetype,
		Name:        filename,
	}
	_, err := w.Write(content)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func (fs *CloudStorage) Delete(fileName string) error {
	return nil
}
