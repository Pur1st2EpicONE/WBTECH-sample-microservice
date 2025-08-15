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
