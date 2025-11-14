package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// TODO: 1. Delete function for store
//       2. Index for store table on database
//       3. Print function for store (kvmcli get sotre)
// 			 4. Connect stores with other instances using labels

// Errors returned by this package.
var (
	// ErrStoreNameEmpty is returned when the store name is not set.
	ErrStoreNameEmpty = errors.New("store: name is empty")
	// ErrStoreExist is returned when attempting to insert a store that already exists.
	ErrStoreExist = errors.New("store: store already exists")
	// ErrNilDBConn is returned when the database connection is missing.
	ErrNilDBConn = errors.New("store: database connection is nil")
)

// Store manages store records in a SQL database.
type Store struct {
	Config StoreConfig
	ctx    context.Context
	db     *sql.DB
}

// StoreOption configures a Store.
type StoreOption func(*Store)

// WithDatabaseConnection injects a *sql.DB (required).
func WithDatabaseConnection(db *sql.DB) StoreOption {
	return func(s *Store) {
		s.db = db
	}
}

// WithContext injects a context.Context; if nil, uses context.Background().
func WithContext(ctx context.Context) StoreOption {
	return func(s *Store) {
		if ctx == nil {
			s.ctx = context.Background()
		} else {
			s.ctx = ctx
		}
	}
}

// NewStore constructs a Store, applies options, and validates dependencies.
func NewStore(cfg StoreConfig, opts ...StoreOption) (*Store, error) {
	if cfg.Metadata.Name == "" {
		return nil, ErrStoreNameEmpty
	}

	s := &Store{
		Config: cfg,
		// default context
		ctx: context.Background(),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.db == nil {
		return nil, ErrNilDBConn
	}

	return s, nil
}

// NOTE: this is just to change later
func (st *Store) Start() error {
	fmt.Println("Start store")
	return nil
}

// func (st *Store) Header() *tabwriter.Writer {
// 	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
// 	// Columns: store name, namespace, backend, artifacts path,
// 	// images path, and how many images are defined
// 	fmt.Fprintln(w, "NAME\tNAMESPACE\tBACKEND\tARTIFACTS_PATH\tIMAGES_PATH\tIMAGE_COUNT")
// 	return w
// }
//
// func (st *Store) PrintInfo(w *tabwriter.Writer) {
// 	imageCount := len(st.Spec.Images)
// 	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\n",
// 		st.Metadata.Name,
// 		st.Metadata.Namespace,
// 		st.Spec.Backend,
// 		st.Spec.Paths.ArtifactsPath,
// 		st.Spec.Paths.ImagesPath,
// 		imageCount,
// 	)
// }
//
// // func (st *Store) GetBaseImagePath(name string) (string, error) {
// // 	return filepath.Join(st.Spec.Paths.ArtifactsPath, st.Spec.Images[name].File)
// // }
//
// // NOTE: I can create Getters here for images, directories .. versions ..etc
// // this will faciliate my operations
