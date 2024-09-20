package routes

import (
	"api/models"
	"api/providers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadReceipt(c *gin.Context) {
	var req models.UploadReceiptRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	url, err := providers.GetPresignedReceiptUploadUrl(req.FileName, req.ContentLength, req.ReceiptContext)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
