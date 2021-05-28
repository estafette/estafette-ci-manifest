package manifest

type TriggerType string

const (
	TriggerTypeUnknown TriggerType = ""
	TriggerTypeBuild   TriggerType = "build"
	TriggerTypeRelease TriggerType = "release"
	TriggerTypeBot     TriggerType = "bot"
)
