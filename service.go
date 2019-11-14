package manifest

// EstafetteService represents a service container to run during a single or multiple stages
type EstafetteService struct {
	Name               string                  `yaml:"name,omitempty"`
	ContainerImage     string                  `yaml:"image,omitempty"`
	EnvVars            map[string]string       `yaml:"env,omitempty"`
	Ports              []*EstafetteServicePort `yaml:"ports,omitempty"`
	Command            string                  `yaml:"command,omitempty"`
	ContinueAfterStage bool                    `yaml:"continueAfterStage,omitempty"`
	Readiness          *ReadinessProbe         `yaml:"readiness,omitempty"`
	CustomProperties   map[string]interface{}  `yaml:",inline"`
}

// EstafetteServicePort represents a port to be opened on the service container and an optional external port if port mapping to a different port is required
type EstafetteServicePort struct {
	Port      int             `yaml:"port,omitempty"`
	HostPort  *int            `yaml:"hostPort,omitempty"`
	Readiness *ReadinessProbe `yaml:"readiness,omitempty"`
}

// ReadinessProbe defines an http readiness probe
type ReadinessProbe struct {
	Path           string `yaml:"path,omitempty"`
	TimeoutSeconds int    `yaml:"timeoutSeconds,omitempty"`
	Port           int    `yaml:"port,omitempty"`
	Protocol       string `yaml:"protocol,omitempty"`
	Hostname       string `yaml:"hostname,omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteService
func (service *EstafetteService) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name               string                  `yaml:"name,omitempty"`
		ContainerImage     string                  `yaml:"image,omitempty"`
		EnvVars            map[string]string       `yaml:"env,omitempty"`
		Ports              []*EstafetteServicePort `yaml:"ports,omitempty"`
		Command            string                  `yaml:"command,omitempty"`
		ContinueAfterStage bool                    `yaml:"continueAfterStage,omitempty"`
		CustomProperties   map[string]interface{}  `yaml:",inline"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	service.Name = aux.Name
	service.ContainerImage = aux.ContainerImage
	service.EnvVars = aux.EnvVars
	service.Ports = aux.Ports
	service.Command = aux.Command
	service.ContinueAfterStage = aux.ContinueAfterStage

	// fix for map[interface{}]interface breaking json.marshal - see https://github.com/go-yaml/yaml/issues/139
	service.CustomProperties = cleanUpStringMap(aux.CustomProperties)

	return nil
}

// SetDefaults sets default values for properties of EstafetteService if not defined
func (service *EstafetteService) SetDefaults() {
	if service.Readiness != nil {
		service.Readiness.SetDefaults(service.Name, nil)
	}

	for _, p := range service.Ports {
		p.SetDefaults(service.Name)
	}
}

// SetDefaults sets default values for properties of EstafetteServicePort if not defined
func (port *EstafetteServicePort) SetDefaults(serviceName string) {
	if port.Readiness != nil {
		port.Readiness.SetDefaults(serviceName, port.HostPort)
	}
}

// SetDefaults sets default values for properties of EstafetteService if not defined
func (readiness *ReadinessProbe) SetDefaults(serviceName string, port *int) {

	if readiness.Hostname == "" && serviceName != "" {
		readiness.Hostname = serviceName
	}

	if readiness.Port == 0 && port != nil {
		readiness.Port = *port
	} else if readiness.Port == 0 {
		readiness.Port = 80
	}

	if readiness.Protocol == "" {
		readiness.Protocol = "http"
	}
}
