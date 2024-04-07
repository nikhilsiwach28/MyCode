package database

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session" // Import credentials package
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nikhilsiwach28/MyCode.git/config"
)

type BlobService interface {
	InsertObject(key string, data []byte) error
	GetObject(key string) ([]byte, error)
}

// S3 represents an S3 client.
type S3 struct {
	config *config.S3Config
}

// NewS3 creates a new S3 instance with the provided configuration.
func NewS3(config *config.S3Config) *S3 {
	return &S3{
		config: config,
	}
}

// InsertObject inserts an object into the S3 bucket.
func (s *S3) InsertObject(key string, data []byte) error {
	// Create a new AWS session with HTTP and no SSL verification

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:                        aws.String(s.config.Region),
			Endpoint:                      aws.String("http://localhost:4566"), // LocalStack endpoint
			CredentialsChainVerboseErrors: aws.Bool(true),
			S3ForcePathStyle:              aws.Bool(true), // Use path-style URLs
			DisableSSL:                    aws.Bool(true), // Disable SSL verification
		},
	}))

	// Create an S3 service client
	svc := s3.New(sess)

	if err := s.CreateBucket(s.config.BucketName); err != nil {
		fmt.Println("Error Creating Bucket", err)
	}

	// Upload input parameters
	params := &s3.PutObjectInput{
		Bucket: aws.String(s.config.BucketName), // S3 bucket name
		Key:    aws.String(key),                 // Object key (filename)
		Body:   bytes.NewReader(data),           // Object data
	}

	// Upload the object to S3
	_, err := svc.PutObject(params)
	if err != nil {
		return err
	}

	if file, err := s.GetObject(key); err != nil {
		fmt.Println("Error Getting s3", err)
	} else {
		fmt.Println("File =", file)
	}

	return nil
}

// GetObject retrieves an object from the S3 bucket.
func (s *S3) GetObject(key string) ([]byte, error) {
	// Create a new AWS session with HTTP and no SSL verification
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:           aws.String(s.config.Region),
			Endpoint:         aws.String("http://localhost:4566"), // LocalStack endpoint
			S3ForcePathStyle: aws.Bool(true),                      // Use path-style URLs
			DisableSSL:       aws.Bool(true),                      // Disable SSL verification
		},
	}))

	// Create an S3 service client
	svc := s3.New(sess)

	// Download input parameters
	params := &s3.GetObjectInput{
		Bucket: aws.String(s.config.BucketName), // S3 bucket name
		Key:    aws.String(key),                 // Object key (filename)
	}

	// Download the object from S3
	resp, err := svc.GetObject(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the object data
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// CreateBucket creates a new S3 bucket.
func (s *S3) CreateBucket(bucketName string) error {
	// Create a new AWS session with HTTP and no SSL verification
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:           aws.String(s.config.Region),
			Endpoint:         aws.String("http://localhost:4566"), // LocalStack endpoint
			S3ForcePathStyle: aws.Bool(true),                      // Use path-style URLs
			DisableSSL:       aws.Bool(true),                      // Disable SSL verification
		},
	}))

	// Create an S3 service client
	svc := s3.New(sess)

	// Create input parameters
	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName), // S3 bucket name
	}

	// Create the bucket
	_, err := svc.CreateBucket(params)
	if err != nil {
		return err
	}

	return nil
}
