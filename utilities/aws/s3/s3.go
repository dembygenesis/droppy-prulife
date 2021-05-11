package s3

import (
	"bytes"
	"github.com/dembygenesis/droppy-prulife/utilities/file"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"os"
)

var (
	s3Session *s3.S3
)

func init() {
	fmt.Println("Horah!")
}

func GetS3Instance() {
	accessId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	fmt.Println("=========== accessId", accessId)
	fmt.Println("=========== secretKey", secretKey)

	awsRegion := os.Getenv("AWS_REGION")

	s3Session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(
			accessId,
			secretKey,
			"",
		),
	})))
}


func UploadObjectBuffer(filename string, f *bytes.Reader, bucketPath string) error {

	fmt.Println("I am here")

	// Initialize s3 instance first - damn, don't know how to handle this because this will be a lambda fn
	GetS3Instance()

	// Read file
	/*f, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("f", f)*/

	// Upload the AWS
	_, err := s3Session.PutObject(&s3.PutObjectInput{
		Body: f,
		Bucket: aws.String(bucketPath),
		// This is the DAMN file name that your file will be named
		// Key: aws.String(filename),
		// Key: aws.String("sample.png"),
		Key: aws.String(filename),
	})

	if err != nil {
		fmt.Println("error uploading to s3", err)
		return err
	} else {
		fmt.Println("No error uploading to s3!")
	}

	return nil
}

func UploadObjectMultiPart(filename string, f *multipart.FileHeader, bucketPath string) error {
	// Initialize s3 instance first - damn, don't know how to handle this because this will be a lambda fn
	GetS3Instance()

	// Upload the AWS
	_, err := s3Session.PutObject(&s3.PutObjectInput{
		Body: bytes.NewReader(file.GetMultiPartAsBuffer(f)),
		Bucket: aws.String(bucketPath),
		ContentType: aws.String(file.GetMultiPartFileType(f)),
		Key: aws.String(filename),
		ACL: aws.String("public-read"),
	})

	if err != nil {
		fmt.Println("error uploading to s3", err)
		return err
	} else {
		fmt.Println("No error uploading to s3!")
	}

	return nil
}

func UploadObject(filename string, f *bytes.Reader, fileType string, bucketPath string) error {

	fmt.Println("I am here")

	// Initialize s3 instance first - damn, don't know how to handle this because this will be a lambda fn
	GetS3Instance()

	// Read file
	/*f, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("f", f)*/

	// Upload the AWS
	_, err := s3Session.PutObject(&s3.PutObjectInput{
		Body: f,
		Bucket: aws.String(bucketPath),
		// This is the DAMN file name that your file will be named
		// Key: aws.String(filename),
		// Key: aws.String("sample.png"),
		ContentType: aws.String(fileType),
		Key: aws.String(filename),
		ACL: aws.String("public-read"),
	})

	if err != nil {
		fmt.Println("error uploading to s3", err)
		return err
	} else {
		fmt.Println("No error uploading to s3!")
	}

	return nil
}