package manifest

// EstafetteTrigger defines an automated trigger to trigger a build or a release
type EstafetteTrigger struct {
	Type      string `yaml:"type" json:"type"`
	Reference string `yaml:"ref" json:"ref"`
	Branch    string `yaml:"branch,omitempty" json:"branch,omitempty"`
}
