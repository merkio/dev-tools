package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ListBuckets print list of buckets
func ListBuckets() {

	// Create a S3 client from just a session.
	svc := s3.New(session.New())

	input := &s3.ListBucketsInput{}

	result, err := svc.ListBuckets(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			default:
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

// CreateBucket create bucket with name
func CreateBucket(bucket string) {
	fmt.Printf("Create bucket %s\n", bucket)
	svc := s3.New(session.New())
	input := &s3.CreateBucketInput{
		Bucket:                    aws.String(bucket),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{},
	}

	result, err := svc.CreateBucket(input)
	if err != nil {
		if error, ok := err.(awserr.Error); ok {
			switch error.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, error.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, error.Error())
			default:
				fmt.Println(error.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

// DeleteBucket delete bucket with name
func DeleteBucket(bucket string) {
	fmt.Printf("Delete s3 bucket %s\n", bucket)
	svc := s3.New(session.New())
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.DeleteBucket(input)
	if err != nil {
		if error, ok := err.(awserr.Error); ok {
			switch error.Code() {
			default:
				fmt.Println(error.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

// ListObjects return slice of objects for bucket
func ListObjects(bucket string) []*s3.Object {
	fmt.Printf("List objects in the bucket %s\n", bucket)
	svc := s3.New(session.New())
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(2),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		if error, ok := err.(awserr.Error); ok {
			switch error.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, error.Error())
			default:
				fmt.Println(error.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	return result.Contents
}

// UploadObject upload file to the bucket with key and tag
func UploadObject(bucket string, filePath string, key string, tag string) {
	fmt.Printf("Upload file %s to the bucket %s\n", filePath, bucket)
	svc := s3.New(session.New())
	input := &s3.PutObjectInput{
		Body:    aws.ReadSeekCloser(strings.NewReader(filePath)),
		Bucket:  aws.String(bucket),
		Key:     aws.String(key),
		Tagging: aws.String(tag),
	}

	result, err := svc.PutObject(input)
	if err != nil {
		if error, ok := err.(awserr.Error); ok {
			switch error.Code() {
			default:
				fmt.Println(error.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

// DownloadObject put file to the path with key
func DownloadObject(bucket string, key string, to string) {
	fmt.Printf("Download file %s from bucket %s to the path %s\n", key, bucket, to)
	downloader := s3manager.NewDownloader(session.New())
	// Download the item from the bucket. If an error occurs, call exitErrorf. Otherwise, notify the user that the download succeeded.

	filePath := filepath.Join(to, key)
	file, error := os.Open(filePath)
	if error != nil {
		fmt.Printf("Unable to open file %s\n%s", filePath, error)
	}
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		fmt.Printf("Unable to download item %s\n%s", key, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

// DeleteObject remove object from the bucket with key
func DeleteObject(bucket string, key string) {
	fmt.Printf("Delete object %s from bucket %s\n", key, bucket)

	svc := s3.New(session.New())
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})

	if err != nil {
		fmt.Println(fmt.Errorf("Unable to delete object %s from bucket %s, %v", key, bucket, err))
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	fmt.Printf("Object %s successfully deleted\n", key)
}
