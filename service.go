package manifest

// EstafetteService represents a service container to run during a single or multiple stages
type EstafetteService struct {
	Name               string                  `yaml:"name,omitempty"`
	ContainerImage     string                  `yaml:"image,omitempty"`
	EnvVars            map[string]string       `yaml:"env,omitempty"`
	Ports              []*EstafetteServicePort `yaml:"ports,omitempty"`
	Command            string                  `yaml:"command,omitempty"`
	ContinueAfterStage bool                    `yaml:"continueAfterStage,omitempty"`
}

// EstafetteServicePort represents a port to be opened on the service container and an optional external port if port mapping to a different port is required
type EstafetteServicePort struct {
	Port     int  `yaml:"port,omitempty"`
	HostPort *int `yaml:"hostPort,omitempty"`
}
