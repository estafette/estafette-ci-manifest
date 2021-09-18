package manifest

// EstafetteManifestPreferences is used to configure validation rules for the manifest
type EstafetteManifestPreferences struct {
	LabelRegexes                    map[string]string            `yaml:"labelRegexes,omitempty" json:"labelRegexes,omitempty"`
	BuilderOperatingSystems         []OperatingSystem            `yaml:"builderOperatingSystems,omitempty" json:"builderOperatingSystems,omitempty"`
	BuilderTracksPerOperatingSystem map[OperatingSystem][]string `yaml:"builderTracksPerOperatingSystem,omitempty" json:"builderTracksPerOperatingSystem,omitempty"`
	DefaultBranch                   string                       `yaml:"defaultBranch,omitempty" json:"defaultBranch,omitempty"`
}

func (p *EstafetteManifestPreferences) SetDefaults() {
	if p.LabelRegexes == nil {
		p.LabelRegexes = make(map[string]string)
	}

	if len(p.BuilderOperatingSystems) == 0 {
		p.BuilderOperatingSystems = []OperatingSystem{OperatingSystemLinux, OperatingSystemWindows}
	}

	if len(p.BuilderTracksPerOperatingSystem) == 0 {
		p.BuilderTracksPerOperatingSystem = map[OperatingSystem][]string{
			OperatingSystemLinux:   {"stable", "beta", "dev"},
			OperatingSystemWindows: {"nanoserver-1809-stable", "nanoserver-1809-beta", "nanoserver-1809-dev"},
		}
	}

	if p.DefaultBranch == "" {
		p.DefaultBranch = "master"
	}
}
