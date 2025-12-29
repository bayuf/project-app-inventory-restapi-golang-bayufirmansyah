package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/handler"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/router"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"go.uber.org/zap"
)

func main() {

	// init config
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	// init Logger
	logger, err := utils.InitLogger(config.PathLogging, config.Debug)
	if err != nil {
		log.Fatal(err)
	}

	// init DB
	dBConn, err := db.Connect(logger, &config.DB)
	if err != nil {
		log.Fatal(err)
	}

	// Init layer
	repo := repository.NewRepository(dBConn, logger)
	service := service.NewService(repo, logger, config)
	handler := handler.NewHandler(service, logger)
	router := router.NewRouter(handler, service, logger)

	fmt.Println("server running on http://localhost:8080")
	if err := http.ListenAndServe(":"+config.Port, router); err != nil {
		logger.Error("failed running server", zap.Error(err))
		panic("failed running server")
	}
}
