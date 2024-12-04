package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Chandra179/go-template/configs"
	"github.com/Chandra179/go-template/internal/db"
)

func StartServer() {
	// -------------
	// Configs
	// -------------
	cfg, err := configs.LoadConfigFromEnv()
	if err != nil {
		fmt.Println("err config", err)
	}
	// -------------
	// Database
	// -------------
	dbase, err := db.NewDatabase(cfg.DbConfig)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	fmt.Println(dbase)
	//---------------
	// Http Server
	// --------------
	log.Fatal(http.ListenAndServe(":8080", nil))
}
