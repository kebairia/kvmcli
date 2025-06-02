package store

import (
	"fmt"

	db "github.com/kebairia/kvmcli/internal/database"
	log "github.com/kebairia/kvmcli/internal/logger"
)

func (s *Store) Delete() error {
	name := s.Metadata.Name
	// Delete the store record from the database.
	if err := db.Delete(db.Ctx, db.DB, name, db.StoreTable); err != nil {
		log.Errorf("failed to delete record for store %q: %v", name, err)
	}

	fmt.Printf("store/%s deleted\n", s.Metadata.Name)
	return nil
}
