package manifest

// EstafetteTrigger defines an automated trigger to trigger a build or a release
type EstafetteTrigger struct {
	Event  string                  `yaml:"event,omitempty" json:"event,omitempty"`
	Filter *EstafetteTriggerFilter `yaml:"filter,omitempty" json:"filter,omitempty"`
	Then   *EstafetteTriggerThen   `yaml:"then,omitempty" json:"then,omitempty"`
}

// EstafetteTriggerFilter filters the triggered event
type EstafetteTriggerFilter struct {
	// pipeline related filtering
	Pipeline string `yaml:"pipeline,omitempty" json:"pipeline,omitempty"`
	Target   string `yaml:"target,omitempty" json:"target,omitempty"`
	Action   string `yaml:"action,omitempty" json:"action,omitempty"`
	Status   string `yaml:"status,omitempty" json:"status,omitempty"`
	Branch   string `yaml:"branch,omitempty" json:"branch,omitempty"`

	// cron related filtering
	Cron string `yaml:"cron,omitempty" json:"cron,omitempty"`

	// docker related filtering
	Image string `yaml:"image,omitempty" json:"image,omitempty"`
	Tag   string `yaml:"tag,omitempty" json:"tag,omitempty"`
}

// EstafetteTriggerThen determines what's triggered when filtered events arive
type EstafetteTriggerThen struct {
	Branch string `yaml:"branch,omitempty" json:"branch,omitempty"`
	Action string `yaml:"action,omitempty" json:"action,omitempty"`
}
