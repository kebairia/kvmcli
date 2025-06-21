package store

import (
	"fmt"
)

func (s *Store) Delete() error {
	record := NewStoreRecord(s)
	// Delete the store record from the database.
	if err := record.Delete(s.ctx, s.db); err != nil {
		return fmt.Errorf("failed to delete record for store %q: %v", s.Config.Metadata.Name, err)
	}

	fmt.Printf("store/%s deleted\n", s.Config.Metadata.Name)
	return nil
}
