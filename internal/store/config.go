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

type Config struct {
	Name      string            `hcl:"name,label"`
	Namespace string            `hcl:"namespace"`
	Labels    map[string]string `hcl:"labels,optional"`

	Backend string `hcl:"backend,optional"`

	Paths Paths `hcl:"paths,block"`

	Images []*Image `hcl:"image,block"` // image "name" { ... }
}

// --------------------------------------------------
// PATHS
// --------------------------------------------------

type Paths struct {
	Artifacts string `hcl:"artifacts,optional"`
	Images    string `hcl:"images,optional"`
}

// --------------------------------------------------
// IMAGE
// --------------------------------------------------

type Image struct {
	Name      string `hcl:"name,label"`
	Display   string `hcl:"display,optional"`
	Version   string `hcl:"version,optional"`
	OSProfile string `hcl:"os_profile,optional"`
	File      string `hcl:"file,optional"`
	Size      string `hcl:"size,optional"`
	Checksum  string `hcl:"checksum,optional"`
}
