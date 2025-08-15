package models

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
