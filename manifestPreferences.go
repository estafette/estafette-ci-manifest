package manifest

// EstafetteManifestPreferences is used to configure validation rules for the manifest
type EstafetteManifestPreferences struct {
	LabelRegexes                    map[string]string   `yaml:"labelRegexes,omitempty" json:"labelRegexes,omitempty"`
	BuilderOperatingSystems         []string            `yaml:"builderOperatingSystems,omitempty" json:"builderOperatingSystems,omitempty"`
	BuilderTracksPerOperatingSystem map[string][]string `yaml:"builderTracksPerOperatingSystem,omitempty" json:"builderTracksPerOperatingSystem,omitempty"`
	DefaultBranch                   string              `yaml:"defaultBranch,omitempty" json:"defaultBranch,omitempty"`
}
