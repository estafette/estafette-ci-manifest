package manifest

import (
	yaml "gopkg.in/yaml.v2"
)

// EstafetteBot allows to respond to any event coming from one of the integrations
type EstafetteBot struct {
	Name   string            `yaml:"-"`
	Events []string          `yaml:"events,omitempty" json:",omitempty"`
	Stages []*EstafetteStage `yaml:"-" json:",omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteBot
func (bot *EstafetteBot) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	var aux struct {
		Name   string        `yaml:"-"`
		Events []string      `yaml:"events"`
		Stages yaml.MapSlice `yaml:"stages"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	// map auxiliary properties
	bot.Events = aux.Events

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
		Name   string        `yaml:"-"`
		Events []string      `yaml:"events,omitempty"`
		Stages yaml.MapSlice `yaml:"stages,omitempty"`
	}

	// map auxiliary properties
	aux.Events = bot.Events

	for _, stage := range bot.Stages {
		aux.Stages = append(aux.Stages, yaml.MapItem{
			Key:   stage.Name,
			Value: stage,
		})
	}

	return aux, err
}
