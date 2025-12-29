package main

import (
	"log"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
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
	_ = dBConn
}
