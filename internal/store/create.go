package store

import (
	"fmt"

	db "github.com/kebairia/kvmcli/internal/database-sql"
)

func (s *Store) Create() error {
	record := NewStoreRecord(s)
	err := db.InsertStore(db.Ctx, db.DB, record)
	if err != nil {
		panic(err)
	}
	fmt.Println("Creating store")
	return nil
}
