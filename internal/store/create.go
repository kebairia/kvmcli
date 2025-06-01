package store

import (
	"fmt"

	db "github.com/kebairia/kvmcli/internal/database"
)

func (s *Store) Create() error {
	record := NewStoreRecord(s)
	err := record.Insert(db.Ctx, db.DB)
	if err != nil {
		return fmt.Errorf("failed to insert new store record: %w", err)
	}
	fmt.Printf("store/%s created\n", s.Metadata.Name)
	return nil
}
