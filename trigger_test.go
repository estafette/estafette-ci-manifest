package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstafetteTriggerValidate(t *testing.T) {
	t.Run("ReturnsErrorIfAllTypesAreEmpty", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Pipeline: nil,
			Git:      nil,
			Docker:   nil,
			Cron:     nil,
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfMoreThanOneTypeIsSet", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Pipeline: &EstafettePipelineTrigger{
				Event:  "succeeded",
				Name:   "github.com/estafette/estafette-ci-api",
				Branch: "master",
			},
			Git: &EstafetteGitTrigger{
				Event:      "push",
				Repository: "github.com/estafette/estafette-ci-builder",
				Branch:     "master",
			},
			Docker: nil,
			Cron:   nil,
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})
}

func TestEstafettePipelineTriggerSetDefaults(t *testing.T) {
	t.Run("SetsEventToSucceededIfEmpty", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Event)
	})

	t.Run("SetsBranchToMasterIfEmpty", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Branch: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Branch)
	})
}

func TestEstafettePipelineTriggerValidate(t *testing.T) {
	t.Run("ReturnsErrorIfEventIsEmpty", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event: "",
			Name:  "github.com/estafette/estafette-ci-manifest",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfNameIsEmpty", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event: "succeeded",
			Name:  "",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfValid", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event: "succeeded",
			Name:  "github.com/estafette/estafette-ci-manifest",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})
}

func TestEstafetteGitTriggerSetDefaults(t *testing.T) {
	t.Run("SetsEventToPushIfEmpty", func(t *testing.T) {

		trigger := EstafetteGitTrigger{
			Event: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "push", trigger.Event)
	})

	t.Run("SetsBranchToMasterIfEmpty", func(t *testing.T) {

		trigger := EstafetteGitTrigger{
			Branch: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Branch)
	})
}

func TestEstafetteGitTriggerValidate(t *testing.T) {
	t.Run("ReturnsErrorIfEventIsEmpty", func(t *testing.T) {

		trigger := EstafetteGitTrigger{
			Event:      "",
			Repository: "github.com/estafette/estafette-ci-manifest",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfRepositoryIsEmpty", func(t *testing.T) {

		trigger := EstafetteGitTrigger{
			Event:      "push",
			Repository: "",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfValid", func(t *testing.T) {

		trigger := EstafetteGitTrigger{
			Event:      "push",
			Repository: "github.com/estafette/estafette-ci-manifest",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})
}

func TestEstafetteTriggerRunSetDefaults(t *testing.T) {
	t.Run("SetsStatusToSucceededIfEmpty", func(t *testing.T) {

		trigger := EstafetteTriggerRun{
			Status: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Status)
	})

	t.Run("SetsBranchToMasterIfEmpty", func(t *testing.T) {

		trigger := EstafetteTriggerRun{
			Branch: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Branch)
	})
}
