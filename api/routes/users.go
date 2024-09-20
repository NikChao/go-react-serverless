package routes

import (
	"api/providers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	user := providers.CreateUser()

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	results := providers.GetUsers(id)

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, results[0])
}
