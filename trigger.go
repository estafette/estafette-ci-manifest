package manifest

// EstafetteTrigger defines an automated trigger to trigger a build or a release
type EstafetteTrigger struct {
	Type      string `yaml:"type" json:"type"`
	Reference string `yaml:"ref" json:"ref"`
	Filter    string `yaml:"filter,omitempty" json:"filter,omitempty"`
}
