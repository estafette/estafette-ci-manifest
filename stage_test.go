package manifest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestUnmarshalStage(t *testing.T) {
	t.Run("ReturnsUnmarshaledStage", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
shell: /bin/bash
workDir: /go/src/github.com/estafette/estafette-ci-manifest
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish
when:
  server == 'estafette'`), &stage)

		assert.Nil(t, err)
		assert.Equal(t, "docker:17.03.0-ce", stage.ContainerImage)
		assert.Equal(t, "/bin/bash", stage.Shell)
		assert.Equal(t, "/go/src/github.com/estafette/estafette-ci-manifest", stage.WorkingDirectory)
		assert.Equal(t, 2, len(stage.Commands))
		assert.Equal(t, "cp Dockerfile ./publish", stage.Commands[0])
		assert.Equal(t, "docker build -t estafette-ci-builder ./publish", stage.Commands[1])
		assert.Equal(t, "server == 'estafette'", stage.When)
	})

	t.Run("DefaultsShellToShIfNotPresent", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish
when:
  server == 'estafette'`), &stage)

		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
		})

		assert.Nil(t, err)
		assert.Equal(t, "/bin/sh", stage.Shell)
	})

	t.Run("DefaultsShellToPowershellIfNotPresentForWindows", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish
when:
  server == 'estafette'`), &stage)

		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "windows",
		})

		assert.Nil(t, err)
		assert.Equal(t, "powershell", stage.Shell)
	})

	t.Run("DefaultsWhenToStatusEqualsSucceededIfNotPresent", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish`), &stage)

		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
		})

		assert.Nil(t, err)
		assert.Equal(t, "status == 'succeeded'", stage.When)
	})

	t.Run("DefaultsWorkingDirectoryToEstafetteWorkIfNotPresent", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish`), &stage)

		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
		})

		assert.Nil(t, err)
		assert.Equal(t, "/estafette-work", stage.WorkingDirectory)
	})

	t.Run("DefaultsWorkingDirectoryToEstafetteWorkIfNotPresentForWindows", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish`), &stage)

		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "windows",
		})

		assert.Nil(t, err)
		assert.Equal(t, "C:/estafette-work", stage.WorkingDirectory)
	})

	t.Run("ReturnsNonReservedSimplePropertyAsCustomProperty", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
unknownProperty1: value1
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish`), &stage)

		assert.Nil(t, err)
		assert.Equal(t, "value1", stage.CustomProperties["unknownProperty1"])
	})

	t.Run("ReturnsNonReservedArrayPropertyAsCustomProperty", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
unknownProperty3:
- supported1
- supported2
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish`), &stage)

		assert.Nil(t, err)
		assert.NotNil(t, stage.CustomProperties["unknownProperty3"])
		assert.Equal(t, "supported1", stage.CustomProperties["unknownProperty3"].([]interface{})[0].(string))
		assert.Equal(t, "supported2", stage.CustomProperties["unknownProperty3"].([]interface{})[1].(string))
	})
}

func TestJSONMarshalStage(t *testing.T) {
	t.Run("MarshalMapStringInterface", func(t *testing.T) {

		property := map[string]interface{}{
			"container": map[string]interface{}{
				"repository": "extension",
				"name":       "gke",
				"tag":        "alpha",
			},
		}

		// act
		bytes, err := json.Marshal(property)

		if assert.Nil(t, err) {
			assert.Equal(t, "{\"container\":{\"name\":\"gke\",\"repository\":\"extension\",\"tag\":\"alpha\"}}", string(bytes))
		}
	})

	t.Run("ReturnsMarshaledStageForNestedCustomProperties", func(t *testing.T) {

		var stage EstafetteStage

		err := yaml.Unmarshal([]byte(`
image: extensions/gke:dev
container:
  repository: extensions`), &stage)

		assert.Nil(t, err)

		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
		})

		// act
		bytes, err := json.Marshal(stage)

		if assert.Nil(t, err) {
			assert.Equal(t, "{\"ContainerImage\":\"extensions/gke:dev\",\"Shell\":\"/bin/sh\",\"WorkingDirectory\":\"/estafette-work\",\"When\":\"status == 'succeeded'\",\"CustomProperties\":{\"container\":{\"repository\":\"extensions\"}}}", string(bytes))
		}
	})

	t.Run("ReturnsMarshaledStage", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
shell: /bin/bash
workDir: /go/src/github.com/estafette/estafette-ci-manifest
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish
when:
  server == 'estafette'`), &stage)

		// act
		bytes, err := json.Marshal(stage)

		if assert.Nil(t, err) {
			assert.Equal(t, "{\"ContainerImage\":\"docker:17.03.0-ce\",\"Shell\":\"/bin/bash\",\"WorkingDirectory\":\"/go/src/github.com/estafette/estafette-ci-manifest\",\"Commands\":[\"cp Dockerfile ./publish\",\"docker build -t estafette-ci-builder ./publish\"],\"When\":\"server == 'estafette'\"}", string(bytes))
		}
	})

}

func TestValidateOnStage(t *testing.T) {
	t.Run("ReturnsErrorIfImageAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := EstafetteStage{
			ContainerImage: "docker",
			ParallelStages: []*EstafetteStage{
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfShellAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := EstafetteStage{
			Shell: "/bin/sh",
			ParallelStages: []*EstafetteStage{
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfWorkingDirectoryAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := EstafetteStage{
			WorkingDirectory: "/estafette-work",
			ParallelStages: []*EstafetteStage{
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfCommandsAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := EstafetteStage{
			Commands: []string{"dotnet build"},
			ParallelStages: []*EstafetteStage{
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfEnvvarsAndParallelStagesAreBothSet", func(t *testing.T) {

		stage := EstafetteStage{
			EnvVars: map[string]string{
				"ENVA": "value a",
			},
			ParallelStages: []*EstafetteStage{
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfImageIsNotSetWithoutParallelStages", func(t *testing.T) {

		stage := EstafetteStage{
			ContainerImage: "",
		}
		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfImageIsNotSetButHasService", func(t *testing.T) {

		stage := EstafetteStage{
			ContainerImage: "",
			Services: []*EstafetteService{
				&EstafetteService{
					Name:           "cockroachdb",
					ContainerImage: "cockroachdb/cockroach:v19.2.0",
				},
			},
		}
		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.Nil(t, err)
	})

	t.Run("ReturnsNoErrorWhenAllFieldsAreValid", func(t *testing.T) {

		stage := EstafetteStage{
			ContainerImage: "docker",
		}
		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.Nil(t, err)
	})

	t.Run("ReturnsNoErrorWhenAllParallelStagesAreValid", func(t *testing.T) {

		stage := EstafetteStage{
			ParallelStages: []*EstafetteStage{
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageA",
				},
				&EstafetteStage{
					ContainerImage: "docker",
					Name:           "StageB",
				},
			},
		}
		stage.SetDefaults(EstafetteBuilder{
			OperatingSystem: "linux",
			Track:           "stable",
		})

		// act
		err := stage.Validate()

		assert.Nil(t, err)
	})
}
