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

func TestEstafetteTriggerRunSetDefaults(t *testing.T) {
	t.Run("SetsBranchToBuildToMasterIfEmpty", func(t *testing.T) {

		trigger := EstafetteTriggerRun{
			BranchToBuild: "",
		}

		// act
		trigger.SetDefaults("build", "")

		assert.Equal(t, "master", trigger.BranchToBuild)
	})

	t.Run("SetsTriggerTypeToTriggerTypeParam", func(t *testing.T) {

		trigger := EstafetteTriggerRun{
			TriggerType: "any",
		}

		// act
		trigger.SetDefaults("build", "")

		assert.Equal(t, "build", trigger.TriggerType)
	})

	t.Run("SetsTargetNameToTargetNameParam", func(t *testing.T) {

		trigger := EstafetteTriggerRun{
			TargetName: "any",
		}

		// act
		trigger.SetDefaults("release", "development")

		assert.Equal(t, "development", trigger.TargetName)
	})
}
