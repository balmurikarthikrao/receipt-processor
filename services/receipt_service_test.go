package services

import (
	"testing"
	"time"

	"receipt-processor/models"

	"github.com/stretchr/testify/assert"
)

func TestStoreReceipt(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		receipt   *models.Receipt
		expectErr bool
	}{
		{
			name: "valid receipt",
			id:   "valid-id",
			receipt: &models.Receipt{
				Retailer: "Retailer",
				Total:    100.00,
			},
			expectErr: false,
		},
		{
			name:      "nil receipt",
			id:        "nil-receipt",
			receipt:   nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := NewReceiptService()
			err := service.StoreReceipt(tt.id, tt.receipt)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the receipt is stored correctly
				storedReceipt, exists := service.(*receiptService).receipts[tt.id]
				assert.True(t, exists)
				assert.Equal(t, *tt.receipt, storedReceipt)
			}
		})
	}
}
func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		receipt   models.Receipt
		expectErr bool
		expected  int
	}{
		{
			name: "valid receipt with round dollar total",
			id:   "valid-id-1",
			receipt: models.Receipt{
				Retailer: "Retailer1",
				Total:    100.00,
				Items: []models.Item{
					{ShortDescription: "Item1", Price: 10.00},
					{ShortDescription: "Item2", Price: 20.00},
				},
				PurchaseDateTime: time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC),
			},
			expectErr: false,
			expected:  50 + 25 + 5 + 6 + 10 + 9, // 50 for round dollar, 25 for multiple of 0.25, 5 for items, 6 for odd day, 10 for time, 9 for retailer name
		},
		{
			name: "receipt with non-round dollar total",
			id:   "valid-id-2",
			receipt: models.Receipt{
				Retailer: "Retailer2",
				Total:    100.50,
				Items: []models.Item{
					{ShortDescription: "Item1", Price: 10.00},
					{ShortDescription: "Item2", Price: 20.00},
				},
				PurchaseDateTime: time.Date(2023, 10, 2, 15, 0, 0, 0, time.UTC),
			},
			expectErr: false,
			expected:  25 + 5 + 10 + 9, // 25 for multiple of 0.25, 5 for items, 10 for time, 9 for retailer name
		},
		{
			name: "receipt with even day",
			id:   "valid-id-3",
			receipt: models.Receipt{
				Retailer: "Retailer3",
				Total:    100.00,
				Items: []models.Item{
					{ShortDescription: "Item1", Price: 10.00},
					{ShortDescription: "Item2", Price: 20.00},
				},
				PurchaseDateTime: time.Date(2023, 10, 2, 15, 0, 0, 0, time.UTC),
			},
			expectErr: false,
			expected:  50 + 25 + 5 + 10 + 9, // 50 for round dollar, 25 for multiple of 0.25, 5 for items, 10 for time, 9 for retailer name
		},
		{
			name: "receipt with no items",
			id:   "valid-id-4",
			receipt: models.Receipt{
				Retailer:         "Retailer4",
				Total:            100.00,
				Items:            []models.Item{},
				PurchaseDateTime: time.Date(2023, 10, 1, 15, 0, 0, 0, time.UTC),
			},
			expectErr: false,
			expected:  50 + 25 + 6 + 10 + 9, // 50 for round dollar, 25 for multiple of 0.25, 6 for odd day, 10 for time, 9 for retailer name
		},
		{
			name:      "nonexistent receipt",
			id:        "nonexistent-id",
			receipt:   models.Receipt{},
			expectErr: true,
			expected:  0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := NewReceiptService()
			if tt.id != "nonexistent-id" {
				err := service.StoreReceipt(tt.id, &tt.receipt)
				assert.NoError(t, err)
			}
			points, err := service.CalculatePoints(tt.id)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, points)
			}
		})
	}
}
