package routes

import (
	"api/providers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateHousehold(c *gin.Context) {
	household := providers.CreateHousehold()

	c.JSON(http.StatusOK, household)
}

func JoinHousehold(c *gin.Context) {
	householdId := c.Param("householdId")
	userId := c.Param("userId")

	if len(userId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userId must not be null",
		})
		return
	}

	if len(householdId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "householdId must not be null",
		})
		return
	}

	err := providers.JoinHousehold(userId, householdId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func LeaveHousehold(c *gin.Context) {
	householdId := c.Param("householdId")
	userId := c.Param("userId")

	if len(userId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "userId must not be null",
		})
		return
	}

	if len(householdId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "householdId must not be null",
		})
		return
	}

	err := providers.LeaveHousehold(userId, householdId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
