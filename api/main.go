package main

import (
	"api/metrics"
	"api/models"
	receiptprocessor "api/receipt-processor"
	"api/routes"
	"api/utils"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func finit() {
	router = gin.Default()

	router.Use(CORSMiddleware())
	router.Use(LoggingMiddleware())

	// Groceries
	router.GET("/groceries/:householdId", routes.GetGroceries)
	router.PUT("/groceries", routes.CreateGroceryItem)
	router.POST("/groceries", routes.UpdateGroceryItem)
	router.DELETE("/groceries/:householdId/:id", routes.DeleteGroceryItem)
	router.POST("/groceries/batchDelete", routes.BatchDeleteGroceryItems)
	router.POST("/groceries/magic", routes.GroceryMagic)

	// Households
	router.PUT("/households", routes.CreateHousehold)
	router.POST("/households/join/:householdId/:userId", routes.JoinHousehold)
	router.POST("/households/leave/:householdId/:userId", routes.LeaveHousehold)

	// Users
	router.PUT("/users", routes.CreateUser)
	router.GET("/users/:id", routes.GetUser)

	// Catalog
	router.GET("/catalog", routes.GetCatalog)

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/receipt/upload", routes.UploadReceipt)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start).Milliseconds()
		status := c.Writer.Status()
		operationName := c.Request.Method + " " + c.FullPath()

		// Create the embedded metric
		metric := metrics.EmbeddedMetric{
			OperationName: operationName,
			StatusCode:    status,
			Latency:       latency,
			CallCount:     1,
		}
		metric.Aws.Timestamp = time.Now().Unix() * 1000
		metric.Aws.CloudWatchMetrics = []metrics.MetricNamespace{
			{
				Namespace: "TaskTote/Lambda",
				Dimensions: [][]string{
					{"OperationName"},
				},
				Metrics: []metrics.MetricDetails{
					{
						Name: "Latency",
						Unit: "Milliseconds",
					},
					{
						Name: "StatusCode",
						Unit: "None",
					},
					{
						Name: "CallCount",
						Unit: "Count",
					},
				},
			},
		}

		metricJSON, _ := json.Marshal(metric)
		log.Println(string(metricJSON))
	}
}

func Handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	// Adapt the API Gateway request to a GIN request
	httpRequest, err := http.NewRequest(strings.ToUpper(req.RequestContext.HTTP.Method), req.RawPath, strings.NewReader(req.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	for key, value := range req.Headers {
		httpRequest.Header.Set(key, value)
	}

	w := &utils.ResponseWriter{}
	router.ServeHTTP(w, httpRequest)

	return events.APIGatewayProxyResponse{
		StatusCode: w.StatusCode,
		Headers:    w.Headers,
		Body:       w.Body,
	}, nil
}

func main() {
	// _, isLambda := os.LookupEnv("LAMBDA_TASK_ROOT")

	// if isLambda {
	// 	lambda.Start(Handler)
	// } else {
	// 	router.Run(":8080")
	// }
	receiptprocessor.ConvertAllReceiptsIntoText(models.BusinessReceipt)
}
