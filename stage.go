package manifest

import (
	"github.com/rs/zerolog/log"
)

// EstafetteStage represents a stage of a build pipeline or release
type EstafetteStage struct {
	Name             string
	ContainerImage   string            `yaml:"image,omitempty"`
	Shell            string            `yaml:"shell,omitempty"`
	WorkingDirectory string            `yaml:"workDir,omitempty"`
	Commands         []string          `yaml:"commands,omitempty"`
	When             string            `yaml:"when,omitempty"`
	EnvVars          map[string]string `yaml:"env,omitempty"`
	AutoInjected     bool              `yaml:"autoInjected,omitempty"`
	CustomProperties map[string]interface{}
}

// UnmarshalYAML customizes unmarshalling an EstafetteStage
func (stage *EstafetteStage) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		ContainerImage   string            `yaml:"image,omitempty"`
		Shell            string            `yaml:"shell,omitempty"`
		WorkingDirectory string            `yaml:"workDir,omitempty"`
		Commands         []string          `yaml:"commands,omitempty"`
		When             string            `yaml:"when,omitempty"`
		EnvVars          map[string]string `yaml:"env,omitempty"`
		AutoInjected     bool              `yaml:"autoInjected,omitempty"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	log.Debug().Interface("aux", aux).Msg("Unmarshalled auxiliary type for EstafetteStage")

	// map auxiliary properties
	stage.ContainerImage = aux.ContainerImage
	stage.Shell = aux.Shell
	stage.WorkingDirectory = aux.WorkingDirectory
	stage.Commands = aux.Commands
	stage.When = aux.When
	stage.EnvVars = aux.EnvVars
	stage.AutoInjected = aux.AutoInjected

	// get all properties in a map
	propertiesMap := make(map[string]interface{})
	if err := unmarshal(&propertiesMap); err != nil {
		return err
	}

	// remove reserved properties
	for _, reservedPropertyName := range getReservedPropertyNames() {
		if _, ok := propertiesMap[reservedPropertyName]; ok {
			delete(propertiesMap, reservedPropertyName)
		}
	}

	// copy remaining (custom) properties
	stage.CustomProperties = propertiesMap

	// set default property values
	stage.SetDefaults()

	return nil
}

// SetDefaults sets default values for properties of EstafetteStage if not defined
func (stage *EstafetteStage) SetDefaults() {
	// set default for Shell if not set
	if stage.Shell == "" {
		stage.Shell = "/bin/sh"
	}

	// set default for WorkingDirectory if not set
	if stage.WorkingDirectory == "" {
		stage.WorkingDirectory = "/estafette-work"
	}

	// set default for When if not set
	if stage.When == "" {
		stage.When = "status == 'succeeded'"
	}

}
