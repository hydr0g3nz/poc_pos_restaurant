package entity

import (
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// DailyRevenue represents daily revenue summary
type DailyRevenue struct {
	Date         time.Time `json:"date"`
	TotalRevenue vo.Money  `json:"total_revenue"`
}

// MonthlyRevenue represents monthly revenue summary
type MonthlyRevenue struct {
	Month        time.Time `json:"month"`
	TotalRevenue vo.Money  `json:"total_revenue"`
}
