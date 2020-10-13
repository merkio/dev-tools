package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/merkio/dev-tools/config"
)

// ListBuckets print list of buckets
func ListBuckets() []string {

	// Create a S3 client from just a session.
	svc := s3.New(createSession())

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
		return []string{}
	}

	return bucketsToStrings(result.Buckets, func(b s3.Bucket) string {
		return *(b.Name)
	})
}

// CreateBucket create bucket with name
func CreateBucket(bucket string) {
	fmt.Printf("Create bucket %s\n", bucket)
	svc := s3.New(createSession())
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.CreateBucket(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, err.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, err.Error())
			default:
				fmt.Println(err.Error())
			}
		} else {
			// Print the err, cast err to awserr.Error to get the Code and
			// Message from an err.
			fmt.Println(err.Error())
		}
	}

	fmt.Println(result)
}

// DeleteBucket delete bucket with name
func DeleteBucket(bucket string, force bool) {
	fmt.Printf("Delete s3 bucket %s\n", bucket)
	svc := s3.New(createSession())

	if force {
		DeleteObjects(bucket, ListObjects(bucket, ""))
	}

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.DeleteBucket(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			default:
				fmt.Println(err.Error())
			}
		} else {
			// Print the err, cast err to awserr.Error to get the Code and
			// Message from an err.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

// ListObjects return slice of objects for bucket
func ListObjects(bucket string, filterKey string) []string {
	fmt.Printf("List objects in the bucket %s\n", bucket)
	svc := s3.New(createSession())
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(2),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Fatal(s3.ErrCodeNoSuchBucket, err.Error())
			default:
				log.Fatal(err.Error())
			}
		} else {
			// Print the err, cast err to awserr.Error to get the Code and
			// Message from an err.
			log.Fatal(err.Error())
		}
	}

	return objectsToStrings(result.Contents, filterKey, func(v s3.Object) string {
		return *(v.Key)
	})
}

// UploadObject upload file to the bucket with key and tag
func UploadObject(bucket string, filePath string, key string, tag string) {
	fmt.Printf("Upload file %s to the bucket %s\n", filePath, bucket)
	svc := s3.New(createSession())
	input := &s3.PutObjectInput{
		Body:    aws.ReadSeekCloser(strings.NewReader(filePath)),
		Bucket:  aws.String(bucket),
		Key:     aws.String(key),
		Tagging: aws.String(tag),
	}

	result, err := svc.PutObject(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			default:
				fmt.Println(err.Error())
			}
		} else {
			// Print the err, cast err to awserr.Error to get the Code and
			// Message from an err.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

// DownloadObject put file to the path with key
func DownloadObject(bucket string, key string, to string) {
	fmt.Printf("Download file %s from bucket %s to the path %s\n", key, bucket, to)

	downloader := s3manager.NewDownloader(createSession())
	// Download the item from the bucket. If an err occurs, call exitErrorf. Otherwise, notify the user that the download succeeded.

	if key == bucket {
		e := os.MkdirAll(filepath.Join(to, bucket), 0770)

		if e != nil {
			fmt.Println(e)
			return
		}
	}

	file, err := createFile(to)

	if err != nil {
		log.Println("Couldn't create file with path\n", err)
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

// DownloadBucket download all files from bucket to the path
func DownloadBucket(bucket string, to string) {
	for _, obj := range ListObjects(bucket, "") {
		DownloadObject(bucket, obj, to)
	}
}

// DeleteObject remove object from the bucket with key
func DeleteObject(bucket string, key string) {
	fmt.Printf("Delete object %s from bucket %s\n", key, bucket)

	svc := s3.New(createSession())
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

// DeleteObjects delete all objects in specific bucket
func DeleteObjects(bucket string, keys []string) {
	objects := make([]*s3.ObjectIdentifier, len(keys))

	for _, key := range keys {
		object := s3.ObjectIdentifier{
			Key: aws.String(key),
		}
		objects = append(objects, &object)
	}
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{
			Objects: objects,
			Quiet:   aws.Bool(false),
		},
	}

	svc := s3.New(createSession())
	result, err := svc.DeleteObjects(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			default:
				fmt.Println(err.Error())
			}
		} else {
			// Print the err, cast err to awserr.Error to get the Code and
			// Message from an err.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}

func createSession() *session.Session {
	cfgMap := config.Config()
	awsSecret := cfgMap.GetString("aws_secret_access_key")
	awsAccessKey := cfgMap.GetString("aws_access_key_id")
	endpointURL := cfgMap.GetString("endpoint_url")

	s, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecret, ""),
		Endpoint:    aws.String(endpointURL),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		log.Fatal(err)
	}

	return s
}

func objectsToStrings(vs []*s3.Object, filterKey string, f func(s3.Object) string) []string {
	vsm := make([]string, len(vs))
	for _, v := range vs {
		if filterKey != "" {
			if strings.Contains(*v.Key, filterKey) {
				vsm = append(vsm, f(*v))
			}
		} else {
			vsm = append(vsm, f(*v))
		}
	}
	return vsm
}

func bucketsToStrings(vs []*s3.Bucket, f func(s3.Bucket) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(*v)
	}
	return vsm
}

func createFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
