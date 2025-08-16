package models

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid" binding:"required"`
	TrackNumber       string    `json:"track_number" binding:"required"`
	Entry             string    `json:"entry" binding:"required"`
	Delivery          Delivery  `json:"delivery" binding:"required"`
	Payment           Payment   `json:"payment" binding:"required"`
	Items             []Item    `json:"items" binding:"required"`
	Locale            string    `json:"locale" binding:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" binding:"required"`
	DeliveryService   string    `json:"delivery_service" binding:"required"`
	ShardKey          string    `json:"shardkey" binding:"required"`
	SmID              int       `json:"sm_id" binding:"required"`
	DateCreated       time.Time `json:"date_created" binding:"required"`
	OofShard          string    `json:"oof_shard" binding:"required"`
}

type Delivery struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Zip     string `json:"zip" binding:"required"`
	City    string `json:"city" binding:"required"`
	Address string `json:"address" binding:"required"`
	Region  string `json:"region" binding:"required"`
	Email   string `json:"email" binding:"required"`
}

type Payment struct {
	Transaction  string  `json:"transaction" binding:"required"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency" binding:"required"`
	Provider     string  `json:"provider" binding:"required"`
	Amount       float64 `json:"amount" binding:"required"`
	PaymentDT    int64   `json:"payment_dt" binding:"required"`
	Bank         string  `json:"bank"`
	DeliveryCost float64 `json:"delivery_cost" binding:"required"`
	GoodsTotal   float64 `json:"goods_total" binding:"required"`
	CustomFee    float64 `json:"custom_fee"`
}

type Item struct {
	ChrtID      int     `json:"chrt_id" binding:"required"`
	TrackNumber string  `json:"track_number" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Rid         string  `json:"rid" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price" binding:"required"`
	NmID        int     `json:"nm_id" binding:"required"`
	Brand       string  `json:"brand" binding:"required"`
	Status      int     `json:"status" binding:"required"`
}
