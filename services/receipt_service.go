package services

import (
	"fmt"
	"math"
	"receipt-processor/models"
	"strings"
	"sync"
	"unicode"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type (
	ReceiptService interface {
		StoreReceipt(id string, receipt *models.Receipt) error
		CalculatePoints(id string) (int, error)
	}

	receiptService struct {
		mu       sync.Mutex
		receipts map[string]models.Receipt
	}
)

func NewReceiptService() ReceiptService {
	return &receiptService{
		receipts: make(map[string]models.Receipt),
	}
}

// StoreReceipt stores a receipt in the service
func (s *receiptService) StoreReceipt(id string, receipt *models.Receipt) error {
	if receipt == nil {
		return fmt.Errorf("receipt cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.receipts == nil {
		s.receipts = make(map[string]models.Receipt)
	}

	s.receipts[id] = *receipt
	return nil
}

// CalculatePoints calculates the points for a receipt with the given ID
func (s *receiptService) CalculatePoints(id string) (int, error) {

	s.mu.Lock()
	receipt, exists := s.receipts[id]
	s.mu.Unlock()

	if !exists {
		return 0, fmt.Errorf("receipt with id %s not found", id)
	}

	return calculatePoints(receipt), nil
}

// calculatePoints calculates the points for a given receipt
func calculatePoints(receipt models.Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if isAlphanumeric(char) {
			points++
		}
	}

	// 50 points if the total is a round dollar amount with no cents
	if receipt.Total == float64(int(receipt.Total)) {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25
	if math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// Points for item descriptions
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	// 6 points if the day in the purchase date is odd
	if receipt.PurchaseDateTime.Day()%2 != 0 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm
	hour := receipt.PurchaseDateTime.Hour()
	if hour >= 14 && hour < 16 {
		points += 10
	}

	return points
}

// isAlphanumeric checks if a character is alphanumeric
func isAlphanumeric(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char)
}
