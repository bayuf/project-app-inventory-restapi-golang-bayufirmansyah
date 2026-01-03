package repository

import (
	"context"
	"errors"

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

func (r *SaleRepository) GetSaleDetailById(ctx context.Context, id uuid.UUID) (*dto.SaleDetailResponse, error) {
	query := `
	SELECT
	s.id, si.item_id, i.name, si.quantity, si.price, s.total_amount, s.status, s.created_at
	FROM sales s
		JOIN sale_items si ON s.id = si.sale_id
		JOIN items i ON si.item_id = i.id
	WHERE s.id = $1;`

	sale := dto.SaleDetailResponse{}
	if err := r.DB.QueryRow(ctx, query, id).Scan(
		&sale.ID, &sale.ItemID, &sale.ItemName, &sale.Quantity,
		&sale.Price, &sale.TotalAmount, &sale.Status, &sale.Created_At,
	); err != nil {
		r.Logger.Error("cant scan sale detail", zap.Error(err))
		return nil, err
	}

	return &sale, nil
}

func (r *SaleRepository) GetStaffSaleDetailById(ctx context.Context, saleId, userId uuid.UUID) (*dto.SaleDetailResponse, error) {
	query := `
	SELECT
	s.id, si.item_id, i.name, si.quantity, si.price, s.total_amount, s.status, s.created_at
	FROM sales s
		JOIN sale_items si ON s.id = si.sale_id
		JOIN items i ON si.item_id = i.id
	WHERE s.id = $1 AND s.user_id = $2; `

	sale := dto.SaleDetailResponse{}
	if err := r.DB.QueryRow(ctx, query, saleId, userId).Scan(
		&sale.ID, &sale.ItemID, &sale.ItemName, &sale.Quantity,
		&sale.Price, &sale.TotalAmount, &sale.Status, &sale.Created_At,
	); err != nil {
		r.Logger.Error("cant scan sale detail", zap.Error(err))
		return nil, err
	}

	return &sale, nil
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
		Scan(&data.ID, &data.TotalAmount, &data.Status, &data.Created_At); err != nil {
		r.Logger.Error("cant scan getSaleById", zap.Error(err))
		return nil, err
	}

	return &data, nil
}

func (r *SaleRepository) GetSales(ctx context.Context, page, limit int) (*[]dto.SaleResponse, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM sales`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query getAll repo ", zap.Error(err))
		return nil, 0, err
	}
	query := `
	SELECT
		id, total_amount, status, created_at
	FROM sales
	ORDER BY created_at ASC
	LIMIT $1 OFFSET $2`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("cant scan get all sales", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	sales := []dto.SaleResponse{}
	for rows.Next() {
		sale := dto.SaleResponse{}
		rows.Scan(&sale.ID, &sale.TotalAmount, &sale.Status, &sale.Created_At)

		sales = append(sales, sale)
	}

	return &sales, total, nil
}

func (r *SaleRepository) GetSalesByUserId(ctx context.Context, page, limit int, userId uuid.UUID) (*[]dto.SaleResponse, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM sales`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query getAll repo ", zap.Error(err))
		return nil, 0, err
	}
	query := `
	SELECT
		id, total_amount, status, created_at
	FROM sales
	WHERE user_id = $3
	ORDER BY created_at ASC
	LIMIT $1 OFFSET $2`

	rows, err := r.DB.Query(ctx, query, limit, offset, userId)
	if err != nil {
		r.Logger.Error("cant scan get all sales", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	sales := []dto.SaleResponse{}
	for rows.Next() {
		sale := dto.SaleResponse{}
		rows.Scan(&sale.ID, &sale.TotalAmount, &sale.Status, &sale.Created_At)

		sales = append(sales, sale)
	}

	if len(sales) <= 0 {
		return nil, 0, errors.New("empty sales")
	}

	return &sales, total, nil
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
