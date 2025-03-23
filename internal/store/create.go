package store

import (
	"fmt"
	"time"

	"github.com/kebairia/kvmcli/internal/config"
)

// Store represents the kvmcli Store object.
type Store struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// Metadata holds basic object metadata.
type Metadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
}

// Spec holds the store specification.
type Spec struct {
	Images []Image `yaml:"images"`
}

// Image represents a VM image with its metadata.
type Image struct {
	Name         string    `yaml:"name"`
	CreationDate time.Time `yaml:"creationDate"`
	Hash         string    `yaml:"hash"`
	Path         string    `yaml:"path"`
	Description  string    `yaml:"description"`
}

func (s Store) Create() error {
	s, err := config.LoadConfig(path)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	fmt.Println(s)
	return nil
}

func (s Store) Delete() error {
	return nil
}
