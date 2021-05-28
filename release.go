package manifest

import (
	"github.com/jinzhu/copier"
	yaml "gopkg.in/yaml.v2"
)

// EstafetteRelease represents a release target that in itself contains one or multiple stages
type EstafetteRelease struct {
	Name            string                    `yaml:"-"`
	Builder         *EstafetteBuilder         `yaml:"builder,omitempty"`
	CloneRepository *bool                     `yaml:"clone,omitempty" json:",omitempty"`
	Actions         []*EstafetteReleaseAction `yaml:"actions,omitempty" json:",omitempty"`
	Triggers        []*EstafetteTrigger       `yaml:"triggers,omitempty" json:",omitempty"`
	Stages          []*EstafetteStage         `yaml:"-" json:",omitempty"`
	Template        string                    `yaml:"template,omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteRelease
func (release *EstafetteRelease) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name            string                    `yaml:"name"`
		Builder         *EstafetteBuilder         `yaml:"builder"`
		CloneRepository *bool                     `yaml:"clone"`
		Actions         []*EstafetteReleaseAction `yaml:"actions"`
		Triggers        []*EstafetteTrigger       `yaml:"triggers"`
		Stages          yaml.MapSlice             `yaml:"stages"`
		Template        string                    `yaml:"template"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	release.Name = aux.Name
	release.Builder = aux.Builder
	release.CloneRepository = aux.CloneRepository
	release.Actions = aux.Actions
	release.Triggers = aux.Triggers
	release.Template = aux.Template

	for _, mi := range aux.Stages {

		bytes, err := yaml.Marshal(mi.Value)
		if err != nil {
			return err
		}

		var stage *EstafetteStage
		if err := yaml.Unmarshal(bytes, &stage); err != nil {
			return err
		}
		if stage == nil {
			stage = &EstafetteStage{}
		}

		// set the stage name, overwriting the name property if set on the stage explicitly
		stage.Name = mi.Key.(string)

		release.Stages = append(release.Stages, stage)
	}

	return nil
}

// MarshalYAML customizes marshalling an EstafetteManifest
func (release EstafetteRelease) MarshalYAML() (out interface{}, err error) {

	var aux struct {
		Name            string                    `yaml:"-"`
		Builder         *EstafetteBuilder         `yaml:"builder,omitempty"`
		CloneRepository *bool                     `yaml:"clone,omitempty"`
		Actions         []*EstafetteReleaseAction `yaml:"actions,omitempty"`
		Triggers        []*EstafetteTrigger       `yaml:"triggers,omitempty"`
		Stages          yaml.MapSlice             `yaml:"stages,omitempty"`
		Template        string                    `yaml:"template,omitempty"`
	}

	// map auxiliary properties
	aux.Builder = release.Builder
	aux.CloneRepository = release.CloneRepository
	aux.Actions = release.Actions
	aux.Triggers = release.Triggers
	aux.Template = release.Template

	for _, stage := range release.Stages {
		aux.Stages = append(aux.Stages, yaml.MapItem{
			Key:   stage.Name,
			Value: stage,
		})
	}

	return aux, err
}

// DeepCopy provides a copy of all nested pointers
func (release EstafetteRelease) DeepCopy() (target EstafetteRelease) {

	copier.CopyWithOption(&target, release, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return
}

// InitFromTemplate uses template values for
func (release *EstafetteRelease) InitFromTemplate(releaseTemplates map[string]*EstafetteReleaseTemplate) {

	if release.Template != "" {
		// check if template with defined name exists, and use its values overridden by this releases values

		if releaseTemplate, found := releaseTemplates[release.Template]; found && releaseTemplate != nil {

			// deep copy so there's no pointers shared with other releases
			template := releaseTemplate.DeepCopy()

			if release.Builder != nil {
				template.Builder = release.Builder
			} else {
				release.Builder = template.Builder
			}

			if release.CloneRepository != nil {
				template.CloneRepository = release.CloneRepository
			} else {
				release.CloneRepository = template.CloneRepository
			}

			if release.Actions != nil && len(release.Actions) > 0 {
				template.Actions = release.Actions
			} else {
				release.Actions = template.Actions
			}

			if release.Triggers != nil && len(release.Triggers) > 0 {
				template.Triggers = release.Triggers
			} else {
				release.Triggers = template.Triggers
			}

			if release.Stages != nil && len(release.Stages) > 0 {
				template.Stages = release.Stages
			} else {
				release.Stages = template.Stages
			}
		}
	}
}
