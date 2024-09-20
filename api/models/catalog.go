package models

type StoreItemData struct {
	StoreName   StorePreference `json:"storeName"`
	ItemName    string          `json:"itemName"`
	Price       string          `json:"price"`
	LastUpdated string          `json:"lastUpdated"`
}

type CatalogItem struct {
	Name      string          `json:"name"`
	StoreData []StoreItemData `json:"storeData"`
}

type Catalog struct {
	Data []CatalogItem `json:"data"`
}
