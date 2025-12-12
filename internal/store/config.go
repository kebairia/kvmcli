package store

// TODO: 1. Delete function for store
//       2. Index for store table on database
//       3. Print function for store (kvmcli get sotre)
// 			 4. Connect stores with other instances using labels

// Root config
// type StoreConfig struct {
// 	Store StoreSpec `hcl:"store,block"`
// }

// --------------------------------------------------
// STORE
// --------------------------------------------------

type StoreConfig struct {
	Name      string            `hcl:"name,label"`
	Namespace string            `hcl:"namespace"`
	Labels    map[string]string `hcl:"labels"`

	Backend string `hcl:"backend"`

	Paths Paths `hcl:"paths,block"`

	Images []*Image `hcl:"image,block"` // image "name" { ... }
}

// --------------------------------------------------
// PATHS
// --------------------------------------------------

type Paths struct {
	Artifacts string `hcl:"artifacts"`
	Images    string `hcl:"images"`
}

// --------------------------------------------------
// IMAGE
// --------------------------------------------------

type Image struct {
	Name      string `hcl:"name,label"`
	Display   string `hcl:"display"`
	Version   string `hcl:"version"`
	OSProfile string `hcl:"os_profile"`
	File      string `hcl:"file"`
	Size      string `hcl:"size"`
	Checksum  string `hcl:"checksum"`
}

// type StoreConfig struct {
// 	APIVersion string   `yaml:"apiVersion"`
// 	Kind       string   `yaml:"kind"`
// 	Metadata   Metadata `yaml:"metadata"`
// 	Spec       Spec     `yaml:"spec"`
// }
//
// type Metadata struct {
// 	Name      string            `yaml:"name"`
// 	Namespace string            `yaml:"namespace"`
// 	Labels    map[string]string `yaml:"labels"`
// }
//
// type Spec struct {
// 	Backend string  `yaml:"backend"`
// 	Paths   Paths   `yaml:"paths"`
// 	Images  []Image `yaml:"images"`
// }
//
// type Paths struct {
// 	ArtifactsPath string `yaml:"artifacts"`
// 	ImagesPath    string `yaml:"images"`
// }
//
// type Image struct {
// 	Name      string `yaml:"name"`
// 	Version   string `yaml:"version"`
// 	OsProfile string `yaml:"osProfile"`
// 	File      string `yaml:"file"`
// 	Size      string `yaml:"size"`
// 	Checksum  string `yaml:"checksum"`
// }
