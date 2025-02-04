package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type ReceiptRequest struct {
	Retailer     string        `json:"retailer" validate:"required"`
	PurchaseDate string        `json:"purchaseDate" validate:"required"`
	PurchaseTime string        `json:"purchaseTime" validate:"required"`
	Items        []ItemRequest `json:"items" validate:"required"`
	Total        string        `json:"total" validate:"required"`
}

type ItemRequest struct {
	ShortDescription string `json:"shortDescription" validate:"required"`
	Price            string `json:"price" validate:"required"`
}

type Receipt struct {
	Retailer         string    `json:"retailer"`
	PurchaseDateTime time.Time `json:"purchaseDate"`
	Items            []Item    `json:"items"`
	Total            float64   `json:"total"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price"`
}

func (r ReceiptRequest) ValidateReceiptRequest() error {
	validate := validator.New()
	return validate.Struct(r)
}
