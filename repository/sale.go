package repository

import (
	"context"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SaleRepository struct {
	DB     db.DBExecutor
	Logger *zap.Logger
}

func NewSaleRepository(db db.DBExecutor, log *zap.Logger) *SaleRepository {
	return &SaleRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *SaleRepository) GetItemById(ctx context.Context, id uuid.UUID) (*dto.ItemResponse, error) {
	query := `
	SELECT
		name, stock, price
	FROM items
	WHERE id = $1 AND
		deleted_at IS NULL;`

	item := dto.ItemResponse{}
	if err := r.DB.QueryRow(ctx, query, id).Scan(&item.Name, &item.Stock, &item.Price); err != nil {
		r.Logger.Error("get item error ", zap.Error(err))
		return nil, err
	}
	return &item, nil
}

func (r *SaleRepository) GetSaleById(ctx context.Context, id uuid.UUID) (*dto.SaleResponse, error) {
	query := `
	SELECT
		id, total_amount, status, created_at
	FROM sales
	WHERE id = $1;`

	data := dto.SaleResponse{}
	if err := r.DB.QueryRow(ctx, query, id).
		Scan(&data.ID, &data.TotalAmmount, &data.Status, &data.Created_At); err != nil {
		r.Logger.Error("cant scan getSaleById", zap.Error(err))
		return nil, err
	}

	return &data, nil
}

func (r *SaleRepository) UpdateStock(ctx context.Context, data dto.StockUpdateFromSale) error {
	query := `
	UPDATE items SET
		stock = $2, updated_at = NOW()
	WHERE id = $1 AND
		deleted_at IS NULL;`

	commandTag, err := r.DB.Exec(ctx, query, data.ID, data.Stock)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("update stock failed, no row affected", zap.Error(err))
			return err
		}
		r.Logger.Error("failed update stock", zap.Error(err))
		return err
	}

	r.Logger.Info("new stock updated with", zap.Any("ID", data.ID))
	return nil
}

func (r *SaleRepository) InsertNewSale(ctx context.Context, data dto.SalesUpdate) error {
	query := `
	INSERT
		INTO sales (id, user_id, total_amount, status)
	VALUES ($1, $2, $3, $4);`

	commandTag, err := r.DB.Exec(ctx, query, data.ID, data.UserID, data.TotalAmount, data.Status)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("create sale failed, no row affected", zap.Error(err))
			return err
		}
		r.Logger.Error("failed create sale", zap.Error(err))
		return err
	}

	r.Logger.Info("new sale created with", zap.Any("ID", data.ID))
	return nil
}

func (r *SaleRepository) InsertSaleItem(ctx context.Context, data dto.SalesItemsUpdate) error {
	query := `
	INSERT INTO sale_items (sale_id, item_id, quantity, price, subtotal)
	VALUES ($1, $2, $3, $4, $5);`

	commandTag, err := r.DB.Exec(ctx, query, data.SaleID, data.ItemID, data.Quantity, data.Price, data.SubTotal)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("create sale_item failed, no row affected", zap.Error(err))
			return err
		}
		r.Logger.Error("failed create sale_item", zap.Error(err))
		return err
	}

	r.Logger.Info("new sale_items created")
	return nil
}
