package manifest

import (
	"fmt"
	"strings"

	foundation "github.com/estafette/estafette-foundation"
)

// EstafetteBuilder contains configuration for the ci-builder component
type EstafetteBuilder struct {
	Track           string        `yaml:"track,omitempty"`
	OperatingSystem string        `yaml:"os,omitempty"`
	StorageMedium   StorageMedium `yaml:"medium,omitempty"`
	MTU             int           `yaml:"mtu,omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteBuilder
func (builder *EstafetteBuilder) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Track           string        `yaml:"track"`
		OperatingSystem string        `yaml:"os"`
		StorageMedium   StorageMedium `yaml:"medium"`
		MTU             int           `yaml:"mtu"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	builder.Track = aux.Track
	builder.OperatingSystem = aux.OperatingSystem
	builder.StorageMedium = aux.StorageMedium
	builder.MTU = aux.MTU

	return nil
}

// SetDefaults sets default values for properties of EstafetteBuilder if not defined
func (builder *EstafetteBuilder) SetDefaults(preferences EstafetteManifestPreferences) {
	// set default for OperatingSystem if not set
	if builder.OperatingSystem == "" {
		builder.OperatingSystem = preferences.BuilderOperatingSystems[0]
	}
	// set default for Track if not set
	if builder.Track == "" {
		builder.Track = preferences.BuilderTracksPerOperatingSystem[builder.OperatingSystem][0]
	}
}

func (builder *EstafetteBuilder) validate(preferences EstafetteManifestPreferences) (err error) {

	if !foundation.StringArrayContains(preferences.BuilderOperatingSystems, builder.OperatingSystem) {
		return fmt.Errorf("builder os should be one of: %v", strings.Join(preferences.BuilderOperatingSystems, ", "))
	}

	tracks, ok := preferences.BuilderTracksPerOperatingSystem[builder.OperatingSystem]

	if !ok {
		return fmt.Errorf("no track preferences have been configured for os %v", builder.OperatingSystem)
	}

	if !foundation.StringArrayContains(tracks, builder.Track) {
		return fmt.Errorf("builder track should be one of: %v", strings.Join(tracks, ", "))
	}

	return nil
}
