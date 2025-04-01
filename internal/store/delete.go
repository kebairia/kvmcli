package store

import "fmt"

func (s *Store) Delete() error {
	fmt.Println("Delete a store object")
	return nil
}
