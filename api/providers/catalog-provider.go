package providers

import (
	"api/models"
	s3proxy "api/proxy/s3"
)

var CatalogBucket = "store-comparison-bucket-001"
var CatalogKey = "catalog.json"

func GetCatalog() models.Catalog {
	return s3proxy.GetDocument[models.Catalog](CatalogBucket, CatalogKey)
}
