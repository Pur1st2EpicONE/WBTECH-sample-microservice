// Package models defines the core data structures used in the service,
// including orders along with their delivery and payment details, and individual items.
//
// These structs represent the shape of JSON payloads for orders and
// include validation tags to ensure data integrity.
package models

import "time"

// Order represents a single customer order with all related details.
type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required,alphanum,min=1,max=255"`
	TrackNumber       string    `json:"track_number" validate:"required,uppercase,alphanum,min=10,max=255"`
	Entry             string    `json:"entry" validate:"required,uppercase,alpha,len=4"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" validate:"required,lowercase,alpha,len=2"`
	InternalSignature string    `json:"internal_signature,omitempty" validate:"max=255"`
	CustomerID        string    `json:"customer_id" validate:"required,lowercase,alphanum,min=1,max=255"`
	DeliveryService   string    `json:"delivery_service" validate:"required,min=2,max=255"`
	ShardKey          string    `json:"shardkey" validate:"required,numeric,min=1,max=10"`
	SmID              int       `json:"sm_id" validate:"required,gt=0"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required,numeric,min=1,max=10"`
}

// Delivery holds the recipient and address information for an order.
type Delivery struct {
	Name    string `json:"name" validate:"required,min=2,max=255,excludesall=0123456789!@#$%^&*()_+={}[]"`
	Phone   string `json:"phone" validate:"required,e164"`
	Zip     string `json:"zip" validate:"required,numeric,len=7"`
	City    string `json:"city" validate:"required,min=2,max=100,excludesall=0123456789!@#$%^&*()_+={}[]"`
	Address string `json:"address" validate:"required,min=5,max=255,excludesall=!@#$%^&*()_+={}[]"`
	Region  string `json:"region" validate:"required,min=2,max=255,excludesall=0123456789!@#$%^&*()_+={}[]"`
	Email   string `json:"email" validate:"required,email,max=100"`
}

// Payment represents payment details for an order.
type Payment struct {
	Transaction  string  `json:"transaction" validate:"required,alphanum,min=1,max=255"`
	RequestID    string  `json:"request_id,omitempty" validate:"max=255"`
	Currency     string  `json:"currency" validate:"required,uppercase,alpha,len=3"`
	Provider     string  `json:"provider" validate:"required,min=2,max=50"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	PaymentDT    int64   `json:"payment_dt" validate:"required,gt=0"`
	Bank         string  `json:"bank,omitempty" validate:"max=50"`
	DeliveryCost float64 `json:"delivery_cost" validate:"required,gte=0"`
	GoodsTotal   float64 `json:"goods_total" validate:"required,gt=0"`
	CustomFee    float64 `json:"custom_fee,omitempty" validate:"gte=0"`
}

// Item represents a single item within an order.
type Item struct {
	ChrtID      int     `json:"chrt_id" validate:"required,gt=0"`
	TrackNumber string  `json:"track_number" validate:"required,min=10,max=255,uppercase,alphanum"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Rid         string  `json:"rid" validate:"required,min=5,max=255"`
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Sale        int     `json:"sale" validate:"gte=0,lte=100"` // does Wildberries ever offer a 100% discount, I wonder?
	Size        string  `json:"size,omitempty" validate:"max=10"`
	TotalPrice  float64 `json:"total_price" validate:"required,gte=0"` // let's assume
	NmID        int     `json:"nm_id" validate:"required,gt=0"`
	Brand       string  `json:"brand" validate:"required,min=2,max=100"`
	Status      int     `json:"status" validate:"required,oneof=100 200 202 300 400"`
}
