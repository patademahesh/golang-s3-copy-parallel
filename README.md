# golang-s3-copy-parallel
Copy folder recursively to aws S3 parallel

# Pre-requisites
  1. Go version > 1.10
  2. AWS credentials
     `\# cat .aws/credentials`
    `[default]`
    `aws_access_key_id = <AWS_ACCESS_KEY>`
    `aws_secret_access_key = <AWS_SECRET_ACCESS_KEY>`
    `region = <DEFAULT_REGION>`

  3. Required go aws modules
     \# go get "github.com/aws/aws-sdk-go/aws"
     \# go get "github.com/aws/aws-sdk-go/aws/session"
     \# go get "github.com/aws/aws-sdk-go/service/s3/s3manager"

# Build instructions
  \# go build s3copy.go
  \# cp s3copy /usr/bin/s3copy

# Usage
  \# s3copy --help
  Usage of s3copy:
    -baseDir string
          Directory to copy s3 contents to. (required)
    -bucket string
          S3 Bucket to copy contents from. (required)
    -concurrency int
          Number of concurrent connections to use. (default 200)
    -region string
          Specify bucket region
    -remoteDir string
          S3 Directory Basepath to copy contents to. (required)

# Example
  \# s3copy -region=ap-southeast-1 -bucket=bucketName -baseDir=/path/to/source/directory/ -concurrency=300 -remoteDir=/destination/path/on/s3/bucket/
