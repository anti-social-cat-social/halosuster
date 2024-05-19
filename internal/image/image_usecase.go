package image

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"github.com/google/uuid"
)
   
// createSession creates a new AWS session
func createSession() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
   
	return session.NewSession(&aws.Config{
	 	Region: aws.String(region),
	 	Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})
}
   
// UploadFileToS3 uploads a file to S3 with a given prefix
func UploadFileToS3(file) (string, error) {
	sess, err := createSession()
	if err != nil {
	 	return nil, err
	}
   
	svc := s3.New(sess)

	bucket := ""
	key := uuid.NewString()
   
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(file.Size),
	})
	if err != nil {
		return nil, err
	}
   
	return fmt.Sprintf("https://awss3.%s.jpeg", key), nil
}