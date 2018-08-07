package manifest

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"

	yaml "gopkg.in/yaml.v2"
)

// EstafetteManifest is the object that the .estafette.yaml deserializes to
type EstafetteManifest struct {
	Builder       EstafetteBuilder    `yaml:"builder,omitempty"`
	Version       EstafetteVersion    `yaml:"version,omitempty"`
	Labels        map[string]string   `yaml:"labels,omitempty"`
	GlobalEnvVars map[string]string   `yaml:"env,omitempty"`
	Stages        []*EstafetteStage   `yaml:"stagesdummy,omitempty" json:"Pipelines,omitempty"`
	Releases      []*EstafetteRelease `yaml:"releasesdummy,omitempty"`
}

// UnmarshalYAML customizes unmarshalling an EstafetteManifest
func (c *EstafetteManifest) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {

	log.Debug().Msg("Unmarshaling EstafetteManifest with custom code")

	var aux struct {
		Builder       EstafetteBuilder  `yaml:"builder,omitempty"`
		Version       EstafetteVersion  `yaml:"version,omitempty"`
		Labels        map[string]string `yaml:"labels,omitempty"`
		GlobalEnvVars map[string]string `yaml:"env,omitempty"`
		Stages        yaml.MapSlice     `yaml:"stages,omitempty"`
		Releases      yaml.MapSlice     `yaml:"releases,omitempty"`
	}

	// unmarshal to auxiliary type
	if err := unmarshal(&aux); err != nil {
		return err
	}

	log.Debug().Interface("aux", aux).Msg("Unmarshalled auxiliary type for EstafetteManifest")

	// map auxiliary properties
	c.Builder = aux.Builder
	c.Version = aux.Version
	c.Labels = aux.Labels
	c.GlobalEnvVars = aux.GlobalEnvVars

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
		c.Stages = append(c.Stages, stage)
	}

	for _, mi := range aux.Releases {

		bytes, err := yaml.Marshal(mi.Value)
		if err != nil {
			return err
		}

		var release *EstafetteRelease
		if err := yaml.Unmarshal(bytes, &release); err != nil {
			return err
		}
		if release == nil {
			release = &EstafetteRelease{}
		}

		release.Name = mi.Key.(string)
		c.Releases = append(c.Releases, release)
	}

	log.Debug().Interface("manifest", c).Msg("Copied auxiliary type properties for EstafetteManifest")

	return nil
}

// // unmarshalYAML parses the .estafette.yaml file into an EstafetteManifest object
// func (c *EstafetteManifest) unmarshalYAML(data []byte) error {

// 	err := yaml.Unmarshal(data, c)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Unmarshalling .estafette.yaml manifest failed")
// 		return err
// 	}

// 	// set default for Builder.Track if not set
// 	if c.Builder.Track == "" {
// 		c.Builder.Track = "stable"
// 	}

// 	// set default version if no version is included
// 	if c.Version.Custom == nil && c.Version.SemVer == nil {
// 		c.Version.SemVer = &EstafetteSemverVersion{}
// 	}
// 	// if version is custom set defaults
// 	if c.Version.Custom != nil {
// 		if c.Version.Custom.LabelTemplate == "" {
// 			c.Version.Custom.LabelTemplate = "{{revision}}"
// 		}
// 	}

// 	// if version is semver set defaults
// 	if c.Version.SemVer != nil {
// 		if c.Version.SemVer.Patch == "" {
// 			c.Version.SemVer.Patch = "{{auto}}"
// 		}
// 		if c.Version.SemVer.LabelTemplate == "" {
// 			c.Version.SemVer.LabelTemplate = "{{branch}}"
// 		}
// 		if c.Version.SemVer.ReleaseBranch == "" {
// 			c.Version.SemVer.ReleaseBranch = "master"
// 		}
// 	}

// 	// create list of reserved property names
// 	reservedPropertyNames := getReservedPropertyNames()

// 	// to preserve order for the pipelines use MapSlice
// 	outerSlice := yaml.MapSlice{}
// 	err = yaml.Unmarshal(data, &outerSlice)
// 	if err != nil {
// 		return err
// 	}

// 	for _, s := range outerSlice {

// 		// support pipelines keyword for backwards compatibility, although it's replaced by the stages keyword
// 		if s.Key == "pipelines" || s.Key == "stages" {

// 			stages, err := unmarshalStagesYAML(s.Value, reservedPropertyNames)
// 			if err != nil {
// 				return err
// 			}

// 			c.Stages = stages
// 		}
// 		if s.Key == "releases" {

// 			releases, err := unmarshalReleasesYAML(s.Value, reservedPropertyNames)
// 			if err != nil {
// 				return err
// 			}

// 			c.Releases = releases
// 		}
// 	}

// 	return nil
// }

func getReservedPropertyNames() (names []string) {
	// create list of reserved property names
	reservedPropertyNames := []string{}
	val := reflect.ValueOf(EstafetteStage{})
	for i := 0; i < val.Type().NumField(); i++ {
		yamlName := val.Type().Field(i).Tag.Get("yaml")
		if yamlName != "" {
			reservedPropertyNames = append(reservedPropertyNames, strings.Replace(yamlName, ",omitempty", "", 1))
		}
		propertyName := val.Type().Field(i).Name
		if propertyName != "" {
			reservedPropertyNames = append(reservedPropertyNames, propertyName)
		}
	}

	return reservedPropertyNames
}

func isReservedPopertyName(s []string, e string) bool {

	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Exists checks whether the .estafette.yaml exists
func Exists(manifestPath string) bool {

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// does not exist
		return false
	}

	// does exist
	return true
}

// ReadManifestFromFile reads the .estafette.yaml into an EstafetteManifest object
func ReadManifestFromFile(manifestPath string) (manifest EstafetteManifest, err error) {

	log.Info().Msgf("Reading %v file...", manifestPath)

	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return manifest, err
	}

	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return manifest, err
	}

	log.Info().Msgf("Finished reading %v file successfully", manifestPath)

	return
}

// ReadManifest reads the string representation of .estafette.yaml into an EstafetteManifest object
func ReadManifest(manifestString string) (manifest EstafetteManifest, err error) {

	log.Info().Msg("Reading manifest from string...")

	if err := yaml.Unmarshal([]byte(manifestString), &manifest); err != nil {
		return manifest, err
	}

	log.Info().Msg("Finished unmarshalling manifest from string successfully")

	return
}
