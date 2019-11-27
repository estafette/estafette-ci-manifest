package manifest

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// EstafetteStage represents a stage of a build pipeline or release
type EstafetteStage struct {
	Name                    string                 `yaml:"-"`
	ContainerImage          string                 `yaml:"image,omitempty"`
	Shell                   string                 `yaml:"shell,omitempty"`
	WorkingDirectory        string                 `yaml:"workDir,omitempty"`
	Commands                []string               `yaml:"commands,omitempty"`
	RunCommandsInForeground bool                   `yaml:"runCommandsInForeground,omitempty"`
	When                    string                 `yaml:"when,omitempty"`
	EnvVars                 map[string]string      `yaml:"env,omitempty"`
	AutoInjected            bool                   `yaml:"autoInjected,omitempty"`
	Retries                 int                    `yaml:"retries,omitempty"`
	ParallelStages          []*EstafetteStage      `yaml:"parallelStages,omitempty"`
	Services                []*EstafetteService    `yaml:"services,omitempty"`
	CustomProperties        map[string]interface{} `yaml:",inline"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteStage
func (stage *EstafetteStage) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name                    string                 `yaml:"-"`
		ContainerImage          string                 `yaml:"image,omitempty"`
		Shell                   string                 `yaml:"shell,omitempty"`
		WorkingDirectory        string                 `yaml:"workDir,omitempty"`
		Commands                []string               `yaml:"commands,omitempty"`
		RunCommandsInForeground bool                   `yaml:"runCommandsInForeground,omitempty"`
		When                    string                 `yaml:"when,omitempty"`
		EnvVars                 map[string]string      `yaml:"env,omitempty"`
		AutoInjected            bool                   `yaml:"autoInjected,omitempty"`
		Retries                 int                    `yaml:"retries,omitempty"`
		ParallelStages          yaml.MapSlice          `yaml:"parallelStages"`
		Services                []*EstafetteService    `yaml:"services,omitempty"`
		CustomProperties        map[string]interface{} `yaml:",inline"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	stage.ContainerImage = aux.ContainerImage
	stage.Shell = aux.Shell
	stage.WorkingDirectory = aux.WorkingDirectory
	stage.Commands = aux.Commands
	stage.RunCommandsInForeground = aux.RunCommandsInForeground
	stage.When = aux.When
	stage.EnvVars = aux.EnvVars
	stage.AutoInjected = aux.AutoInjected
	stage.Retries = aux.Retries
	stage.Services = aux.Services

	for _, mi := range aux.ParallelStages {

		bytes, err := yaml.Marshal(mi.Value)
		if err != nil {
			return err
		}

		var innerStage *EstafetteStage
		if err := yaml.Unmarshal(bytes, &innerStage); err != nil {
			return err
		}
		if innerStage == nil {
			innerStage = &EstafetteStage{}
		}

		innerStage.Name = mi.Key.(string)
		stage.ParallelStages = append(stage.ParallelStages, innerStage)
	}

	// fix for map[interface{}]interface breaking json.marshal - see https://github.com/go-yaml/yaml/issues/139
	stage.CustomProperties = cleanUpStringMap(aux.CustomProperties)

	return nil
}

// SetDefaults sets default values for properties of EstafetteStage if not defined
func (stage *EstafetteStage) SetDefaults(builder EstafetteBuilder) {
	// set default for Shell if not set
	if len(stage.ParallelStages) == 0 && stage.Shell == "" {
		if builder.OperatingSystem == "windows" {
			stage.Shell = "powershell"
		} else {
			stage.Shell = "/bin/sh"
		}
	}

	// set default for WorkingDirectory if not set
	if len(stage.ParallelStages) == 0 && stage.WorkingDirectory == "" {
		if builder.OperatingSystem == "windows" {
			stage.WorkingDirectory = "C:/estafette-work"
		} else {
			stage.WorkingDirectory = "/estafette-work"
		}
	}

	// set default for When if not set
	if stage.When == "" {
		stage.When = "status == 'succeeded'"
	}

	// set defaults for inner stages
	for _, s := range stage.ParallelStages {
		s.SetDefaults(builder)
	}

	for _, svc := range stage.Services {
		svc.SetDefaults(builder, *stage)
	}
}

// Validate checks whether the stage has valid parameters
func (stage *EstafetteStage) Validate() (err error) {

	if len(stage.ParallelStages) > 0 {
		if stage.ContainerImage != "" {
			return fmt.Errorf("Stage %v cannot use parameters parallelStages and image at the same time", stage.Name)
		}
		if stage.Shell != "" {
			return fmt.Errorf("Stage %v cannot use parameters parallelStages and shell at the same time", stage.Name)
		}
		if stage.WorkingDirectory != "" {
			return fmt.Errorf("Stage %v cannot use parameters parallelStages and workDir at the same time", stage.Name)
		}
		if len(stage.Commands) > 0 {
			return fmt.Errorf("Stage %v cannot use parameters parallelStages and commands at the same time", stage.Name)
		}
		if len(stage.EnvVars) > 0 {
			return fmt.Errorf("Stage %v cannot use parameters parallelStages and env at the same time", stage.Name)
		}
	} else {
		if stage.ContainerImage == "" && len(stage.Services) == 0 {
			return fmt.Errorf("Stage %v has no image set", stage.Name)
		}
		if stage.Retries < 0 {
			return fmt.Errorf("Stage %v has no negative retries; needs to be zero or greater", stage.Name)
		}
	}

	return nil
}
