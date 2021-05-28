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

	t.Run("SetsBranchToMasterOrMainIfEmpty", func(t *testing.T) {

		trigger := EstafettePipelineTrigger{
			Branch: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master|main", trigger.Branch)
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

	t.Run("SetsBranchToMasterOrMainIfEmpty", func(t *testing.T) {

		trigger := EstafetteGitTrigger{
			Branch: "",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master|main", trigger.Branch)
	})
}

func TestEstafetteTriggerBuildActionSetDefaults(t *testing.T) {
	t.Run("SetsBranchToMasterIfEmpty", func(t *testing.T) {

		trigger := EstafetteTriggerBuildAction{
			Branch: "",
		}

		preferences := EstafetteManifestPreferences{
			DefaultBranch: "main",
		}

		// act
		trigger.SetDefaults(preferences)

		assert.Equal(t, "main", trigger.Branch)
	})
}

func TestEstafetteTriggerReleaseActionSetDefaults(t *testing.T) {
	t.Run("SetsTargetToTargetParam", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Pipeline: &EstafettePipelineTrigger{
				Name: "self",
			},
			ReleaseAction: &EstafetteTriggerReleaseAction{
				Target: "any",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "development", trigger.ReleaseAction.Target)
	})

	t.Run("SetsVersionToLatestIfEmptyAndTriggersOnOtherPipeline", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Pipeline: &EstafettePipelineTrigger{
				Name: "github.com/estafette/estafette-ci-builder",
			},
			ReleaseAction: &EstafetteTriggerReleaseAction{
				Target:  "any",
				Version: "",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "latest", trigger.ReleaseAction.Version)
	})

	t.Run("KeepsVersionIfNotEmpty", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Pipeline: &EstafettePipelineTrigger{
				Name: "self",
			},
			ReleaseAction: &EstafetteTriggerReleaseAction{
				Target:  "any",
				Version: "current",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "current", trigger.ReleaseAction.Version)
	})

	t.Run("SetsVersionToSameIfPipelineTriggerIsTheSelfPipeline", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Pipeline: &EstafettePipelineTrigger{
				Name: "self",
			},
			ReleaseAction: &EstafetteTriggerReleaseAction{
				Target:  "development",
				Version: "",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "same", trigger.ReleaseAction.Version)
	})

	t.Run("SetsVersionToSameIfReleaseTriggerIsTheSelfPipeline", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Release: &EstafetteReleaseTrigger{
				Name: "self",
			},
			ReleaseAction: &EstafetteTriggerReleaseAction{
				Target:  "development",
				Version: "",
			},
		}

		// act
		trigger.ReleaseAction.SetDefaults(&trigger, "development")

		assert.Equal(t, "same", trigger.ReleaseAction.Version)
	})

}
