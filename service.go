package manifest

// EstafetteService represents a service container to run during a single or multiple stages
type EstafetteService struct {
	Name               string                 `yaml:"name,omitempty"`
	ContainerImage     string                 `yaml:"image,omitempty"`
	Shell              string                 `yaml:"shell,omitempty"`
	Commands           []string               `yaml:"commands,omitempty"`
	When               string                 `yaml:"when,omitempty"`
	EnvVars            map[string]string      `yaml:"env,omitempty"`
	ContinueAfterStage bool                   `yaml:"continueAfterStage,omitempty"`
	Readiness          *ReadinessProbe        `yaml:"readiness,omitempty"`
	CustomProperties   map[string]interface{} `yaml:",inline"`
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
		Name               string                 `yaml:"name,omitempty"`
		ContainerImage     string                 `yaml:"image,omitempty"`
		Shell              string                 `yaml:"shell,omitempty"`
		Commands           []string               `yaml:"commands,omitempty"`
		When               string                 `yaml:"when,omitempty"`
		EnvVars            map[string]string      `yaml:"env,omitempty"`
		ContinueAfterStage bool                   `yaml:"continueAfterStage,omitempty"`
		Readiness          *ReadinessProbe        `yaml:"readiness,omitempty"`
		CustomProperties   map[string]interface{} `yaml:",inline"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	service.Name = aux.Name
	service.ContainerImage = aux.ContainerImage
	service.Shell = aux.Shell
	service.Commands = aux.Commands
	service.When = aux.When
	service.EnvVars = aux.EnvVars
	service.ContinueAfterStage = aux.ContinueAfterStage
	service.Readiness = aux.Readiness

	// fix for map[interface{}]interface breaking json.marshal - see https://github.com/go-yaml/yaml/issues/139
	service.CustomProperties = cleanUpStringMap(aux.CustomProperties)

	return nil
}

// SetDefaults sets default values for properties of EstafetteService if not defined
func (service *EstafetteService) SetDefaults(builder EstafetteBuilder) {
	if service.Readiness != nil {
		service.Readiness.SetDefaults(service.Name)
	}

	// set default for Shell if not set
	if service.Shell == "" {
		if builder.OperatingSystem == "windows" {
			service.Shell = "powershell"
		} else {
			service.Shell = "/bin/sh"
		}
	}

	// set default for When if not set
	if service.When == "" {
		service.When = "status == 'succeeded'"
	}
}

// SetDefaults sets default values for properties of EstafetteService if not defined
func (readiness *ReadinessProbe) SetDefaults(serviceName string) {

	if readiness.Hostname == "" && serviceName != "" {
		readiness.Hostname = serviceName
	}

	if readiness.Port == 0 {
		readiness.Port = 80
	}

	if readiness.Protocol == "" {
		readiness.Protocol = "http"
		if readiness.Port == 443 {
			readiness.Protocol = "https"
		}
	}

	if readiness.TimeoutSeconds == 0 {
		readiness.TimeoutSeconds = 60
	}
}
