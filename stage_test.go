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

		assert.Nil(t, err)
		assert.Equal(t, "/bin/sh", stage.Shell)
	})

	t.Run("DefaultsWhenToStatusEqualsSucceededIfNotPresent", func(t *testing.T) {

		var stage EstafetteStage

		// act
		err := yaml.Unmarshal([]byte(`
image: docker:17.03.0-ce
commands:
- cp Dockerfile ./publish
- docker build -t estafette-ci-builder ./publish`), &stage)

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

		assert.Nil(t, err)
		assert.Equal(t, "/estafette-work", stage.WorkingDirectory)
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

		// act
		bytes, err := json.Marshal(stage)

		if assert.Nil(t, err) {
			assert.Equal(t, "{\"Name\":\"\",\"ContainerImage\":\"extensions/gke:dev\",\"Shell\":\"/bin/sh\",\"WorkingDirectory\":\"/estafette-work\",\"Commands\":null,\"When\":\"status == 'succeeded'\",\"EnvVars\":null,\"AutoInjected\":false,\"Retries\":0,\"CustomProperties\":{\"container\":{\"repository\":\"extensions\"}}}", string(bytes))
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
			assert.Equal(t, "{\"Name\":\"\",\"ContainerImage\":\"docker:17.03.0-ce\",\"Shell\":\"/bin/bash\",\"WorkingDirectory\":\"/go/src/github.com/estafette/estafette-ci-manifest\",\"Commands\":[\"cp Dockerfile ./publish\",\"docker build -t estafette-ci-builder ./publish\"],\"When\":\"server == 'estafette'\",\"EnvVars\":null,\"AutoInjected\":false,\"Retries\":0,\"CustomProperties\":{}}", string(bytes))
		}
	})

}
