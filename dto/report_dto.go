package dto

import "github.com/shopspring/decimal"

type ItemsReport struct {
	TotalItems    int `json:"total_items"`
	TotalStock    int `json:"total_stock"`
	TotalLowStock int `json:"low_stock_items"`
}

type SalesReport struct {
	TotalSales          int `json:"total_sales"`
	TotalCompletedSales int `json:"total_completed"`
	TotalCanceledSales  int `json:"total_canceled"`
	TotalOnProcess      int `json:"total_on-process"`
}

type RevenueReport struct {
	TotalRevenue   decimal.Decimal `json:"total_revenue"`
	AvgTransaction decimal.Decimal `json:"average_transaction"`
}
