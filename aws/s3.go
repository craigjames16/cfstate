package aws

import (
	"bytes"
	"fmt"
	"os"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	S3Downloader    *s3manager.Downloader
	S3Uploader      *s3manager.Uploader
	S3Service       *s3.S3
	BUCKET_NAME     string = os.Getenv("CFSTATE_BUCKET_NAME")
	STATE_FILE_NAME string = os.Getenv("CFSTATE_STATE_FILE_NAME")
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1"),
	}))

	S3Service = s3.New(sess)

	S3Uploader = s3manager.NewUploader(sess)

	S3Downloader = s3manager.NewDownloader(sess)
}

func UploadObject(stateFile []byte) (err error) {
	_, err = S3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(STATE_FILE_NAME),
		Body:   bytes.NewReader(stateFile),
	})

	return err
}

func DownloadStateFile() (output []byte, err error) {
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

func RenameObject(oldName string, newKey string) (err error) {
	_, err = S3Service.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(BUCKET_NAME),
		Key:        aws.String(newKey),
		CopySource: aws.String(fmt.Sprintf("%s/%s", BUCKET_NAME, oldName)),
	})

	return err
}
