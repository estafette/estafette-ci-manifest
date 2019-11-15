package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestToYamlMarshalling(t *testing.T) {

	t.Run("ReturnsSameYaml", func(t *testing.T) {
		var service EstafetteService

		input := `name: kubernetes
image: bsycorp/kind:latest-1.15
env:
  SOME_ENVIRONMENT_VAR: some value with spaces
continueAfterStage: true
readiness:
  path: /kubernetes-ready
  timeoutSeconds: 60
  port: 80
  protocol: http
  hostname: kubernetes.kube-system.svc.cluster.local
`
		// act
		err := yaml.Unmarshal([]byte(input), &service)
		assert.Nil(t, err)

		// act
		output, err := yaml.Marshal(service)

		assert.Nil(t, err)
		assert.Equal(t, input, string(output))
	})
}
