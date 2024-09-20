package providers

import (
	"api/models"
	ddbproxy "api/proxy/ddb"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var groceriesTableName = "Groceries"

func GetGroceryItems(householdId string) []models.GroceryItem {
	hashKeyAttributeValues := map[string]types.AttributeValue{
		":hId": &types.AttributeValueMemberS{Value: householdId},
	}

	return ddbproxy.QueryTable[models.GroceryItem](groceriesTableName, "householdId = :hId", hashKeyAttributeValues)
}

func CreateGroceryItem(groceryItem models.GroceryItem) error {
	return ddbproxy.CreateItem(groceriesTableName, groceryItem)
}

func UpdateGroceryItem(groceryItem models.GroceryItem) error {
	key := map[string]types.AttributeValue{
		"householdId": &types.AttributeValueMemberS{Value: groceryItem.HouseholdId},
		"id":          &types.AttributeValueMemberS{Value: groceryItem.Id},
	}
	ignoreKeys := []string{"id", "householdId"}

	return ddbproxy.UpdateItem(groceriesTableName, key, groceryItem, ignoreKeys)
}

func DeleteGroceryItem(householdId string, groceryItemId string) error {
	key := map[string]types.AttributeValue{
		"householdId": &types.AttributeValueMemberS{Value: householdId},
		"id":          &types.AttributeValueMemberS{Value: groceryItemId},
	}

	return ddbproxy.DeleteItem(groceriesTableName, key)
}

func BatchDeleteGroceryItems(groceryItems []models.GroceryItem) error {
	keys := make([]map[string]types.AttributeValue, len(groceryItems))
	for index, groceryItem := range groceryItems {
		keys[index] = map[string]types.AttributeValue{
			"householdId": &types.AttributeValueMemberS{Value: groceryItem.HouseholdId},
			"id":          &types.AttributeValueMemberS{Value: groceryItem.Id},
		}
	}

	return ddbproxy.BatchDeleteItems(groceriesTableName, keys)
}
