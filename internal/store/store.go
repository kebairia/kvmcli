package store

import db "github.com/kebairia/kvmcli/internal/database"

type Store struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
}

type Spec struct {
	Backend string                   `yaml:"backend"`
	Config  Config                   `yaml:"config"`
	Images  map[string]db.StoreImage `yaml:"images"`
}

type Config struct {
	ArtifactsPath string `yaml:"artifactsPath"`
	ImagesPath    string `yaml:"imagesPath"`
}

type Image struct {
	Version   string `yaml:"version"   json:"version"`
	OsProfile string `yaml:"osProfile" json:"osProfile"`
	Directory string `yaml:"directory" json:"directory"`
	File      string `yaml:"file"      json:"file"`
	Checksum  string `yaml:"checksum"  json:"checksum"`
	Size      string `yaml:"size"      json:"size"`
}
