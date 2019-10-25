package manifest

// EstafetteService represents a service container to run during a single or multiple stages
type EstafetteService struct {
	ContainerImage string            `yaml:"image,omitempty"`
	EnvVars        map[string]string `yaml:"env,omitempty"`
}
