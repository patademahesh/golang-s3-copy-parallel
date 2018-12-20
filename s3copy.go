/*
   Author: Mahesh Patade

   Copyright Network 18 Med. Ind. Ltd. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License.

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.

*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func s3copy(filename string, bucket string, baseDir string, remoteDir string, wg *sync.WaitGroup, ses *session.Session) {
	// fmt.Println("started Goroutine ", filename)

	uploader := s3manager.NewUploader(ses)

	uploadPath := remoteDir + strings.Split(filename, baseDir)[1]

	// Upload the file's body to S3 bucket as an object
	file, err := os.Open(filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(uploadPath),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,

		// Make file publicaly available
		ACL: aws.String("public-read"),
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	// fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
	// fmt.Printf("Goroutine %s ended\n", filename)
	wg.Done()
}

var baseDir = flag.String("baseDir", "", "Directory to copy s3 contents to. (required)")
var remoteDir = flag.String("remoteDir", "", "S3 Directory Basepath to copy contents to. (required)")
var bucket = flag.String("bucket", "", "S3 Bucket to copy contents from. (required)")
var concurrency = flag.Int("concurrency", 200, "Number of concurrent connections to use.")
var region = flag.String("region", "", "Specify bucket region")

func main() {
	start := time.Now()

	flag.Parse()
	if len(*baseDir) == 0 || len(*remoteDir) == 0 || len(*bucket) == 0 || len(*region) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(*region)},
	)

	var wg sync.WaitGroup
	count, counter := 0, 0

	err := filepath.Walk(*baseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			count++
			wg.Add(1)
			go s3copy(path, *bucket, *baseDir, *remoteDir, &wg, sess)
		}

		// wait for #concurrency uploads to complete
		if count == *concurrency {
			wg.Wait()
			counter = count + counter
			fmt.Printf("Uploaded %d files..\n", counter)
			count = 0
		}
		return nil
	})
	if err != nil {
		exitErrorf("Unable to scan directory %q , %v", baseDir, err)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("\nTotal files uploaded:", count+counter, ", Time taken: ", elapsed)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
