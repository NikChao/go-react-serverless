package models

type GroceryMagicRequest struct {
	HouseholdId     string            `json:"householdId"`
	GroceryList     GroceryList       `json:"groceryList"`
	PreferredStores []StorePreference `json:"preferredStores"`
}

type GroceryMagicResponse struct {
	GroceryList GroceryList `json:"groceryList"`
}
