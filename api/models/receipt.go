package models

type UploadReceiptRequest struct {
	FileName       string         `json:"fileName"`
	ContentLength  int64          `json:"contentLength"`
	ReceiptContext ReceiptContext `json:"receiptContext"`
}

type ReceiptContext string

const (
	BusinessReceipt ReceiptContext = "Business"
	RetailReceipt   ReceiptContext = "Retail"
)
