package manifest

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

	// map auxiliary properties
	builder.Track = aux.Track

	// set default property values
	builder.SetDefaults()

	return nil
}

// SetDefaults sets default values for properties of EstafetteBuilder if not defined
func (builder *EstafetteBuilder) SetDefaults() {
	// set default for Track if not set
	if builder.Track == "" {
		builder.Track = "stable"
	}
}
