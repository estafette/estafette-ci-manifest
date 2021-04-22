package manifest

import (
	"github.com/jinzhu/copier"
	yaml "gopkg.in/yaml.v2"
)

// EstafetteReleaseTemplate represents a template for a release target
type EstafetteReleaseTemplate struct {
	Name            string                    `yaml:"-"`
	Builder         *EstafetteBuilder         `yaml:"builder,omitempty"`
	CloneRepository *bool                     `yaml:"clone,omitempty" json:",omitempty"`
	Actions         []*EstafetteReleaseAction `yaml:"actions,omitempty" json:",omitempty"`
	Triggers        []*EstafetteTrigger       `yaml:"triggers,omitempty" json:",omitempty"`
	Stages          []*EstafetteStage         `yaml:"-"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteRelease
func (releaseTemplate *EstafetteReleaseTemplate) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name            string                    `yaml:"name"`
		Builder         *EstafetteBuilder         `yaml:"builder"`
		CloneRepository *bool                     `yaml:"clone"`
		Actions         []*EstafetteReleaseAction `yaml:"actions"`
		Triggers        []*EstafetteTrigger       `yaml:"triggers"`
		Stages          yaml.MapSlice             `yaml:"stages"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	releaseTemplate.Name = aux.Name
	releaseTemplate.Builder = aux.Builder
	releaseTemplate.CloneRepository = aux.CloneRepository
	releaseTemplate.Actions = aux.Actions
	releaseTemplate.Triggers = aux.Triggers

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

		releaseTemplate.Stages = append(releaseTemplate.Stages, stage)
	}

	return nil
}

// MarshalYAML customizes marshalling an EstafetteManifest
func (releaseTemplate EstafetteReleaseTemplate) MarshalYAML() (out interface{}, err error) {

	var aux struct {
		Name            string                    `yaml:"-"`
		Builder         *EstafetteBuilder         `yaml:"builder,omitempty"`
		CloneRepository *bool                     `yaml:"clone,omitempty"`
		Actions         []*EstafetteReleaseAction `yaml:"actions,omitempty"`
		Triggers        []*EstafetteTrigger       `yaml:"triggers,omitempty"`
		Stages          yaml.MapSlice             `yaml:"stages,omitempty"`
	}

	// map auxiliary properties
	aux.Builder = releaseTemplate.Builder
	aux.CloneRepository = releaseTemplate.CloneRepository
	aux.Actions = releaseTemplate.Actions
	aux.Triggers = releaseTemplate.Triggers

	for _, stage := range releaseTemplate.Stages {
		aux.Stages = append(aux.Stages, yaml.MapItem{
			Key:   stage.Name,
			Value: stage,
		})
	}

	return aux, err
}

// DeepCopy provides a copy of all nested pointers
func (releaseTemplate EstafetteReleaseTemplate) DeepCopy() (target EstafetteReleaseTemplate) {

	copier.CopyWithOption(&target, releaseTemplate, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return
}
