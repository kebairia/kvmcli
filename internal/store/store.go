package store

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
	Backend string           `yaml:"backend"`
	Config  Config           `yaml:"config"`
	Images  map[string]Image `yaml:"images"`
}

type Config struct {
	ArtifactsPath string `yaml:"artifactsPath"`
	ImagesPath    string `yaml:"imagesPath"`
}

type Image struct {
	Version   string `yaml:"version"`
	Directory string `yaml:"directory"`
	File      string `yaml:"file"`
	Checksum  string `yaml:"checksum"`
	Size      string `yaml:"size"`
}
