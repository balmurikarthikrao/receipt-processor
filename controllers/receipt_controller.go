package controllers

import (
	"fmt"
	"net/http"
	"receipt-processor/models"
	"receipt-processor/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	ReceiptControllerInterface interface {
		ProcessReceipt(c *gin.Context)
		GetPoints(c *gin.Context)
	}

	receiptController struct {
		Service services.ReceiptService
	}
)

func NewReceiptController(service services.ReceiptService) ReceiptControllerInterface {
	return &receiptController{Service: service}
}

// ProcessReceipt processes a receipt and returns a unique ID
func (rc *receiptController) ProcessReceipt(c *gin.Context) {
	var receiptRequest models.ReceiptRequest
	if err := c.ShouldBindJSON(&receiptRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := receiptRequest.ValidateReceiptRequest()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dateTime, err := parseDateTime(receiptRequest.PurchaseDate, receiptRequest.PurchaseTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date or time format"})
		return
	}

	totalFloat, err := strconv.ParseFloat(receiptRequest.Total, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid total amount"})
		return
	}

	var items []models.Item
	for _, value := range receiptRequest.Items {

		price, err := strconv.ParseFloat(value.Price, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
			return
		}

		items = append(items, models.Item{ShortDescription: value.ShortDescription, Price: price})
	}

	receipt := &models.Receipt{
		Retailer:         receiptRequest.Retailer,
		PurchaseDateTime: dateTime,
		Items:            items,
		Total:            totalFloat,
	}

	id := uuid.New().String()
	rc.Service.StoreReceipt(id, receipt)

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// GetPoints calculates and returns the points for a given receipt ID
func (rc *receiptController) GetPoints(c *gin.Context) {
	id := c.Param("id")
	points, err := rc.Service.CalculatePoints(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": points})
}

func parseDateTime(purchaseDate, purchaseTime string) (time.Time, error) {
	dateTimeStr := fmt.Sprintf("%sT%s:00Z", purchaseDate, purchaseTime)
	return time.Parse(time.RFC3339, dateTimeStr)
}
