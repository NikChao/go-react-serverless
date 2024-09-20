package routes

import (
	"api/data"
	"api/models"
	"api/parsing"
	"api/providers"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/gin-gonic/gin"
)

func GroceryMagic(c *gin.Context) {
	var request models.GroceryMagicRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	catalog := providers.GetCatalog()

	var groceryItems []models.GroceryItem
	layoutBlockMap := make(map[models.StorePreference][]models.LayoutBlock)

	var wg sync.WaitGroup

	for _, item := range request.GroceryList.Items {
		recipeUrl, isRecipeUrl := parseUrl(item.Name)

		if isRecipeUrl {
			wg.Add(1)
			providers.DeleteGroceryItem(item.HouseholdId, item.Id)
			go func() {
				defer wg.Done()
				recipeGroceryItems, extractedLayoutBlockMap := extractAndCreateGroceryItemsFromRecipeUrl(recipeUrl, request.HouseholdId, groceryItems, request.PreferredStores)

				groceryItems = append(groceryItems, recipeGroceryItems...)

				for key, value := range extractedLayoutBlockMap {
					layoutBlockMap[key] = append(layoutBlockMap[key], value...)
				}
			}()
			continue
		}

		storePreference := getStorePreferenceForItem(item, catalog, request.PreferredStores)

		layoutBlockMap[storePreference] = append(layoutBlockMap[storePreference], models.LayoutBlock{
			Value: item.Id,
			Type:  models.GroceryItemId,
		})

		groceryItems = append(groceryItems, models.GroceryItem{
			Id:          item.Id,
			Name:        item.Name,
			HouseholdId: item.HouseholdId,
			Checked:     item.Checked,
		})
	}

	wg.Wait()
	var layout []models.LayoutBlock
	for key, blocks := range layoutBlockMap {
		layout = append(layout, models.LayoutBlock{Value: string(key), Type: models.Text})
		layout = append(layout, blocks...)
	}

	groceryList := models.GroceryList{
		Items:  groceryItems,
		Layout: layout,
	}

	response := models.GroceryMagicResponse{
		GroceryList: groceryList,
	}

	c.JSON(http.StatusOK, response)
}

func getStorePreferenceForItem(item models.GroceryItem, catalog models.Catalog, preferredStores []models.StorePreference) models.StorePreference {
	if len(item.StoreOverride) > 0 {
		return item.StoreOverride
	}

	itemName := parseItemName(item.Name)
	return getCheapestStoreForItemOrStorePreference(itemName, catalog, preferredStores)
}

func getStorePreferenceForItemName(itemName string) models.StorePreference {
	store, storeExists := data.StorePreferences[itemName]

	if storeExists {
		return store
	}

	store, storeExists = data.StorePreferences[itemName+"s"]

	if storeExists {
		return store
	}

	store, storeExists = data.StorePreferences[itemName[:len(itemName)-1]]

	if storeExists {
		return store
	}

	log.Printf("No store preference for item %s\n", itemName)
	return models.Aldi
}

func removeEmojis(s string) string {
	return strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1
		}
		return r
	}, s)
}

func parseItemName(itemName string) string {
	parsedItemName := removeEmojis(itemName)
	parsedItemName = strings.TrimSpace(parsedItemName)
	return strings.ToLower(parsedItemName)
}

func parseUrl(itemName string) (string, bool) {
	u, err := url.ParseRequestURI(itemName)

	if err != nil {
		return "", false
	}

	return u.String(), true
}

func extractAndCreateGroceryItemsFromRecipeUrl(recipeUrl string, householdId string, existingGroceryItems []models.GroceryItem, preferredStores []models.StorePreference) ([]models.GroceryItem, map[models.StorePreference][]models.LayoutBlock) {
	recipe, _ := parsing.NewFromURL(recipeUrl)
	ingredients := recipe.IngredientList().Ingredients

	var groceryItems []models.GroceryItem = make([]models.GroceryItem, len(ingredients))
	layoutBlockMap := make(map[models.StorePreference][]models.LayoutBlock)

	for i, ingredient := range ingredients {
		if isIngredientAlreadyInGroceryList(ingredient, existingGroceryItems) {
			continue
		}

		storePreference := getCheapestStoreForItemOrStorePreference(ingredient.Name, catalog, preferredStores)

		groceryItem := models.GroceryItem{
			HouseholdId:   householdId,
			Name:          ingredient.Name,
			Checked:       false,
			StoreOverride: "",
		}

		groceryItem.GenerateID()

		providers.CreateGroceryItem(groceryItem)
		groceryItems[i] = groceryItem

		layoutBlockMap[storePreference] = append(layoutBlockMap[storePreference], models.LayoutBlock{
			Value: groceryItem.Id,
			Type:  models.GroceryItemId,
		})
	}

	return groceryItems, layoutBlockMap
}

func isIngredientAlreadyInGroceryList(ingredient parsing.Ingredient, groceryItems []models.GroceryItem) bool {
	for _, existingGroceryItem := range groceryItems {
		if ingredient.Name == existingGroceryItem.Name {
			return true
		}
	}

	return false
}

func getCheapestStoreForItemOrStorePreference(itemName string, catalog models.Catalog, preferredStores []models.StorePreference) models.StorePreference {
	itemStoreData := make([]models.StoreItemData, 0)

	for _, item := range catalog.Data {
		if strings.EqualFold(item.Name, itemName) {
			itemStoreData = item.StoreData
			break
		}
	}

	if len(itemStoreData) == 0 {
		fmt.Println("No item found for " + itemName)
		return getStorePreferenceForItemName(itemName)
	}

	storePreference := models.Unknown
	minPrice := float64(99999)

	if len(itemStoreData) > 0 {
		for _, itemData := range itemStoreData {
			if len(preferredStores) > 0 && !slices.Contains(preferredStores, itemData.StoreName) {
				continue
			}

			priceAtStore, err := extractNumber(itemData.Price)
			if err == nil && priceAtStore < minPrice {
				minPrice = priceAtStore
				storePreference = itemData.StoreName
			}
		}
	}

	fmt.Println("Best store for " + itemName + " is " + string(storePreference))

	return storePreference
}

func extractNumber(s string) (float64, error) {
	// Regular expression to match the first numeric value in the string
	re := regexp.MustCompile(`[-+]?[0-9]*\.?[0-9]+`)
	match := re.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("no numeric value found")
	}
	return strconv.ParseFloat(match, 64)
}
