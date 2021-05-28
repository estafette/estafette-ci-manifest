package manifest

import (
	"fmt"
	"strings"

	foundation "github.com/estafette/estafette-foundation"
)

// EstafetteBuilder contains configuration for the ci-builder component
type EstafetteBuilder struct {
	Track           string          `yaml:"track,omitempty"`
	OperatingSystem OperatingSystem `yaml:"os,omitempty"`
	StorageMedium   StorageMedium   `yaml:"medium,omitempty"`
	BuilderType     BuilderType     `yaml:"type,omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteBuilder
func (builder *EstafetteBuilder) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Track           string          `yaml:"track"`
		OperatingSystem OperatingSystem `yaml:"os"`
		StorageMedium   StorageMedium   `yaml:"medium"`
		BuilderType     BuilderType     `yaml:"type,omitempty"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	builder.Track = aux.Track
	builder.OperatingSystem = aux.OperatingSystem
	builder.StorageMedium = aux.StorageMedium
	builder.BuilderType = aux.BuilderType

	return nil
}

// SetDefaults sets default values for properties of EstafetteBuilder if not defined
func (builder *EstafetteBuilder) SetDefaults(preferences EstafetteManifestPreferences) {
	// set default for OperatingSystem if not set
	if builder.OperatingSystem == OperatingSystemUnknown {
		builder.OperatingSystem = preferences.BuilderOperatingSystems[0]
	}
	// set default for Track if not set
	if tracks, ok := preferences.BuilderTracksPerOperatingSystem[builder.OperatingSystem]; ok && len(tracks) > 0 && builder.Track == "" {
		builder.Track = preferences.BuilderTracksPerOperatingSystem[builder.OperatingSystem][0]
	}
	if builder.BuilderType == BuilderTypeUnknown {
		builder.BuilderType = BuilderTypeDocker
	}
}

func (builder *EstafetteBuilder) validate(preferences EstafetteManifestPreferences) (err error) {

	if !OperatingSystemArrayContains(preferences.BuilderOperatingSystems, builder.OperatingSystem) {

		var allowedOperatingSystems []string
		for _, os := range preferences.BuilderOperatingSystems {
			if os != OperatingSystemUnknown {
				allowedOperatingSystems = append(allowedOperatingSystems, string(os))
			}
		}

		return fmt.Errorf("builder os should be one of: %v", strings.Join(allowedOperatingSystems, ", "))
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
