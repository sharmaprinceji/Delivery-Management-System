package db

import (
	"github.com/sharmaprinceji/delivery-management-system/internal/config"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage/sqlite"
)

func Mydb(cfg *config.Config) (storage.Storage, error) {
	st, err := sqlite.New(cfg) // Make sure sqlite.New accepts *config.Config too
	if err != nil {
		return nil, err
	}
	return st, nil
}
