package manifest

import "github.com/rs/zerolog/log"

// EstafetteBuilder contains configuration for the ci-builder component
type EstafetteBuilder struct {
	Track string `yaml:"track,omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteBuilder
func (builder *EstafetteBuilder) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Track string `yaml:"track,omitempty"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	log.Debug().Interface("aux", aux).Msg("Unmarshalled auxiliary type for EstafetteBuilder")

	// map auxiliary properties
	builder.Track = aux.Track

	// set default property values
	builder.SetDefaults()

	log.Debug().Interface("builder", builder).Msg("Copied auxiliary type properties for EstafetteBuilder")

	return nil
}

// SetDefaults sets default values for properties of EstafetteBuilder if not defined
func (builder *EstafetteBuilder) SetDefaults() {
	// set default for Track if not set
	if builder.Track == "" {
		builder.Track = "stable"
	}
}
