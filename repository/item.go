package repository

import (
	"context"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ItemRepository struct {
	DB     db.DBExecutor
	Logger *zap.Logger
}

func NewItemRepository(db db.DBExecutor, log *zap.Logger) *ItemRepository {
	return &ItemRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *ItemRepository) Create(ctx context.Context, new model.Item) error {
	query := `INSERT
	INTO items (id, name, category_id, rack_id, stock, min_stock, price)
	VALUES ($1, $2, $3, $4, $5, $6, $7);`

	commandTag, err := r.DB.Exec(ctx, query,
		new.ID, new.Name, new.CategoryID, new.RackID, new.Stock, new.MinStock, new.Price)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("create item failed, no row affected", zap.Error(err))
			return err
		}
		r.Logger.Error("failed create item", zap.Error(err))
		return err
	}

	r.Logger.Info("new item created with", zap.Any("ID", new.ID))
	return nil
}

func (r *ItemRepository) GetAllItems(ctx context.Context, page, limit int) (*[]model.ItemDetail, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM items WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query getAll repo ", zap.Error(err))
		return nil, 0, err
	}

	query := `
	SELECT
		i.id, i.name, i.category_id ,c.name, i.rack_id , r.code, w.name,
		w.location, i.stock, i.price ,i.created_at, i.updated_at
	FROM items i
		JOIN categories c ON c.id = i.category_id
		JOIN racks r ON r.id = i.rack_id
		JOIN warehouses w ON r.warehouse_id = w.id
	WHERE i.deleted_at IS NULL AND c.deleted_at IS NULL AND r.is_active = TRUE AND w.is_active = TRUE
		ORDER BY i.updated_at DESC
		LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("get all items error ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	items := []model.ItemDetail{}
	for rows.Next() {
		item := model.ItemDetail{}
		rows.Scan(
			&item.ID, &item.Name, &item.CategoryID, &item.CategoryName, &item.RackID, &item.RackCode,
			&item.RackName, &item.Location, &item.Stock, &item.Price, &item.Created_At, &item.Updated_At,
		)
		items = append(items, item)
	}
	return &items, total, nil
}

func (r *ItemRepository) GetItem(ctx context.Context, id uuid.UUID) (*model.ItemDetail, error) {
	query := `
	SELECT
		i.id, i.name, i.category_id ,c.name, i.rack_id , r.code, w.name,
		w.location, i.stock, i.price ,i.created_at, i.updated_at
	FROM items i
		JOIN categories c ON c.id = i.category_id
		JOIN racks r ON r.id = i.rack_id
		JOIN warehouses w ON r.warehouse_id = w.id
	WHERE i.id = $1 AND
		i.deleted_at IS NULL AND
		c.deleted_at IS NULL AND
		r.is_active = TRUE AND
		w.is_active = TRUE;`

	item := model.ItemDetail{}
	if err := r.DB.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.Name, &item.CategoryID, &item.CategoryName, &item.RackID, &item.RackCode,
		&item.RackName, &item.Location, &item.Stock, &item.Price, &item.Created_At, &item.Updated_At,
	); err != nil {
		r.Logger.Error("get item error ", zap.Error(err))
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE items
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL;`

	commandTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Info("item not found", zap.Error(err))
			return err
		}

		r.Logger.Error("error delete item ", zap.Error(err))
		return err
	}

	r.Logger.Info("item deleted ", zap.Any("ID", id))
	return nil
}
