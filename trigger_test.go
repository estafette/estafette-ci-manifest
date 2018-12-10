package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDefaults(t *testing.T) {

	t.Run("DefaultsFilterBranchToMasterForEventPipelineBuildStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Filter.Branch)
	})

	t.Run("DefaultsFilterBranchToMasterForEventPipelineBuildFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Filter.Branch)
	})

	t.Run("DefaultsFilterStatusToSucceededForEventPipelineBuildFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Filter.Status)
	})

	t.Run("DefaultsFilterBranchToAnyForEventPipelineReleaseStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, ".+", trigger.Filter.Branch)
	})

	t.Run("DefaultsFilterBranchToAnyForEventPipelineReleaseFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, ".+", trigger.Filter.Branch)
	})

	t.Run("DefaultsFilterStatusToSucceededForEventPipelineReleaseFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "succeeded", trigger.Filter.Status)
	})

	t.Run("DefaultsFilterPipelineToThisPipelineToAnyForEventPipelineReleaseStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "this", trigger.Filter.Pipeline)
	})

	t.Run("DefaultsFilterPipelineToThisPipelineToAnyForEventPipelineReleaseFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "this", trigger.Filter.Pipeline)
	})

	t.Run("DefaultsThenBranchToMasterForEventPipelineBuildStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Run.Branch)
	})

	t.Run("DefaultsThenBranchToMasterForEventPipelineBuildFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-build-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, "master", trigger.Run.Branch)
	})

	t.Run("DefaultsThenBranchToAnyForEventPipelineReleaseStarted", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-started",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, ".+", trigger.Run.Branch)
	})

	t.Run("DefaultsThenBranchToAnyForEventPipelineReleaseFinished", func(t *testing.T) {

		trigger := EstafetteTrigger{
			Event: "pipeline-release-finished",
		}

		// act
		trigger.SetDefaults()

		assert.Equal(t, ".+", trigger.Run.Branch)
	})
}

func TestRunEquals(t *testing.T) {

	t.Run("ReturnsTrueIfBothRunObjectsAreNil", func(t *testing.T) {

		var run1 *EstafetteTriggerRun
		var run2 *EstafetteTriggerRun

		// act
		equal := run1.Equals(run2)

		assert.True(t, equal)
	})

	t.Run("ReturnsFalseIfOneIsNilAndTheOtherIsNot", func(t *testing.T) {

		var run1 *EstafetteTriggerRun
		run2 := &EstafetteTriggerRun{}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfOneIsNilAndTheOtherIsNotReversed", func(t *testing.T) {

		run1 := &EstafetteTriggerRun{}
		var run2 *EstafetteTriggerRun

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfBranchIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerRun{
			Branch: "master",
		}
		run2 := &EstafetteTriggerRun{
			Branch: "development",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfActionIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerRun{
			Action: "deploy-canary",
		}
		run2 := &EstafetteTriggerRun{
			Action: "deploy-stable",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsTrueIfAllOptionsAreEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerRun{
			Branch: "master",
			Action: "deploy-canary",
		}
		run2 := &EstafetteTriggerRun{
			Branch: "master",
			Action: "deploy-canary",
		}

		// act
		equal := run1.Equals(run2)

		assert.True(t, equal)
	})
}

func TestFilterEquals(t *testing.T) {

	t.Run("ReturnsTrueIfBothRunObjectsAreNil", func(t *testing.T) {

		var run1 *EstafetteTriggerFilter
		var run2 *EstafetteTriggerFilter

		// act
		equal := run1.Equals(run2)

		assert.True(t, equal)
	})

	t.Run("ReturnsFalseIfOneIsNilAndTheOtherIsNot", func(t *testing.T) {

		var run1 *EstafetteTriggerFilter
		run2 := &EstafetteTriggerFilter{}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfOneIsNilAndTheOtherIsNotReversed", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{}
		var run2 *EstafetteTriggerFilter

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfPipelineIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{
			Pipeline: "this",
		}
		run2 := &EstafetteTriggerFilter{
			Pipeline: "github.com/estafette/estafette-extensions-gke",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfBranchIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{
			Branch: "master",
		}
		run2 := &EstafetteTriggerFilter{
			Branch: "development",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfActionIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{
			Action: "deploy-canary",
		}
		run2 := &EstafetteTriggerFilter{
			Action: "deploy-stable",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfTargetIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{
			Target: "beta",
		}
		run2 := &EstafetteTriggerFilter{
			Target: "stable",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfStatusIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{
			Status: "succeeded",
		}
		run2 := &EstafetteTriggerFilter{
			Status: "failed",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfCronIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{
			Cron: "*/5 * * * *",
		}
		run2 := &EstafetteTriggerFilter{
			Cron: "*/6 * * * *",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfImageIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{
			Image: "golang",
		}
		run2 := &EstafetteTriggerFilter{
			Image: "node",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfTagIsNotEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{
			Tag: "1.11.2-alpine3.8",
		}
		run2 := &EstafetteTriggerFilter{
			Tag: "1.11-alpine3.8",
		}

		// act
		equal := run1.Equals(run2)

		assert.False(t, equal)
	})

	t.Run("ReturnsTrueIfAllOptionsAreEqual", func(t *testing.T) {

		run1 := &EstafetteTriggerFilter{
			Pipeline: "github.com/estafette/estafette-extensions-gke",
			Branch:   "master",
			Action:   "deploy-canary",
			Target:   "stable",
			Status:   "succeeded",
			Cron:     "*/5 * * * *",
			Image:    "golang",
			Tag:      "1.11.2-alpine3.8",
		}
		run2 := &EstafetteTriggerFilter{
			Pipeline: "github.com/estafette/estafette-extensions-gke",
			Branch:   "master",
			Action:   "deploy-canary",
			Target:   "stable",
			Status:   "succeeded",
			Cron:     "*/5 * * * *",
			Image:    "golang",
			Tag:      "1.11.2-alpine3.8",
		}

		// act
		equal := run1.Equals(run2)

		assert.True(t, equal)
	})
}

func TestTriggerEquals(t *testing.T) {

	t.Run("ReturnsTrueIfBothRunObjectsAreNil", func(t *testing.T) {

		var trigger1 *EstafetteTrigger
		var trigger2 *EstafetteTrigger

		// act
		equal := trigger1.Equals(trigger2)

		assert.True(t, equal)
	})

	t.Run("ReturnsFalseIfOneIsNilAndTheOtherIsNot", func(t *testing.T) {

		var trigger1 *EstafetteTrigger
		trigger2 := &EstafetteTrigger{}

		// act
		equal := trigger1.Equals(trigger2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfOneIsNilAndTheOtherIsNotReversed", func(t *testing.T) {

		trigger1 := &EstafetteTrigger{}
		var trigger2 *EstafetteTrigger

		// act
		equal := trigger1.Equals(trigger2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfEventIsNotEqual", func(t *testing.T) {

		trigger1 := &EstafetteTrigger{
			Event: "pipeline-build-finished",
		}
		trigger2 := &EstafetteTrigger{
			Event: "pipeline-build-started",
		}

		// act
		equal := trigger1.Equals(trigger2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfFilterIsNotEqual", func(t *testing.T) {

		trigger1 := &EstafetteTrigger{
			Filter: &EstafetteTriggerFilter{
				Pipeline: "this",
			},
		}
		trigger2 := &EstafetteTrigger{
			Filter: &EstafetteTriggerFilter{
				Pipeline: "github.com/estafette/estafette-extension-gke",
			},
		}

		// act
		equal := trigger1.Equals(trigger2)

		assert.False(t, equal)
	})

	t.Run("ReturnsFalseIfRunIsNotEqual", func(t *testing.T) {

		trigger1 := &EstafetteTrigger{
			Run: &EstafetteTriggerRun{
				Branch: "master",
			},
		}
		trigger2 := &EstafetteTrigger{
			Run: &EstafetteTriggerRun{
				Branch: "development",
			},
		}

		// act
		equal := trigger1.Equals(trigger2)

		assert.False(t, equal)
	})

	t.Run("ReturnsTrueIfAllOptionsAreEqual", func(t *testing.T) {

		trigger1 := &EstafetteTrigger{
			Event: "pipeline-build-finished",
			Filter: &EstafetteTriggerFilter{
				Pipeline: "github.com/estafette/estafette-extensions-gke",
				Branch:   "master",
				Action:   "deploy-canary",
				Target:   "stable",
				Status:   "succeeded",
				Cron:     "*/5 * * * *",
				Image:    "golang",
				Tag:      "1.11.2-alpine3.8",
			},
			Run: &EstafetteTriggerRun{
				Branch: "master",
				Action: "deploy-canary",
			},
		}
		trigger2 := &EstafetteTrigger{
			Event: "pipeline-build-finished",
			Filter: &EstafetteTriggerFilter{
				Pipeline: "github.com/estafette/estafette-extensions-gke",
				Branch:   "master",
				Action:   "deploy-canary",
				Target:   "stable",
				Status:   "succeeded",
				Cron:     "*/5 * * * *",
				Image:    "golang",
				Tag:      "1.11.2-alpine3.8",
			},
			Run: &EstafetteTriggerRun{
				Branch: "master",
				Action: "deploy-canary",
			},
		}

		// act
		equal := trigger1.Equals(trigger2)

		assert.True(t, equal)
	})

}
