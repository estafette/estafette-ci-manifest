package manifest

import (
	yaml "gopkg.in/yaml.v2"
)

// EstafetteBot allows to respond to any event coming from one of the integrations
type EstafetteBot struct {
	Name            string              `yaml:"-"`
	Builder         *EstafetteBuilder   `yaml:"builder,omitempty"`
	CloneRepository *bool               `yaml:"clone,omitempty" json:",omitempty"`
	Triggers        []*EstafetteTrigger `yaml:"triggers,omitempty" json:",omitempty"`
	Stages          []*EstafetteStage   `yaml:"-" json:",omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteBot
func (bot *EstafetteBot) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name            string              `yaml:"-"`
		Builder         *EstafetteBuilder   `yaml:"builder"`
		CloneRepository *bool               `yaml:"clone"`
		Triggers        []*EstafetteTrigger `yaml:"triggers"`
		Stages          yaml.MapSlice       `yaml:"stages"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	bot.Name = aux.Name
	bot.Builder = aux.Builder
	bot.CloneRepository = aux.CloneRepository
	bot.Triggers = aux.Triggers

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

		stage.Name = mi.Key.(string)
		bot.Stages = append(bot.Stages, stage)
	}

	return nil
}

// MarshalYAML customizes marshalling an EstafetteBot
func (bot *EstafetteBot) MarshalYAML() (out interface{}, err error) {

	var aux struct {
		Name            string              `yaml:"-"`
		Builder         *EstafetteBuilder   `yaml:"builder,omitempty"`
		CloneRepository *bool               `yaml:"clone,omitempty"`
		Triggers        []*EstafetteTrigger `yaml:"triggers,omitempty"`
		Stages          yaml.MapSlice       `yaml:"stages,omitempty"`
	}

	// map auxiliary properties
	aux.Builder = bot.Builder
	aux.CloneRepository = bot.CloneRepository
	aux.Triggers = bot.Triggers

	for _, stage := range bot.Stages {
		aux.Stages = append(aux.Stages, yaml.MapItem{
			Key:   stage.Name,
			Value: stage,
		})
	}

	return aux, err
}
