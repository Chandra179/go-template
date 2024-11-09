package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Chandra179/go-template/configs"
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
	_, err := configs.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	//---------------
	// Http Server
	// --------------
	log.Fatal(http.ListenAndServe(":8080", nil))
}
