package models

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
