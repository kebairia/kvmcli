package store

import (
	"fmt"

	db "github.com/kebairia/kvmcli/internal/database"
)

func (s *Store) Create() error {
	record := NewStoreRecord(s)
	err := db.InsertStore(db.Ctx, db.DB, record)
	if err != nil {
		panic(err)
	}
	fmt.Printf("store/%s created\n", s.Metadata.Name)
	return nil
}
