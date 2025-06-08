package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/digitalocean/go-libvirt"
)

// TODO: 1. Delete function for store
//       2. Index for store table on database
//       3. Print function for store (kvmcli get sotre)
// 			 4. Connect stores with other instances using labels

var ErrStoreExist = errors.New("failed to insert new store record, store exist")

type Store struct {
	DB      *sql.DB         `yaml:"-"`
	Context context.Context `yaml:"_"`

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

func (store *Store) SetConnection(ctx context.Context, db *sql.DB, conn *libvirt.Libvirt) {
	_ = conn
	store.DB = db
	store.Context = ctx
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
		st.Spec.Paths.ArtifactsPath,
		st.Spec.Paths.ImagesPath,
		imageCount,
	)
}

// func (st *Store) GetBaseImagePath(name string) (string, error) {
// 	return filepath.Join(st.Spec.Paths.ArtifactsPath, st.Spec.Images[name].File)
// }

// NOTE: I can create Getters here for images, directories .. versions ..etc
// this will faciliate my operations
