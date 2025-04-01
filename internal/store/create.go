package store

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/database"
)

func (s *Store) Create() error {
	record := NewStoreRecord(s)
	_, err := database.InsertStore(record)
	if err != nil {
		panic(err)
	}
	fmt.Println("Creating store")
	return nil
}
