package store

import (
	"fmt"
	"os"
	"text/tabwriter"

	db "github.com/kebairia/kvmcli/internal/database"
)

// TODO: 1. Delete function for store
//       2. Index for store table on database
//       3. Print function for store (kvmcli get sotre)
// 			 4. Connect stores with other instances using labels

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

func (st *Store) Header() *tabwriter.Writer {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	// Columns: store name, namespace, backend, artifacts path,
	// images path, and how many images are defined
	fmt.Fprintln(w, "NAME\tNAMESPACE\tBACKEND\tARTIFACTS_PATH\tIMAGES_PATH\tIMAGE_COUNT")
	return w
}

func (st *Store) PrintRow(w *tabwriter.Writer) {
	imageCount := len(st.Spec.Images)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\n",
		st.Metadata.Name,
		st.Metadata.Namespace,
		st.Spec.Backend,
		st.Spec.Config.ArtifactsPath,
		st.Spec.Config.ImagesPath,
		imageCount,
	)
}
