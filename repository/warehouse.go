package repository

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"go.uber.org/zap"
)

type WarehouseRepository struct {
	DB     db.PgxIface
	Logger *zap.Logger
}

func NewWarehousesRepository(db db.PgxIface, log *zap.Logger) *WarehouseRepository {
	return &WarehouseRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *WarehouseRepository) Create(data model.Warehouse) error {
	query := `INSERT INTO warehouses (name, location) VALUES ($1, $2);`

	_, err := r.DB.Exec(context.Background(), query, data.Name, data.Location)
	if err != nil {
		r.Logger.Error("error create new warehouses :", zap.Error(err))
	}

	return nil
}

func (r *WarehouseRepository) GetAll(page, limit int) (*[]model.Warehouse, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM warehouses WHERE is_active=TRUE`
	err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query findAll repo ", zap.Error(err))
		return nil, 0, err
	}

	query := `SELECT id, name, location, created_at, updated_at
	FROM warehouses
	WHERE is_active = TRUE
	ORDER BY id ASC
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		if rows.CommandTag().RowsAffected() == 0 {
			return nil, 0, errors.New("warehouses is empty")
		}

		r.Logger.Error("error get all warehouses ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	warehouses := []model.Warehouse{}
	for rows.Next() {
		warehouse := model.Warehouse{}
		rows.Scan(&warehouse.ID, &warehouse.Name, &warehouse.Location, &warehouse.Created_at, &warehouse.Updated_at)

		warehouses = append(warehouses, warehouse)
	}

	return &warehouses, total, nil
}

func (r *WarehouseRepository) GetById(id int) (*model.Warehouse, error) {
	query := `SELECT id, name, location, created_at, updated_at
	FROM warehouses
	WHERE id = $1 AND is_active = TRUE;`

	warehouse := model.Warehouse{}
	if err := r.DB.QueryRow(context.Background(), query, id).Scan(&warehouse.ID, &warehouse.Name,
		&warehouse.Location, &warehouse.Created_at, &warehouse.Updated_at); err != nil {
		r.Logger.Error("error get by id ", zap.Error(err))
		return nil, err
	}

	return &warehouse, nil
}
