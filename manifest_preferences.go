package manifest

// EstafetteManifestPreferences is used to configure validation rules for the manifest
type EstafetteManifestPreferences struct {
	LabelRegexes                    map[string]string            `yaml:"labelRegexes,omitempty" json:"labelRegexes,omitempty"`
	BuilderOperatingSystems         []OperatingSystem            `yaml:"builderOperatingSystems,omitempty" json:"builderOperatingSystems,omitempty"`
	BuilderTracksPerOperatingSystem map[OperatingSystem][]string `yaml:"builderTracksPerOperatingSystem,omitempty" json:"builderTracksPerOperatingSystem,omitempty"`
	DefaultBranch                   string                       `yaml:"defaultBranch,omitempty" json:"defaultBranch,omitempty"`
}
