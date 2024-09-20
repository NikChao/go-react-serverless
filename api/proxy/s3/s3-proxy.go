package s3proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var svc *s3.Client
var presignClient *s3.PresignClient

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))

	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	svc = s3.NewFromConfig(cfg)
	presignClient = s3.NewPresignClient(svc)
}

// GetDocument retrieves a document from the specified S3 bucket and key, and unmarshals it into the provided struct
func GetDocument[T interface{}](bucketName string, key string) T {
	result, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Fatalf("unable to download item %q, %v", key, err)
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Fatalf("failed to read object body: %v", err)
	}

	var v T
	if err := json.Unmarshal(body, &v); err != nil {
		log.Fatalf("failed to unmarshal json: %v", err)
	}

	return v
}

func GetKeys(bucket string, prefix string) ([]string, error) {
	var keys []string
	var continuationToken *string

	for {
		input := &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			Prefix:            aws.String(prefix),
			ContinuationToken: continuationToken,
		}

		// List objects in the bucket
		result, err := svc.ListObjectsV2(context.TODO(), input)
		if err != nil {
			return keys, fmt.Errorf("unable to list items in bucket %q, %v", bucket, err)
		}

		// Append keys to the slice
		for _, item := range result.Contents {
			keys = append(keys, *item.Key)
		}

		// Check if there are more objects to list
		if *result.IsTruncated {
			continuationToken = result.NextContinuationToken
		} else {
			break
		}
	}

	return keys, nil
}

func GetDocumentFile(bucketName string, key string) ([]byte, error) {
	result, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("unable to download item %q, %v", key, err)
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %v", err)
	}

	return body, nil
}

func MoveObject(key string, fromBucket string, toBucket string) error {
	copySource := fmt.Sprintf("%s/%s", fromBucket, key)

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(toBucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(key),
	}
	_, err := svc.CopyObject(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("unable to copy object %q from bucket %q to bucket %q, %v", key, fromBucket, toBucket, err)
	}

	_, err = svc.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(fromBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to delete object %q from bucket %q, %v", key, fromBucket, err)
	}

	return nil
}

func PutObject(key string, bucket string, fileContents []byte) error {
	body := bytes.NewReader(fileContents)
	contentLength := body.Size()

	input := &s3.PutObjectInput{
		Key:           aws.String(key),
		Bucket:        aws.String(bucket),
		Body:          body,
		ContentLength: &contentLength,
	}

	_, err := svc.PutObject(context.TODO(), input)
	return err
}

func GeneratePresignedUrl(bucketName string, key string, contentLength *int64) (string, error) {
	request := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		ContentLength: contentLength,
	}

	result, err := presignClient.PresignPutObject(context.TODO(), request, setPresignedUrlExpiration)

	return result.URL, err
}

func setPresignedUrlExpiration(opts *s3.PresignOptions) {
	opts.Expires = time.Duration(5 * time.Minute)
}
