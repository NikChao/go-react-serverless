package models

import "github.com/google/uuid"

type BatchDeleteGroceryItemsRequest struct {
	ItemsToDelete []GroceryItem `json:"itemsToDelete" dynamodbav:"itemsToDelete"`
}

type GroceryItem struct {
	HouseholdId   string          `json:"householdId" dynamodbav:"householdId"`
	Id            string          `json:"id" dynamodbav:"id"`
	Name          string          `json:"name" dynamodbav:"name"`
	StoreOverride StorePreference `json:"storeOverride" dynamodbav:"storeOverride"`
	Checked       bool            `json:"checked" dynamodbav:"checked"`
}

type LayoutBlockType string

const (
	Text          LayoutBlockType = "Text"
	GroceryItemId LayoutBlockType = "GroceryItemId"
)

type LayoutBlock struct {
	Value string          `json:"value" dynamodbav:"value"`
	Type  LayoutBlockType `json:"type" dynamodbav:"type"`
}

type GroceryList struct {
	Name   string
	Items  []GroceryItem `json:"items" dynamodbav:"items"`
	Layout []LayoutBlock `json:"layout" dynamodbav:"layout"`
}

// Function to generate UUID for ID field
func (item *GroceryItem) GenerateID() {
	item.Id = uuid.NewString()
}

func (item *GroceryItem) GetOrGenerateID() string {
	if item.Id == "" {
		item.GenerateID()
	}

	return item.Id
}
