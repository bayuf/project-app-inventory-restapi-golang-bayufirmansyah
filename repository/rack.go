package repository

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"go.uber.org/zap"
)

type RackRepository struct {
	DB     db.DBExecutor
	Logger *zap.Logger
}

func NewRackRepository(db db.DBExecutor, log *zap.Logger) *RackRepository {
	return &RackRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *RackRepository) Create(ctx context.Context, newData model.Rack) error {
	query := `INSERT INTO racks (warehouse_id, code, description)
	VALUES ($1, $2, $3);`

	commandPgx, err := r.DB.Exec(ctx, query, newData.WarehouseID, newData.Code, newData.Description)
	if err != nil {
		if commandPgx.RowsAffected() == 0 {
			r.Logger.Error("error create new rack", zap.Error(err))
			return errors.New("0 row affected")
		}

		r.Logger.Error("error create new rack", zap.Error(err))
		return err
	}

	r.Logger.Info("new rack created")
	return nil
}

func (r *RackRepository) GetAllRacks(ctx context.Context, page, limit int) (*[]model.Rack, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM racks WHERE is_active=TRUE`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query findAll repo ", zap.Error(err))
		return nil, 0, err
	}

	query := `SELECT r.id, r.description, r.code, w.name, w.location, r.created_at, r.updated_at
	FROM racks r
	JOIN warehouses w ON r.warehouse_id = w.id
	WHERE r.is_active = TRUE AND w.is_active = TRUE
	ORDER BY warehouse_id
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("get all rack error ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	racks := []model.Rack{}
	for rows.Next() {
		rack := model.Rack{}
		rows.Scan(&rack.ID, &rack.Description, &rack.Code, &rack.WarehouseName,
			&rack.WarehouseLocation, &rack.Created_At, &rack.Updated_At)

		racks = append(racks, rack)

	}

	return &racks, total, nil
}

func (r *RackRepository) GetRackById(ctx context.Context, id int) (*model.Rack, error) {
	query := `SELECT r.id, r.description, r.code, w.name, w.location, r.created_at, r.updated_at
	FROM racks r
	JOIN warehouses w ON r.warehouse_id = w.id
	WHERE r.is_active = TRUE AND w.is_active = TRUE AND r.id = $1`

	rack := model.Rack{}
	if err := r.DB.QueryRow(ctx, query, id).
		Scan(&rack.ID, &rack.Description, &rack.Code, &rack.WarehouseName, &rack.WarehouseLocation,
			&rack.Created_At, &rack.Updated_At); err != nil {
		r.Logger.Error("error scan id get rack by id ", zap.Error(err))
		return nil, err
	}

	return &rack, nil
}

func (r *RackRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE racks SET is_active = FALSE, updated_at = NOW()
	WHERE is_active = TRUE AND id = $1;`

	commandTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Info("rack not found", zap.Error(err))
			return err
		}

		r.Logger.Error("error delete rack ", zap.Error(err))
		return err
	}

	r.Logger.Info("rack deleted ", zap.Int("ID:", id))
	return nil
}
