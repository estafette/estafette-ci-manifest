package manifest

import (
	yaml "gopkg.in/yaml.v2"
)

// EstafetteRelease represents a release target that in itself contains one or multiple stages
type EstafetteRelease struct {
	Name            string                    `yaml:"-"`
	Builder         *EstafetteBuilder         `yaml:"builder,omitempty"`
	CloneRepository bool                      `yaml:"clone,omitempty" json:",omitempty"`
	Actions         []*EstafetteReleaseAction `yaml:"actions,omitempty" json:",omitempty"`
	Triggers        []*EstafetteTrigger       `yaml:"triggers,omitempty" json:",omitempty"`
	Stages          []*EstafetteStage         `yaml:"-"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteRelease
func (release *EstafetteRelease) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name            string                    `yaml:"name"`
		Builder         *EstafetteBuilder         `yaml:"builder"`
		CloneRepository bool                      `yaml:"clone"`
		Actions         []*EstafetteReleaseAction `yaml:"actions"`
		Triggers        []*EstafetteTrigger       `yaml:"triggers"`
		Stages          yaml.MapSlice             `yaml:"stages"`
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
		CloneRepository bool                      `yaml:"clone,omitempty"`
		Actions         []*EstafetteReleaseAction `yaml:"actions,omitempty"`
		Triggers        []*EstafetteTrigger       `yaml:"triggers,omitempty"`
		Stages          yaml.MapSlice             `yaml:"stages,omitempty"`
	}

	// map auxiliary properties
	aux.Builder = release.Builder
	aux.CloneRepository = release.CloneRepository
	aux.Actions = release.Actions
	aux.Triggers = release.Triggers

	for _, stage := range release.Stages {
		aux.Stages = append(aux.Stages, yaml.MapItem{
			Key:   stage.Name,
			Value: stage,
		})
	}

	return aux, err
}
