package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstafetteTriggerPipelineFires(t *testing.T) {
	t.Run("ReturnsTrueIfEventStatusNameAndBranchMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
		}

		// act
		fires, err := trigger.Fires(&event)

		assert.Nil(t, err)
		assert.True(t, fires)
	})

	t.Run("ReturnsFalseIfEventDoesNotMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "started",
			Status: "",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
		}

		// act
		fires, err := trigger.Fires(&event)

		assert.Nil(t, err)
		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfStatusDoesNotMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "failed",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
		}

		// act
		fires, err := trigger.Fires(&event)

		assert.Nil(t, err)
		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfNameDoesNotMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-builder",
			Branch: "master",
		}

		// act
		fires, err := trigger.Fires(&event)

		assert.Nil(t, err)
		assert.False(t, fires)
	})

	t.Run("ReturnsFalseIfBranchDoesNotMatch", func(t *testing.T) {

		event := EstafettePipelineEvent{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "master",
		}

		trigger := EstafettePipelineTrigger{
			Event:  "finished",
			Status: "succeeded",
			Name:   "github.com/estafette/estafette-ci-api",
			Branch: "development",
		}

		// act
		fires, err := trigger.Fires(&event)

		assert.Nil(t, err)
		assert.False(t, fires)
	})
}
