package store

// TODO: 1. Delete function for store
//       2. Index for store table on database
//       3. Print function for store (kvmcli get sotre)
// 			 4. Connect stores with other instances using labels

type StoreConfig struct {
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
	Backend string  `yaml:"backend"`
	Paths   Paths   `yaml:"paths"`
	Images  []Image `yaml:"images"`
}

type Paths struct {
	ArtifactsPath string `yaml:"artifacts"`
	ImagesPath    string `yaml:"images"`
}

type Image struct {
	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	OsProfile string `yaml:"osProfile"`
	File      string `yaml:"file"`
	Size      string `yaml:"size"`
	Checksum  string `yaml:"checksum"`
}
