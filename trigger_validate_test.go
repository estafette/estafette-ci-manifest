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

func TestEstafettePipelineTriggerValidate(t *testing.T) {
	t.Run("ReturnsErrorIfEventIsEmpty", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event:  "",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-manifest",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfStatusIsEmptyWhenEventIsFinished", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "",
			Name:   "github.com/estafette/estafette-ci-manifest",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfStatusIsEmptyWhenEventIsStarted", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/estafette/estafette-ci-manifest",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})

	t.Run("ReturnsErrorIfNameIsEmpty", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfValid", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-manifest",
			Branch: "master",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})
}

func TestEstafetteReleaseTriggerValidate(t *testing.T) {
	t.Run("ReturnsErrorIfEventIsEmpty", func(t *testing.T) {

		trigger := EstafetteReleaseTrigger{
			Event:  "",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-manifest",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfStatusIsEmptyWhenEventIsFinished", func(t *testing.T) {

		trigger := EstafetteReleaseTrigger{
			Event:  "finished",
			Status: "",
			Name:   "github.com/estafette/estafette-ci-manifest",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfStatusIsEmptyWhenEventIsStarted", func(t *testing.T) {

		trigger := EstafetteReleaseTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/estafette/estafette-ci-manifest",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
	})

	t.Run("ReturnsErrorIfNameIsEmpty", func(t *testing.T) {

		trigger := EstafetteReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsErrorIfTargetIsEmpty", func(t *testing.T) {

		trigger := EstafetteReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-manifest",
			Target: "",
		}

		// act
		err := trigger.Validate()

		assert.NotNil(t, err)
	})

	t.Run("ReturnsNoErrorIfValid", func(t *testing.T) {

		trigger := EstafetteReleaseTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-manifest",
			Target: "development",
		}

		// act
		err := trigger.Validate()

		assert.Nil(t, err)
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
