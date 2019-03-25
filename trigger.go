package manifest

import (
	"fmt"
)

// EstafetteTrigger represents a trigger of any supported type and what action to take if the trigger fired
type EstafetteTrigger struct {
	Pipeline *EstafettePipelineTrigger `yaml:"pipeline,omitempty"`
	Git      *EstafetteGitTrigger      `yaml:"git,omitempty"`
	Docker   *EstafetteDockerTrigger   `yaml:"docker,omitempty"`
	Cron     *EstafetteCronTrigger     `yaml:"cron,omitempty"`
}

// EstafettePipelineTrigger fires for pipeline changes and applies filtering to limit when this results in an action
type EstafettePipelineTrigger struct {
	Event  string `yaml:"event,omitempty"`
	Name   string `yaml:"name,omitempty"`
	Branch string `yaml:"branch,omitempty"`
}

// EstafetteGitTrigger fires for git repository changes and applies filtering to limit when this results in an action
type EstafetteGitTrigger struct {
	Event      string `yaml:"event,omitempty"`
	Repository string `yaml:"repository,omitempty"`
	Branch     string `yaml:"branch,omitempty"`
}

// EstafetteDockerTrigger fires for docker image changes and applies filtering to limit when this results in an action
type EstafetteDockerTrigger struct {
	Event string `yaml:"event,omitempty"`
	Image string `yaml:"image,omitempty"`
	Tag   string `yaml:"tag,omitempty"`
}

// EstafetteCronTrigger fires at intervals specified by the cron expression
type EstafetteCronTrigger struct {
	Expression string `yaml:"expression,omitempty"`
}

// SetDefaults sets defaults for EstafetteTrigger
func (t *EstafetteTrigger) SetDefaults() {
	if t.Pipeline != nil {
		t.Pipeline.SetDefaults()
	}
	if t.Git != nil {
		t.Git.SetDefaults()
	}
	if t.Docker != nil {
		t.Docker.SetDefaults()
	}
	if t.Cron != nil {
		t.Cron.SetDefaults()
	}
}

// SetDefaults sets defaults for EstafettePipelineTrigger
func (p *EstafettePipelineTrigger) SetDefaults() {
	if p.Event == "" {
		p.Event = "succeeded"
	}
	if p.Branch == "" {
		p.Branch = "master"
	}
}

// SetDefaults sets defaults for EstafetteGitTrigger
func (g *EstafetteGitTrigger) SetDefaults() {
	if g.Event == "" {
		g.Event = "push"
	}
	if g.Branch == "" {
		g.Branch = "master"
	}
}

// SetDefaults sets defaults for EstafetteDockerTrigger
func (d *EstafetteDockerTrigger) SetDefaults() {
}

// SetDefaults sets defaults for EstafetteCronTrigger
func (c *EstafetteCronTrigger) SetDefaults() {
}

// Validate checks if EstafetteTrigger is valid
func (t *EstafetteTrigger) Validate() (err error) {

	numberOfTypes := 0

	if t.Pipeline == nil &&
		t.Git == nil &&
		t.Docker == nil &&
		t.Cron == nil {
		return fmt.Errorf("Set at least a pipeline, git, docker or cron trigger")
	}

	if t.Pipeline != nil {
		err = t.Pipeline.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.Git != nil {
		err = t.Git.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.Docker != nil {
		err = t.Docker.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.Cron != nil {
		err = t.Cron.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}

	if numberOfTypes != 1 {
		return fmt.Errorf("Specify at least and at most one type of trigger, 'pipeline', 'git', 'docker' or 'cron'")
	}

	return nil
}

// Validate checks if EstafettePipelineTrigger is valid
func (p *EstafettePipelineTrigger) Validate() (err error) {
	if p.Event != "failed" && p.Event != "succeeded" && p.Event != "finished" {
		return fmt.Errorf("Set pipeline.event in your trigger to 'failed', 'succeeded' or 'finished'")
	}
	if p.Name == "" {
		return fmt.Errorf("Set pipeline.name in your trigger to a full qualified pipeline name, i.e. github.com/estafette/estafette-ci-manifest")
	}
	return nil
}

// Validate checks if EstafetteGitTrigger is valid
func (g *EstafetteGitTrigger) Validate() (err error) {
	if g.Event != "push" {
		return fmt.Errorf("Set git.event in your trigger to 'push'")
	}
	if g.Repository == "" {
		return fmt.Errorf("Set git.repository in your trigger to a full qualified git repository name, i.e. github.com/estafette/estafette-ci-manifest")
	}
	return nil
}

// Validate checks if EstafetteDockerTrigger is valid
func (d *EstafetteDockerTrigger) Validate() (err error) {
	return nil
}

// Validate checks if EstafetteCronTrigger is valid
func (c *EstafetteCronTrigger) Validate() (err error) {
	return nil
}