package routes

import (
	"api/models"
	"api/providers"
	"net/http"

	"github.com/gin-gonic/gin"
)

var catalog models.Catalog

func init() {
	catalog = providers.GetCatalog()
}

func GetCatalog(c *gin.Context) {
	c.JSON(http.StatusOK, catalog)
}
