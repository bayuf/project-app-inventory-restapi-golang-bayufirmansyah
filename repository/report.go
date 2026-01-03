package repository

import (
	"context"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"go.uber.org/zap"
)

type ReportRepository struct {
	DB     db.DBExecutor
	Logger *zap.Logger
}

func NewReportRepository(db db.DBExecutor, log *zap.Logger) *ReportRepository {
	return &ReportRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *ReportRepository) GetItemsReport(ctx context.Context) (*dto.ItemsReport, error) {
	query := `
	SELECT
    	COUNT(*),
     	SUM(stock),
      	COUNT(*) FILTER (WHERE stock <= min_stock)
    FROM items
    WHERE deleted_at IS NULL;`

	var report = dto.ItemsReport{}
	if err := r.DB.QueryRow(ctx, query).Scan(&report.TotalItems, &report.TotalStock, &report.TotalLowStock); err != nil {
		r.Logger.Error("failed to get total of all items", zap.Error(err))
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) GetSalesReport(ctx context.Context) (*dto.SalesReport, error) {
	query := `
	SELECT
    	COUNT(*),
    	COUNT(*) FILTER (WHERE status = 'COMPLETED'),
    	COUNT(*) FILTER (WHERE status = 'CANCELED'),
    	COUNT(*) FILTER (WHERE status = 'PROCESS')
    FROM sales;`

	var report = dto.SalesReport{}
	if err := r.DB.QueryRow(ctx, query).
		Scan(&report.TotalSales, &report.TotalCompletedSales, &report.TotalCanceledSales, &report.TotalOnProcess); err != nil {
		r.Logger.Error("failed to get total of all sales", zap.Error(err))
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) GetRevenueReport(ctx context.Context) (*dto.RevenueReport, error) {
	query := `
	SELECT
    	COALESCE(SUM(total_amount), 0),
     	COALESCE(AVG(total_amount), 0)
    FROM sales
    WHERE status = 'COMPLETED';`

	var report = dto.RevenueReport{}
	if err := r.DB.QueryRow(ctx, query).Scan(&report.TotalRevenue, &report.AvgTransaction); err != nil {
		r.Logger.Error("failed to get total revenue", zap.Error(err))
		return nil, err
	}
	return &report, nil
}
