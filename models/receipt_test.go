package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateReceiptRequest(t *testing.T) {
	// Test cases for ValidateReceiptRequest
	t.Run("should pass validation with all required fields", func(t *testing.T) {
		receipt := ReceiptRequest{
			Retailer:     "Retailer",
			PurchaseDate: "2023-01-01",
			PurchaseTime: "12:00",
			Total:        "100.00",
			Items: []ItemRequest{
				{ShortDescription: "Item 1", Price: "50.00"},
				{ShortDescription: "Item 2", Price: "50.00"},
			},
		}
		err := receipt.ValidateReceiptRequest()
		assert.NoError(t, err)
	})

	t.Run("should fail validation with missing retailer", func(t *testing.T) {
		receipt := ReceiptRequest{
			PurchaseDate: "2023-01-01",
			PurchaseTime: "12:00",
			Total:        "100.00",
			Items: []ItemRequest{
				{ShortDescription: "Item 1", Price: "50.00"},
				{ShortDescription: "Item 2", Price: "50.00"},
			},
		}
		err := receipt.ValidateReceiptRequest()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Retailer")
	})

	t.Run("should fail validation with missing purchase date", func(t *testing.T) {
		receipt := ReceiptRequest{
			Retailer:     "Retailer",
			PurchaseTime: "12:00",
			Total:        "100.00",
			Items: []ItemRequest{
				{ShortDescription: "Item 1", Price: "50.00"},
				{ShortDescription: "Item 2", Price: "50.00"},
			},
		}
		err := receipt.ValidateReceiptRequest()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PurchaseDate")
	})

	t.Run("should fail validation with missing purchase time", func(t *testing.T) {
		receipt := ReceiptRequest{
			Retailer:     "Retailer",
			PurchaseDate: "2023-01-01",
			Total:        "100.00",
			Items: []ItemRequest{
				{ShortDescription: "Item 1", Price: "50.00"},
				{ShortDescription: "Item 2", Price: "50.00"},
			},
		}
		err := receipt.ValidateReceiptRequest()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PurchaseTime")
	})

	t.Run("should fail validation with missing total", func(t *testing.T) {
		receipt := ReceiptRequest{
			Retailer:     "Retailer",
			PurchaseDate: "2023-01-01",
			PurchaseTime: "12:00",
			Items: []ItemRequest{
				{ShortDescription: "Item 1", Price: "50.00"},
				{ShortDescription: "Item 2", Price: "50.00"},
			},
		}
		err := receipt.ValidateReceiptRequest()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Total")
	})
}
