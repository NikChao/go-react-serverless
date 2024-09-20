package receiptprocessor

import (
	"api/models"
	businessmodel "api/models/business"
	"api/providers"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/gen2brain/go-fitz"
	"github.com/otiai10/gosseract/v2"
)

func ConvertAllReceiptsIntoText(receiptContext models.ReceiptContext) (string, error) {
	var prefix string
	if receiptContext == models.BusinessReceipt {
		prefix = "business/"
	} else {
		prefix = ""
	}

	keys, err := providers.GetUnprocessedReceiptKeys(prefix)
	if err != nil {
		return "", fmt.Errorf("could not get keys in bucket: %v", err)
	}

	completeText := ""
	var wg sync.WaitGroup

	someKeys := keys[0:5]
	wg.Add(len(someKeys))

	for _, key := range someKeys {
		go func(key string) {
			text, err := ConvertUnprocessedReceiptToText(key)
			if err == nil {
				completeText += text
			}
			wg.Done()
		}(key)
	}

	wg.Wait()

	completeJsonText := parseBidFoodCatalog(completeText)

	outFile, err := os.Create("all-text.json")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	defer outFile.Close()
	outFile.WriteString(completeJsonText)

	return completeJsonText, nil
}

/**
 *
 */
func ConvertUnprocessedReceiptToText(key string) (string, error) {
	extension := filepath.Ext(key)
	data, err := providers.GetUnprocessedReceipt(key)
	if err != nil {
		return "", fmt.Errorf("could not get unprocessed receipt: %v", err)
	}

	switch extension {
	case "png":
		return ConvertUnprocessedReceiptPngToText(data)
	case "pdf":
		return ConvertUnprocessedReceiptPdfToText(data)
	default:
		return ConvertUnprocessedReceiptPdfToText(data)
	}
}

/**
 *
 */
func ConvertUnprocessedReceiptPngToText(image []byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImageFromBytes(image)
	return client.Text()
}

/**
 * Converts a receipt in an s3 bucket to text by:
 * 1. converting it to an image using go-fitz
 * 2. converting it to text using tesseract-ocr (system dependency)
 */
func ConvertUnprocessedReceiptPdfToText(pdfData []byte) (string, error) {
	doc, err := fitz.NewFromMemory(pdfData)
	if err != nil {
		return "", fmt.Errorf("could not open PDF document: %v", err)
	}
	defer doc.Close()

	if err != nil {
		return "", fmt.Errorf("could not create output file: %v", err)
	}

	text := ""

	for i := 0; i < doc.NumPage(); i++ {
		image, err := doc.ImagePNG(i, 300.0)
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}

		ocrText, _ := ConvertUnprocessedReceiptPngToText(image)

		text += ocrText
	}

	return text, nil
}

func parseBidFoodCatalog(catalogText string) string {
	lines := strings.Split(catalogText, "\n")
	var items []businessmodel.BidFoodItem
	var currentItem businessmodel.BidFoodItem
	itemLineRegex := regexp.MustCompile(`^(.*)\| .* \|`)
	sizeLineRegex := regexp.MustCompile(`^Size:\s+(\d+ X \d+GR)`)
	uomLineRegex := regexp.MustCompile(`^UOM:\s+(\w+)`)
	priceLineRegex := regexp.MustCompile(`\$(\d+\.\d+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for item name line
		itemNameMatch := itemLineRegex.FindStringSubmatch(line)
		if itemNameMatch != nil {
			if currentItem.Name != "" {
				items = append(items, currentItem)
				currentItem = businessmodel.BidFoodItem{}
			}
			currentItem.Name = strings.TrimSpace(itemNameMatch[1])
		}

		// Check for size line
		sizeMatch := sizeLineRegex.FindStringSubmatch(line)
		if sizeMatch != nil {
			currentItem.Size = sizeMatch[1]
		}

		// Check for UOM line
		uomMatch := uomLineRegex.FindStringSubmatch(line)
		if uomMatch != nil {
			currentItem.UOM = uomMatch[1]
		}

		// Check for price line
		priceMatch := priceLineRegex.FindStringSubmatch(line)
		if priceMatch != nil {
			fmt.Sscanf(priceMatch[1], "%f", &currentItem.Price)
		}
	}

	if currentItem.Name != "" {
		items = append(items, currentItem)
	}

	jsonOutput, _ := json.MarshalIndent(map[string][]businessmodel.BidFoodItem{"items": items}, "", "    ")
	return string(jsonOutput)
}
