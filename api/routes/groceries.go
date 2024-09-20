package routes

import (
	"api/models"
	"api/providers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGroceries(c *gin.Context) {
	householdId := c.Param("householdId")

	groceryItems := providers.GetGroceryItems(householdId)

	layout := make([]models.LayoutBlock, len(groceryItems))
	for i, item := range groceryItems {
		layout[i] = models.LayoutBlock{Type: models.GroceryItemId, Value: item.Id}
	}

	groceryList := models.GroceryList{
		Items:  groceryItems,
		Layout: layout,
	}

	c.IndentedJSON(http.StatusOK, groceryList)
}

func CreateGroceryItem(c *gin.Context) {
	var groceryItem models.GroceryItem

	if err := c.ShouldBindJSON(&groceryItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groceryItem.GenerateID()

	err := providers.CreateGroceryItem(groceryItem)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UpdateGroceryItem(c *gin.Context) {
	var groceryItem models.GroceryItem

	if err := c.ShouldBindJSON(&groceryItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := providers.UpdateGroceryItem(groceryItem)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteGroceryItem(c *gin.Context) {
	householdId := c.Param("householdId")
	groceryItemId := c.Param("id")
	err := providers.DeleteGroceryItem(householdId, groceryItemId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func BatchDeleteGroceryItems(c *gin.Context) {
	var request models.BatchDeleteGroceryItemsRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	providers.BatchDeleteGroceryItems(request.ItemsToDelete)

	c.JSON(http.StatusOK, gin.H{})
}
