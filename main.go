package main

import (
	"receipt-processor/controllers"
	"receipt-processor/services"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// Initialize services and controllers
	receiptService := services.NewReceiptService()
	receiptController := controllers.NewReceiptController(receiptService)

	router.POST("/receipts/process", receiptController.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptController.GetPoints)

	router.Run(":8080")
}
