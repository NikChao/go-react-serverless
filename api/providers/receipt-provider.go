package providers

import (
	"api/models"
	s3proxy "api/proxy/s3"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

var unprocessedReceiptBucketName = "unprocessed-receipts-001"
var processedReceiptBucketName = "processed-receipts-001"
var maxObjectSize int64 = 10 * 1024 * 1024 // 10MB

func GetPresignedReceiptUploadUrl(fileName string, contentLength int64, receiptContext models.ReceiptContext) (string, error) {
	if contentLength == 0 || contentLength > maxObjectSize {
		return "", fmt.Errorf("invalid content length")
	}

	extension := filepath.Ext(fileName)
	t := time.Now()
	randomId := uuid.NewString()
	key := fmt.Sprintf("%d-%d-%d-%s%s", t.Year(), t.Month(), t.Day(), randomId, extension)

	if receiptContext == models.BusinessReceipt {
		key = fmt.Sprintf("business/%s", key)
	}

	return s3proxy.GeneratePresignedUrl(unprocessedReceiptBucketName, key, &contentLength)
}

func GetUnprocessedReceiptKeys(prefix string) ([]string, error) {
	return s3proxy.GetKeys(unprocessedReceiptBucketName, prefix)
}

func GetUnprocessedReceipt(key string) ([]byte, error) {
	return s3proxy.GetDocumentFile(unprocessedReceiptBucketName, key)
}

func MarkReceiptAsProcessed(key string) error {
	return s3proxy.MoveObject(key, unprocessedReceiptBucketName, processedReceiptBucketName)
}
