package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstafettePipelineTriggerSetDefaults(t *testing.T) {
	t.Run("SetsEventToFinishedIfEmpty", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Event: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "finished", trigger.Event)
	})

	t.Run("SetsStatusToSucceededIfEmpty", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Status: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Status)
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

func TestEstafetteReleaseTriggerSetDefaults(t *testing.T) {
	t.Run("SetsEventToFinishedIfEmpty", func(t *testing.T) {

		trigger := EstafetteReleaseTrigger{
			Event: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "finished", trigger.Event)
	})
	t.Run("SetsStatusToSucceededIfEmpty", func(t *testing.T) {

		trigger := EstafetteReleaseTrigger{
			Status: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Status)
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

func TestEstafetteTriggerBuildActionSetDefaults(t *testing.T) {
	t.Run("SetsBranchToMasterIfEmpty", func(t *testing.T) {

		trigger := EstafetteTriggerBuildAction{
			Branch: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Branch)
	})
}

func TestEstafetteTriggerReleaseActionSetDefaults(t *testing.T) {
	t.Run("SetsTargetToTargetParam", func(t *testing.T) {

		trigger := EstafetteTriggerReleaseAction{
			Target: "any",
		}

		// act
		trigger.SetDefaults("development")

		assert.Equal(t, "development", trigger.Target)
	})

	t.Run("SetsVersionToLatestIfEmpty", func(t *testing.T) {

		trigger := EstafetteTriggerReleaseAction{
			Version: "",
		}

		// act
		trigger.SetDefaults("development")

		assert.Equal(t, "latest", trigger.Version)
	})

	t.Run("KeepsVersionIfNotEmpty", func(t *testing.T) {

		trigger := EstafetteTriggerReleaseAction{
			Version: "current",
		}

		// act
		trigger.SetDefaults("development")

		assert.Equal(t, "current", trigger.Version)
	})
}
