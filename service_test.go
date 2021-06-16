package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestToYamlMarshalling(t *testing.T) {

	t.Run("ReturnsSameYamlLegacy", func(t *testing.T) {
		var service EstafetteService

		input := `name: kubernetes
image: bsycorp/kind:latest-1.15
env:
  SOME_ENVIRONMENT_VAR: some value with spaces
readiness:
  timeoutSeconds: 60
  path: /kubernetes-ready
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

	t.Run("ReturnsSameYamlHttpGet", func(t *testing.T) {
		var service EstafetteService

		input := `name: kubernetes
image: bsycorp/kind:latest-1.15
env:
  SOME_ENVIRONMENT_VAR: some value with spaces
readinessProbe:
  httpGet:
    path: /kubernetes-ready
    port: 80
    host: kubernetes.kube-system.svc.cluster.local
    scheme: http
  timeoutSeconds: 60
`
		// act
		err := yaml.Unmarshal([]byte(input), &service)
		assert.Nil(t, err)

		// act
		output, err := yaml.Marshal(service)

		assert.Nil(t, err)
		assert.Equal(t, input, string(output))
	})

	t.Run("ReturnsSameYamlExec", func(t *testing.T) {
		var service EstafetteService

		input := `name: kubernetes
image: bsycorp/kind:latest-1.15
env:
  SOME_ENVIRONMENT_VAR: some value with spaces
readinessProbe:
  exec:
    command:
    - /bin/sh
    - -c
    - -e
    - |
      exec pg_isready -U "postgres" -h 127.0.0.1 -p 5432
      [ -f /opt/bitnami/postgresql/tmp/.initialized ] || [ -f /bitnami/postgresql/.initialized ]
  timeoutSeconds: 60
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

func TestSetDefaults(t *testing.T) {

	t.Run("SetsShellToBinShIfEmpty", func(t *testing.T) {

		service := EstafetteService{
			Shell: "",
		}
		builder := EstafetteBuilder{}
		parentStage := EstafetteStage{}

		// act
		service.SetDefaults(builder, parentStage)

		assert.Equal(t, "/bin/sh", service.Shell)
	})

	t.Run("SetsShellToPowershellIfEmptyAndOperatingSystemIsWindows", func(t *testing.T) {

		service := EstafetteService{
			Shell: "",
		}
		builder := EstafetteBuilder{
			OperatingSystem: "windows",
		}
		parentStage := EstafetteStage{}

		// act
		service.SetDefaults(builder, parentStage)

		assert.Equal(t, "powershell", service.Shell)
	})

	t.Run("KeepsShellIfNotEmpty", func(t *testing.T) {

		service := EstafetteService{
			Shell: "/bin/bash",
		}
		builder := EstafetteBuilder{}
		parentStage := EstafetteStage{}

		// act
		service.SetDefaults(builder, parentStage)

		assert.Equal(t, "/bin/bash", service.Shell)
	})

	t.Run("SetsMultiStageToFalseIfNotSetAndParentStageHasImage", func(t *testing.T) {

		service := EstafetteService{
			Shell: "/bin/bash",
		}
		builder := EstafetteBuilder{}
		parentStage := EstafetteStage{
			ContainerImage: "alpine",
		}

		// act
		service.SetDefaults(builder, parentStage)

		assert.Equal(t, false, *service.MultiStage)
	})

	t.Run("SetsMultiStageToTrueIfNotSetAndParentStageHasNoImage", func(t *testing.T) {

		service := EstafetteService{
			Shell: "/bin/bash",
		}
		builder := EstafetteBuilder{}
		parentStage := EstafetteStage{
			ContainerImage: "",
		}

		// act
		service.SetDefaults(builder, parentStage)

		assert.Equal(t, true, *service.MultiStage)
	})

	t.Run("KeepsMultiStageIfSet", func(t *testing.T) {

		trueValue := true
		service := EstafetteService{
			Shell:      "/bin/bash",
			MultiStage: &trueValue,
		}
		builder := EstafetteBuilder{}
		parentStage := EstafetteStage{
			ContainerImage: "alpine",
		}

		// act
		service.SetDefaults(builder, parentStage)

		assert.Equal(t, true, *service.MultiStage)
	})
}
