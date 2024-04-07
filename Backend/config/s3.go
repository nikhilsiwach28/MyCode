package config

// Config holds the configuration for S3.
type S3Config struct {
	Region          string // AWS region
	AccessKeyID     string // AWS access key ID
	SecretAccessKey string // AWS secret access key
	BucketName      string // S3 bucket name
}

// NewConfig creates a new Config instance with the provided values.
func NewS3Config() *S3Config {
	return &S3Config{
		Region:          GetEnvWithDefault("REGION", "us-east-1"),
		AccessKeyID:     GetEnvWithDefault("ACCESS_KEY_ID", ""),
		SecretAccessKey: GetEnvWithDefault("SECRET_ACCESS_KEY", ""),
		BucketName:      GetEnvWithDefault("BUCKET_NAME", "my-bucket"),
	}
}
