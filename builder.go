package manifest

import "fmt"

// EstafetteBuilder contains configuration for the ci-builder component
type EstafetteBuilder struct {
	Track           string `yaml:"track,omitempty"`
	OperatingSystem string `yaml:"os,omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteBuilder
func (builder *EstafetteBuilder) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Track           string `yaml:"track"`
		OperatingSystem string `yaml:"os,omitempty"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	builder.Track = aux.Track
	builder.OperatingSystem = aux.OperatingSystem

	// set default property values
	builder.SetDefaults()

	return nil
}

// SetDefaults sets default values for properties of EstafetteBuilder if not defined
func (builder *EstafetteBuilder) SetDefaults() {
	// set default for Track if not set
	if builder.OperatingSystem == "" {
		builder.OperatingSystem = "linux"
	}
	if builder.Track == "" && builder.OperatingSystem == "linux" {
		builder.Track = "stable"
	}
	if builder.Track == "" && builder.OperatingSystem == "windows" {
		builder.Track = "windowsservercore-1809"
	}
}

func (builder *EstafetteBuilder) validate() (err error) {
	if builder.OperatingSystem == "linux" && builder.Track != "dev" && builder.Track != "beta" && builder.Track != "stable" {
		return fmt.Errorf("builder track should be one of: dev, beta or stable")
	}
	if builder.OperatingSystem != "linux" && builder.OperatingSystem != "windows" {
		return fmt.Errorf("builder os should be one of: linux or windows")
	}
	return nil
}
