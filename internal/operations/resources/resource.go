package resources

type Resource interface {
	Create(name string, xml []byte) error
	Delete(name string) error
	// Get() ([]string, error)
	// Edit(name string, params map[string]any) error
}
