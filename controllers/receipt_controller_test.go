package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"receipt-processor/models"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockReceiptService struct{}

func (m *MockReceiptService) StoreReceipt(id string, receipt *models.Receipt) error {
	return nil
}

func (m *MockReceiptService) CalculatePoints(id string) (int, error) {

	if id == "nonexistent-id" {
		return 0, fmt.Errorf("receipt not found")
	}
	return 100, nil
}

func TestProcessReceipt(t *testing.T) {
	// Mock the ReceiptService
	mockService := new(MockReceiptService)
	rc := NewReceiptController(mockService)

	// Create a new Gin context
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/process", rc.ProcessReceipt)

	t.Run("should return 400 if JSON binding fails", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/process", strings.NewReader(`invalid json`))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("should return 400 if validation fails", func(t *testing.T) {
		invalidReceipt := models.ReceiptRequest{
			Retailer: "Retailer",
			// Missing required fields
		}
		body, _ := json.Marshal(invalidReceipt)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("should return 400 if date or time format is invalid", func(t *testing.T) {
		invalidReceipt := models.ReceiptRequest{
			Retailer:     "Retailer",
			PurchaseDate: "invalid-date",
			PurchaseTime: "invalid-time",
			Total:        "100.00",
			Items:        []models.ItemRequest{},
		}
		body, _ := json.Marshal(invalidReceipt)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid date or time format")
	})

	t.Run("should return 400 if total amount is invalid", func(t *testing.T) {
		invalidReceipt := models.ReceiptRequest{
			Retailer:     "Retailer",
			PurchaseDate: "2023-01-01",
			PurchaseTime: "12:00",
			Total:        "invalid-total",
			Items:        []models.ItemRequest{},
		}
		body, _ := json.Marshal(invalidReceipt)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid total amount")
	})

	t.Run("should return 400 if item price format is invalid", func(t *testing.T) {
		invalidReceipt := models.ReceiptRequest{
			Retailer:     "Retailer",
			PurchaseDate: "2023-01-01",
			PurchaseTime: "12:00",
			Total:        "100.00",
			Items: []models.ItemRequest{
				{ShortDescription: "Item 1", Price: "invalid-price"},
			},
		}
		body, _ := json.Marshal(invalidReceipt)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid price format")
	})

	t.Run("should return 200 and store receipt if request is valid", func(t *testing.T) {
		validReceipt := models.ReceiptRequest{
			Retailer:     "Retailer",
			PurchaseDate: "2023-01-01",
			PurchaseTime: "12:00",
			Total:        "100.00",
			Items: []models.ItemRequest{
				{ShortDescription: "Item 1", Price: "50.00"},
				{ShortDescription: "Item 2", Price: "50.00"},
			},
		}
		body, _ := json.Marshal(validReceipt)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "id")
	})
}
func TestGetPoints(t *testing.T) {
	// Mock the ReceiptService
	mockService := new(MockReceiptService)
	rc := NewReceiptController(mockService)

	// Create a new Gin context
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/points/:id", rc.GetPoints)

	t.Run("should return 404 if receipt ID is not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/points/nonexistent-id", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "receipt not found")
	})

	t.Run("should return 200 and points if receipt ID is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/points/valid-id", nil)
		// mockService.On("CalculatePoints", "valid-id").Return(100, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "100")
	})
}
