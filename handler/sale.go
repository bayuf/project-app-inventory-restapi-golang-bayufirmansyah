package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/middleware"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

type SaleHandler struct {
	Service *service.SaleService
	Logger  *zap.Logger
	Config  *utils.Configuration
}

func NewSaleHandler(service *service.SaleService, log *zap.Logger, config *utils.Configuration) *SaleHandler {
	return &SaleHandler{
		Service: service,
		Logger:  log,
		Config:  config,
	}
}

func (h *SaleHandler) InsertNewSale(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(r)
	if !ok {
		utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if r.Method != "POST" {
		utils.ResponseFailed(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}

	Sale := dto.NewSale{}
	if err := json.NewDecoder(r.Body).Decode(&Sale); err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input format", err)
		return
	}

	// validate
	messageInvalid, err := utils.ValidateInput(&Sale)
	if err != nil {
		utils.ResponseFailed(w, http.StatusBadRequest, "invalid input data", messageInvalid)
		return
	}

	saleInfo, err := h.Service.NewSaleTX(ctx, Sale, authUser.UserID)
	if err != nil {
		h.Logger.Error("failed to create new item", zap.Error(err))
		utils.ResponseFailed(w, http.StatusBadRequest, "failed to create new item", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusCreated, "success", saleInfo)
}
