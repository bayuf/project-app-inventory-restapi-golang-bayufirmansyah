package service

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type SaleService struct {
	Repo   *repository.SaleRepository
	Tx     db.TxManager
	Logger *zap.Logger
}

func NewSaleService(repo *repository.SaleRepository, log *zap.Logger, tx db.TxManager) *SaleService {
	return &SaleService{
		Repo:   repo,
		Tx:     tx,
		Logger: log,
	}
}

func (s *SaleService) NewSaleTX(ctx context.Context, newSale dto.NewSale, userId uuid.UUID) (*dto.SaleResponse, error) {
	tx, err := s.Tx.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// cancel all changes if error
	defer tx.Rollback(ctx)

	// init repo with TX
	repoTx := repository.NewSaleRepository(tx, s.Logger)

	// get item info
	itemInfo, err := repoTx.GetItemById(ctx, newSale.ItemID)
	if err != nil {
		return nil, err
	}

	// count total ammount
	if itemInfo.Stock <= 0 {
		return nil, errors.New("item not enough")
	}

	decQty := decimal.NewFromInt(int64(newSale.Quantity))
	currentPrice := itemInfo.Price
	Total := currentPrice.Mul(decQty)

	// Insert to sales
	newSales := dto.SalesUpdate{
		ID:          uuid.New(),
		UserID:      userId,
		TotalAmount: Total,
		Status:      "COMPLETED",
	}
	if err := repoTx.InsertNewSale(ctx, newSales); err != nil {
		return nil, err
	}

	// Insert to sales_items
	if err := repoTx.InsertSaleItem(ctx, dto.SalesItemsUpdate{
		SaleID:   newSales.ID,
		ItemID:   newSale.ItemID,
		Quantity: newSale.Quantity,
		Price:    itemInfo.Price,
		SubTotal: Total,
	}); err != nil {
		return nil, err
	}

	// update items stock
	newStock := (itemInfo.Stock - newSale.Quantity)
	if newStock <= 0 {
		return nil, errors.New("invalid quantity")
	}

	if err := repoTx.UpdateStock(ctx, dto.StockUpdateFromSale{
		ID:    newSale.ItemID,
		Stock: newStock,
	}); err != nil {
		return nil, err
	}

	// get sale data
	finalSale, err := repoTx.GetSaleById(ctx, newSales.ID)
	if err != nil {
		s.Logger.Error("cant get sale in service", zap.Error(err))
		return nil, err
	}

	// commit changes
	if err := tx.Commit(ctx); err != nil {
		s.Logger.Error("transaction failed", zap.Error(err))
		return nil, err
	}

	return &dto.SaleResponse{
		ID:          finalSale.ID,
		TotalAmount: finalSale.TotalAmount,
		Status:      finalSale.Status,
		Created_At:  finalSale.Created_At,
	}, nil
}

func (s *SaleService) GetSaleDetailById(ctx context.Context, id uuid.UUID) (*dto.SaleDetailResponse, error) {

	return s.Repo.GetSaleDetailById(ctx, id)
}
