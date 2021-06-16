package manifest

// EstafetteService represents a service container to run during a single or multiple stages
type EstafetteService struct {
	Name                    string                 `yaml:"name,omitempty"`
	ContainerImage          string                 `yaml:"image,omitempty"`
	Shell                   string                 `yaml:"shell,omitempty"`
	Commands                []string               `yaml:"commands,omitempty"`
	RunCommandsInForeground bool                   `yaml:"runCommandsInForeground,omitempty"`
	MultiStage              *bool                  `yaml:"multiStage,omitempty"`
	When                    string                 `yaml:"when,omitempty"`
	EnvVars                 map[string]string      `yaml:"env,omitempty"`
	Readiness               *ReadinessProbe        `yaml:"readiness,omitempty"`
	ReadinessProbe          *ReadinessProbe        `yaml:"readinessProbe,omitempty"`
	CustomProperties        map[string]interface{} `yaml:",inline"`
}

// ReadinessProbe defines an http readiness probe
type ReadinessProbe struct {
	HttpGet        *HttpGetProbe `yaml:"httpGet,omitempty"`
	Exec           *ExecProbe    `yaml:"exec,omitempty"`
	TimeoutSeconds int           `yaml:"timeoutSeconds,omitempty"`

	// deprecated
	Path     string `yaml:"path,omitempty"`     // httpGet.path
	Port     int    `yaml:"port,omitempty"`     // httpGet.port
	Protocol string `yaml:"protocol,omitempty"` // httpGet.scheme
	Hostname string `yaml:"hostname,omitempty"` // httpGet.host
}

type HttpGetProbe struct {
	Path   string `yaml:"path,omitempty"`
	Port   int    `yaml:"port,omitempty"`
	Host   string `yaml:"host,omitempty"`
	Scheme string `yaml:"scheme,omitempty"`
}

type ExecProbe struct {
	Command []string `yaml:"command,omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteService
func (service *EstafetteService) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name                    string                 `yaml:"name,omitempty"`
		ContainerImage          string                 `yaml:"image,omitempty"`
		Shell                   string                 `yaml:"shell,omitempty"`
		Commands                []string               `yaml:"commands,omitempty"`
		RunCommandsInForeground bool                   `yaml:"runCommandsInForeground,omitempty"`
		MultiStage              *bool                  `yaml:"multiStage,omitempty"`
		When                    string                 `yaml:"when,omitempty"`
		EnvVars                 map[string]string      `yaml:"env,omitempty"`
		Readiness               *ReadinessProbe        `yaml:"readiness,omitempty"`
		ReadinessProbe          *ReadinessProbe        `yaml:"readinessProbe,omitempty"`
		CustomProperties        map[string]interface{} `yaml:",inline"`
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
	service.RunCommandsInForeground = aux.RunCommandsInForeground
	service.MultiStage = aux.MultiStage
	service.When = aux.When
	service.EnvVars = aux.EnvVars
	service.Readiness = aux.Readiness
	service.ReadinessProbe = aux.ReadinessProbe

	// fix for map[interface{}]interface breaking json.marshal - see https://github.com/go-yaml/yaml/issues/139
	service.CustomProperties = cleanUpStringMap(aux.CustomProperties)

	return nil
}

// SetDefaults sets default values for properties of EstafetteService if not defined
func (service *EstafetteService) SetDefaults(builder EstafetteBuilder, parentStage EstafetteStage) {
	if service.Readiness != nil {
		service.Readiness.SetDefaults(service.Name)
	}
	if service.ReadinessProbe != nil {
		service.ReadinessProbe.SetDefaults(service.Name)
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

	// set default for multistage depending on whether parent stage has image and commands or not
	if service.MultiStage == nil {
		trueValue := true
		falseValue := false
		if parentStage.ContainerImage == "" {
			service.MultiStage = &trueValue
		} else {
			service.MultiStage = &falseValue
		}
	}
}

// SetDefaults sets default values for properties of EstafetteService if not defined
func (readiness *ReadinessProbe) SetDefaults(serviceName string) {

	if readiness.TimeoutSeconds == 0 {
		readiness.TimeoutSeconds = 60
	}

	if readiness.HttpGet != nil {
		if readiness.HttpGet.Host == "" && serviceName != "" {
			readiness.HttpGet.Host = serviceName
		}
		if readiness.HttpGet.Port == 0 {
			readiness.HttpGet.Port = 80
		}
		if readiness.HttpGet.Scheme == "" {
			readiness.HttpGet.Scheme = "http"
			if readiness.HttpGet.Port == 443 {
				readiness.HttpGet.Scheme = "https"
			}
		}

		return
	}

	if readiness.Exec != nil {
		return
	}

	// legacy settings
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
}
