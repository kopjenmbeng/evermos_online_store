package dto

import "time"

type Chart struct {
	ChartId    string `json:"chart_id"`
	Qty        int  `json:"qty"`
	TotalPrice float64
	ProductId  string
	CreatedAt  time.Time
	CreatedBy  string
	UpdatedAt  time.Time
	UpdatedBy  string
	DeletedAt  time.Time
	DeletedBy  string
}
