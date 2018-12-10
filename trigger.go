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

// UnmarshalYAML customizes unmarshalling an EstafetteTrigger
func (t *EstafetteTrigger) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Event  string                  `yaml:"event,omitempty" json:"event,omitempty"`
		Filter *EstafetteTriggerFilter `yaml:"filter,omitempty" json:"filter,omitempty"`
		Then   *EstafetteTriggerThen   `yaml:"then,omitempty" json:"then,omitempty"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	t.Event = aux.Event
	t.Filter = aux.Filter
	t.Then = aux.Then

	t.SetDefaults()

	return nil
}

// SetDefaults sets event-specific defaults
func (t *EstafetteTrigger) SetDefaults() {

	// set filter default
	switch t.Event {
	case "pipeline-build-started",
		"pipeline-build-finished",
		"pipeline-release-started",
		"pipeline-release-finished":
		if t.Filter == nil {
			t.Filter = &EstafetteTriggerFilter{}
		}
		if t.Then == nil {
			t.Then = &EstafetteTriggerThen{}
		}
	}

	// set filter branch default
	switch t.Event {
	case "pipeline-build-started",
		"pipeline-build-finished":
		if t.Filter.Branch == "" {
			t.Filter.Branch = "master"
		}

	case "pipeline-release-started",
		"pipeline-release-finished":
		if t.Filter.Branch == "" {
			t.Filter.Branch = ".+"
		}
	}

	// set then branch default
	switch t.Event {
	case "pipeline-build-started",
		"pipeline-build-finished":
		if t.Then.Branch == "" {
			t.Then.Branch = "master"
		}

	case "pipeline-release-started",
		"pipeline-release-finished":
		if t.Then.Branch == "" {
			t.Then.Branch = ".+"
		}
	}

	// set filter pipeline default
	switch t.Event {
	case "pipeline-release-started",
		"pipeline-release-finished":
		if t.Filter.Pipeline == "" {
			t.Filter.Pipeline = "this"
		}
	}

	// set filter pipeline default
	switch t.Event {
	case "pipeline-build-finished",
		"pipeline-release-finished":
		if t.Filter.Status == "" {
			t.Filter.Status = "succeeded"
		}
	}
}
