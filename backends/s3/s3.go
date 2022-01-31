package s3

import (
	"bytes"
	"fmt"
	"os"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	S3Downloader    *s3manager.Downloader
	S3Uploader      *s3manager.Uploader
	S3Service       *s3.S3
	BUCKET_NAME     string
	STATE_FILE_NAME string
)

type S3Backend struct{}

func init() {
	var (
		bExist bool
		sExist bool
	)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1"),
	}))

	S3Service = s3.New(sess)

	S3Uploader = s3manager.NewUploader(sess)

	S3Downloader = s3manager.NewDownloader(sess)

	BUCKET_NAME, bExist = os.LookupEnv("CFSTATE_BUCKET_NAME")
	STATE_FILE_NAME, sExist = os.LookupEnv("CFSTATE_STATE_FILE_NAME")

	if !bExist || !sExist {
		panic(fmt.Errorf("CFSTATE_BUCKET_NAME or CFSTATE_STATE_FILE_NAME not set"))
	}
}

func (s3b S3Backend) NewBackend() (backend S3Backend) {
	return S3Backend{}
}

func (s3b S3Backend) UpdateState(stateFile []byte) (err error) {
	now := time.Now()
	sec := now.Unix()
	newKey := fmt.Sprintf("prev_states/%s-%d.json", STATE_FILE_NAME, sec)
	stateFileName := fmt.Sprintf("%s.json", STATE_FILE_NAME)

	if _, err = S3Service.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(BUCKET_NAME),
		Key:        aws.String(newKey),
		CopySource: aws.String(fmt.Sprintf("%s/%s", BUCKET_NAME, stateFileName)),
	}); err != nil {
		return err
	}

	if _, err = S3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(stateFileName),
		Body:   bytes.NewReader(stateFile),
	}); err != nil {
		return err
	}

	return err
}

func (s3b S3Backend) GetState() (output []byte, err error) {
	// Create a file to write the S3 Object contents to.
	buff := &aws.WriteAtBuffer{}

	// Write the contents of S3 Object to the file
	_, err = S3Downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(STATE_FILE_NAME),
	})

	if err != nil {
		return output, fmt.Errorf("failed to download file, %v", err)
	}

	return buff.Bytes(), err
}
