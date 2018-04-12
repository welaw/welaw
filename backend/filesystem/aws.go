package filesystem

import (
	"bytes"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3FS struct {
	cfg    *aws.Config
	region string
	bucket string
}

//var (
//defaultRegion   = "us-west-1"
//defaultFiletype = "application/octet-stream"
//imageFiletype = "image/jpg"
//)

func NewS3FS(region, bucket string) Filesystem {
	//creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, token)
	creds := credentials.NewEnvCredentials()
	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	// create bucket if doesn't exist
	_, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucket})
	if err != nil {
		if !strings.HasPrefix(err.Error(), "BucketAlreadyOwnedByYou") {
			panic(err)
		}
	}

	return &S3FS{
		cfg:    cfg,
		region: region,
		bucket: bucket,
	}
}

func (fs *S3FS) Get(filename string) ([]byte, error) {
	return nil, nil
}

func (fs *S3FS) Delete(filename string) error {
	return nil
}

func (fs *S3FS) Exists(filename string) bool {
	return false
}

func (fs *S3FS) Put(filename, filetype string, content []byte) error {
	svc := s3.New(session.New(), fs.cfg)

	//content := []byte("test")

	//params := &s3.PutObjectInput{
	//Bucket:        aws.String(bucketName),
	//Key:           aws.String(awsAccessKeyID),
	//Body:          content,
	//ContentLength: aws.Int64(size),
	//ContentType:   aws.String(fileType),
	//}
	//resp, err := svc.PutObject(params)
	//if err != nil {
	//return err
	//}

	uploader := s3manager.NewUploaderWithClient(svc)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      &fs.bucket,
		Key:         &filename,
		Body:        bytes.NewReader(content),
		ContentType: &filetype,
	})
	if err != nil {
		return err
	}

	return nil
}

func (fs *S3FS) createBucket(region, bucket string) error {
	svc := s3.New(session.New(), fs.cfg)
	_, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucket})
	if err != nil {
		return err
	}
	return nil
}

func (fs *S3FS) deleteBucket(region, bucket string) error {
	svc := s3.New(session.New(), fs.cfg)
	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{Bucket: &bucket})
	if err != nil {
		return err
	}
	return nil
}
