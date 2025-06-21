package store

import (
	// "context"
	// "database/sql"
	"fmt"
	// db "github.com/kebairia/kvmcli/internal/database"
)

func (s *Store) Create() error {
	record := NewStoreRecord(s)
	err := record.Insert(s.ctx, s.db)
	if err != nil {
		return fmt.Errorf("failed to insert new store record: %w", err)
	}
	fmt.Printf("store/%s created\n", s.Config.Metadata.Name)
	return nil
}
