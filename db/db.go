package db

import (
	"log"

	"github.com/sharmaprinceji/delivery-management-system/internal/config"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage/sqlite"
)

func Mydb(cfg *config.Config) (storage.Storage, error) {
	st, err := sqlite.New(cfg) // Make sure sqlite.New accepts *config.Config too
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	//log.Printf("DB Path: %s\n", cfg.StoragePath)
	return st,nil
}

