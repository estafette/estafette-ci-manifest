package manifest

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/robfig/cron"
)

// EstafetteTrigger represents a trigger of any supported type and what action to take if the trigger fired
type EstafetteTrigger struct {
	Name      string                     `yaml:"name,omitempty" json:"name,omitempty"`
	Pipeline  *EstafettePipelineTrigger  `yaml:"pipeline,omitempty" json:"pipeline,omitempty"`
	Release   *EstafetteReleaseTrigger   `yaml:"release,omitempty" json:"release,omitempty"`
	Git       *EstafetteGitTrigger       `yaml:"git,omitempty" json:"git,omitempty"`
	Docker    *EstafetteDockerTrigger    `yaml:"docker,omitempty" json:"docker,omitempty"`
	Cron      *EstafetteCronTrigger      `yaml:"cron,omitempty" json:"cron,omitempty"`
	PubSub    *EstafettePubSubTrigger    `yaml:"pubsub,omitempty" json:"pubsub,omitempty"`
	Github    *EstafetteGithubTrigger    `yaml:"github,omitempty" json:"github,omitempty"`
	Bitbucket *EstafetteBitbucketTrigger `yaml:"bitbucket,omitempty" json:"bitbucket,omitempty"`

	BuildAction   *EstafetteTriggerBuildAction   `yaml:"builds,omitempty" json:"builds,omitempty"`
	ReleaseAction *EstafetteTriggerReleaseAction `yaml:"releases,omitempty" json:"releases,omitempty"`
	BotAction     *EstafetteTriggerBotAction     `yaml:"runs,omitempty" json:"runs,omitempty"`
}

// EstafettePipelineTrigger fires for pipeline changes and applies filtering to limit when this results in an action
type EstafettePipelineTrigger struct {
	Event  string `yaml:"event,omitempty" json:"event,omitempty"`
	Status string `yaml:"status,omitempty" json:"status,omitempty"`
	Name   string `yaml:"name,omitempty" json:"name,omitempty"`
	Branch string `yaml:"branch,omitempty" json:"branch,omitempty"`
}

// EstafetteReleaseTrigger fires for pipeline releases and applies filtering to limit when this results in an action
type EstafetteReleaseTrigger struct {
	Event  string `yaml:"event,omitempty" json:"event,omitempty"`
	Status string `yaml:"status,omitempty" json:"status,omitempty"`
	Name   string `yaml:"name,omitempty" json:"name,omitempty"`
	Target string `yaml:"target,omitempty" json:"target,omitempty"`
}

// EstafetteGitTrigger fires for git repository changes and applies filtering to limit when this results in an action
type EstafetteGitTrigger struct {
	Event      string `yaml:"event,omitempty" json:"event,omitempty"`
	Repository string `yaml:"repository,omitempty" json:"repository,omitempty"`
	Branch     string `yaml:"branch,omitempty" json:"branch,omitempty"`
}

// EstafetteDockerTrigger fires for docker image changes and applies filtering to limit when this results in an action
type EstafetteDockerTrigger struct {
	Event string `yaml:"event,omitempty" json:"event,omitempty"`
	Image string `yaml:"image,omitempty" json:"image,omitempty"`
	Tag   string `yaml:"tag,omitempty" json:"tag,omitempty"`
}

// EstafettePubSubTrigger fires for pubsub events in a certain project and topic
type EstafettePubSubTrigger struct {
	Project string `yaml:"project,omitempty" json:"project,omitempty"`
	Topic   string `yaml:"topic,omitempty" json:"topic,omitempty"`
}

// EstafetteGithubTrigger fires for github events
type EstafetteGithubTrigger struct {
	Events     []string `yaml:"events,omitempty" json:"events,omitempty"`
	Repository string   `yaml:"repository,omitempty" json:"repository,omitempty"`
}

// EstafetteBitbucketTrigger fires for bitbucket events
type EstafetteBitbucketTrigger struct {
	Events     []string `yaml:"events,omitempty" json:"events,omitempty"`
	Repository string   `yaml:"repository,omitempty" json:"repository,omitempty"`
}

// EstafetteCronTrigger fires at intervals specified by the cron schedule
type EstafetteCronTrigger struct {
	Schedule string `yaml:"schedule,omitempty" json:"schedule,omitempty"`
}

// EstafetteTriggerBuildAction determines what builds when the trigger fires
type EstafetteTriggerBuildAction struct {
	Branch string `yaml:"branch,omitempty" json:"branch,omitempty"`
}

// EstafetteTriggerReleaseAction determines what releases when the trigger fires
type EstafetteTriggerReleaseAction struct {
	Target  string `yaml:"target,omitempty" json:"target,omitempty"`
	Action  string `yaml:"action,omitempty" json:"action,omitempty"`
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
}

type EstafetteTriggerBotAction struct {
	Bot    string `yaml:"bot,omitempty" json:"bot,omitempty"`
	Branch string `yaml:"branch,omitempty" json:"branch,omitempty"`
}

// SetDefaults sets defaults for EstafetteTrigger
func (t *EstafetteTrigger) SetDefaults(preferences EstafetteManifestPreferences, triggerType TriggerType, targetName string) {
	if t.Pipeline != nil {
		t.Pipeline.SetDefaults()
	}
	if t.Release != nil {
		t.Release.SetDefaults()
	}
	if t.Git != nil {
		t.Git.SetDefaults()
	}
	if t.Docker != nil {
		t.Docker.SetDefaults()
	}
	if t.Cron != nil {
		t.Cron.SetDefaults()
	}
	if t.PubSub != nil {
		t.PubSub.SetDefaults()
	}
	if t.Github != nil {
		t.Github.SetDefaults()
	}
	if t.Bitbucket != nil {
		t.Bitbucket.SetDefaults()
	}

	switch triggerType {
	case TriggerTypeBuild:
		if t.BuildAction == nil {
			t.BuildAction = &EstafetteTriggerBuildAction{}
		}
		t.BuildAction.SetDefaults(preferences)
	case TriggerTypeRelease:
		if t.ReleaseAction == nil {
			t.ReleaseAction = &EstafetteTriggerReleaseAction{}
		}
		t.ReleaseAction.SetDefaults(t, targetName)
	case TriggerTypeBot:
		if t.BotAction == nil {
			t.BotAction = &EstafetteTriggerBotAction{}
		}
		t.BotAction.SetDefaults(preferences, targetName)
	}
}

// SetDefaults sets defaults for EstafettePipelineTrigger
func (p *EstafettePipelineTrigger) SetDefaults() {
	if p.Event == "" {
		p.Event = "finished"
	}
	if p.Status == "" {
		p.Status = "succeeded"
	}
	if p.Branch == "" {
		p.Branch = "master|main"
	}
}

// SetDefaults sets defaults for EstafetteReleaseTrigger
func (r *EstafetteReleaseTrigger) SetDefaults() {
	if r.Event == "" {
		r.Event = "finished"
	}
	if r.Status == "" {
		r.Status = "succeeded"
	}
}

// SetDefaults sets defaults for EstafetteGitTrigger
func (g *EstafetteGitTrigger) SetDefaults() {
	if g.Event == "" {
		g.Event = "push"
	}
	if g.Branch == "" {
		g.Branch = "master|main"
	}
}

// SetDefaults sets defaults for EstafetteDockerTrigger
func (d *EstafetteDockerTrigger) SetDefaults() {
}

// SetDefaults sets defaults for EstafetteCronTrigger
func (c *EstafetteCronTrigger) SetDefaults() {
}

// SetDefaults sets defaults for EstafettePubSubTrigger
func (p *EstafettePubSubTrigger) SetDefaults() {
}

// SetDefaults sets defaults for EstafetteGithubTrigger
func (p *EstafetteGithubTrigger) SetDefaults() {
	if p.Repository == "" {
		p.Repository = "self"
	}
}

// SetDefaults sets defaults for EstafetteBitbucketTrigger
func (p *EstafetteBitbucketTrigger) SetDefaults() {
	if p.Repository == "" {
		p.Repository = "self"
	}
}

// SetDefaults sets defaults for EstafetteTriggerBuildAction
func (b *EstafetteTriggerBuildAction) SetDefaults(preferences EstafetteManifestPreferences) {
	if b.Branch == "" {
		b.Branch = preferences.DefaultBranch
	}
}

// SetDefaults sets defaults for EstafetteTriggerReleaseAction
func (r *EstafetteTriggerReleaseAction) SetDefaults(t *EstafetteTrigger, targetName string) {
	r.Target = targetName
	if r.Version == "" {
		if t.Pipeline != nil && t.Pipeline.Name == "self" {
			r.Version = "same"
		} else if t.Release != nil && t.Release.Name == "self" {
			r.Version = "same"
		} else {
			r.Version = "latest"
		}
	}
}

// SetDefaults sets defaults for EstafetteTriggerReleaseAction
func (b *EstafetteTriggerBotAction) SetDefaults(preferences EstafetteManifestPreferences, botName string) {
	b.Bot = botName
	if b.Branch == "" {
		b.Branch = preferences.DefaultBranch
	}
}

// Validate checks if EstafetteTrigger is valid
func (t *EstafetteTrigger) Validate(triggerType TriggerType, targetName string) (err error) {

	numberOfTypes := 0

	if t.Pipeline == nil &&
		t.Release == nil &&
		t.Git == nil &&
		t.Docker == nil &&
		t.Cron == nil &&
		t.PubSub == nil &&
		t.Github == nil &&
		t.Bitbucket == nil {
		return fmt.Errorf("Set at least a 'pipeline', 'release', 'git', 'docker', 'cron', 'pubsub', 'github' or 'bitbucket' trigger")
	}

	if t.Pipeline != nil {
		err = t.Pipeline.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.Release != nil {
		err = t.Release.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.Git != nil {
		err = t.Git.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.Docker != nil {
		err = t.Docker.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.Cron != nil {
		err = t.Cron.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.PubSub != nil {
		err = t.PubSub.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.Github != nil {
		err = t.Github.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}
	if t.Bitbucket != nil {
		err = t.Bitbucket.Validate()
		if err != nil {
			return err
		}
		numberOfTypes++
	}

	if numberOfTypes != 1 {
		return fmt.Errorf("Do not specify more than one type of trigger 'pipeline', 'release', 'git', 'docker', 'cron', 'pubsub', 'github' or 'bitbucket' per trigger object")
	}

	switch triggerType {
	case TriggerTypeBuild:
		if t.BuildAction == nil {
			return fmt.Errorf("For a build trigger set the 'builds' property")
		}
		if t.ReleaseAction != nil {
			return fmt.Errorf("For a build trigger do not set the 'releases' property")
		}
		err = t.BuildAction.Validate()
		if err != nil {
			return err
		}
	case TriggerTypeRelease:
		if t.ReleaseAction == nil {
			return fmt.Errorf("For a release trigger set the 'releases' property")
		}
		if t.BuildAction != nil {
			return fmt.Errorf("For a release trigger do not set the 'builds' property")
		}
		err = t.ReleaseAction.Validate(targetName)
		if err != nil {
			return err
		}
	case TriggerTypeBot:
		if t.BotAction == nil {
			return fmt.Errorf("For a bot trigger set the 'runs' property")
		}
		err = t.BotAction.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate checks if EstafettePipelineTrigger is valid
func (p *EstafettePipelineTrigger) Validate() (err error) {
	if p.Event != "started" && p.Event != "finished" {
		return fmt.Errorf("Set pipeline.event in your trigger to 'started' or 'finished'")
	}
	if p.Event == "finished" && p.Status != "succeeded" && p.Status != "failed" {
		return fmt.Errorf("Set pipeline.status in your trigger to 'succeeded' or 'failed' for event 'finished'")
	}
	if p.Name == "" {
		return fmt.Errorf("Set pipeline.name in your trigger to 'self' or a full qualified pipeline name, i.e. github.com/estafette/estafette-ci-manifest")
	}
	return nil
}

// Validate checks if EstafetteReleaseTrigger is valid
func (r *EstafetteReleaseTrigger) Validate() (err error) {
	if r.Event != "started" && r.Event != "finished" {
		return fmt.Errorf("Set release.event in your trigger to 'started' or 'finished'")
	}
	if r.Event == "finished" && r.Status != "succeeded" && r.Status != "failed" {
		return fmt.Errorf("Set release.status in your trigger to 'succeeded' or 'failed' for event 'finished'")
	}
	if r.Name == "" {
		return fmt.Errorf("Set release.name in your trigger to 'self' or a full qualified pipeline name, i.e. github.com/estafette/estafette-ci-manifest")
	}
	if r.Target == "" {
		return fmt.Errorf("Set release.target in your trigger to a release target name on the pipeline set by release.name")
	}
	return nil
}

// Validate checks if EstafetteGitTrigger is valid
func (g *EstafetteGitTrigger) Validate() (err error) {
	if g.Event != "push" {
		return fmt.Errorf("Set git.event in your trigger to 'push'")
	}
	if g.Repository == "" {
		return fmt.Errorf("Set git.repository in your trigger to a full qualified git repository name, i.e. github.com/estafette/estafette-ci-manifest")
	}
	return nil
}

// Validate checks if EstafetteDockerTrigger is valid
func (d *EstafetteDockerTrigger) Validate() (err error) {
	return nil
}

// Validate checks if EstafetteCronTrigger is valid
func (c *EstafetteCronTrigger) Validate() (err error) {

	if c.Schedule == "" {
		return fmt.Errorf("Set cron.schedule in your trigger to '<minute> <hour> <day of month> <month> <day of week>'")
	}
	_, err = cron.ParseStandard(c.Schedule)
	if err != nil {
		return fmt.Errorf("Invalid cron.schedule in your trigger: %v", err)
	}

	return nil
}

// Validate checks if EstafettePubSubTrigger is valid
func (p *EstafettePubSubTrigger) Validate() (err error) {
	if p.Project == "" {
		return fmt.Errorf("Set pubsub.project in your trigger to the google cloud project id containing the pubsub topic")
	}
	if p.Topic == "" {
		return fmt.Errorf("Set pubsub.topic in your trigger to the pubsub topic you want this pipeline to subscribe to")
	}

	return nil
}

// Validate checks if EstafetteGithubTrigger is valid
func (p *EstafetteGithubTrigger) Validate() (err error) {
	if len(p.Events) == 0 {
		return fmt.Errorf("Set array github.events in your trigger to at least one github event")
	}

	return nil
}

// Validate checks if EstafetteBitbucketTrigger is valid
func (p *EstafetteBitbucketTrigger) Validate() (err error) {
	if len(p.Events) == 0 {
		return fmt.Errorf("Set array bitbucket.events in your trigger to at least one bitbucket event")
	}

	return nil
}

// Validate checks if EstafetteTriggerBuildAction is valid
func (b *EstafetteTriggerBuildAction) Validate() (err error) {
	return nil
}

// Validate checks if EstafetteTriggerReleaseAction is valid
func (r *EstafetteTriggerReleaseAction) Validate(targetName string) (err error) {
	if r.Target != targetName {
		return fmt.Errorf("The target in your releases action should have defaulted to '%v'", targetName)
	}

	return nil
}

// Validate checks if EstafetteTriggerBotAction is valid
func (b *EstafetteTriggerBotAction) Validate() (err error) {
	return nil
}

// ReplaceSelf replaces pipeline names set to "self" with the actual pipeline name
func (t *EstafetteTrigger) ReplaceSelf(pipeline string) {
	if t.Pipeline != nil && t.Pipeline.Name == "self" {
		t.Pipeline.Name = pipeline
	}
	if t.Release != nil && t.Release.Name == "self" {
		t.Release.Name = pipeline
	}
	if t.Github != nil && t.Github.Repository == "self" {
		t.Github.Repository = pipeline
	}
	if t.Bitbucket != nil && t.Bitbucket.Repository == "self" {
		t.Bitbucket.Repository = pipeline
	}
}

// Fires indicates whether EstafettePipelineTrigger fires for an EstafettePipelineEvent
func (p *EstafettePipelineTrigger) Fires(e *EstafettePipelineEvent) bool {

	// compare event as regex
	eventMatched, err := regexMatch(p.Event, e.Event)
	if err != nil || !eventMatched {
		return false
	}

	if p.Event == "finished" {
		// compare status as regex
		statusMatched, err := regexMatch(p.Status, e.Status)
		if err != nil || !statusMatched {
			return false
		}
	}

	// compare name case insensitive
	nameMatches := strings.EqualFold(p.Name, fmt.Sprintf("%v/%v/%v", e.RepoSource, e.RepoOwner, e.RepoName))
	if !nameMatches {
		return false
	}

	// compare branch as regex
	branchMatched, err := regexMatch(p.Branch, e.Branch)
	if err != nil || !branchMatched {
		return false
	}

	return true
}

// Fires indicates whether EstafetteReleaseTrigger fires for an EstafetteReleaseEvent
func (r *EstafetteReleaseTrigger) Fires(e *EstafetteReleaseEvent) bool {
	// compare event as regex
	eventMatched, err := regexMatch(r.Event, e.Event)
	if err != nil || !eventMatched {
		return false
	}

	if r.Event == "finished" {
		// compare status as regex
		statusMatched, err := regexMatch(r.Status, e.Status)
		if err != nil || !statusMatched {
			return false
		}
	}

	// compare name case insensitive
	nameMatches := strings.EqualFold(r.Name, fmt.Sprintf("%v/%v/%v", e.RepoSource, e.RepoOwner, e.RepoName))
	if !nameMatches {
		return false
	}

	// compare target as regex
	branchMatched, err := regexMatch(r.Target, e.Target)
	if err != nil || !branchMatched {
		return false
	}

	return true
}

// Fires indicates whether EstafetteGitTrigger fires for an EstafetteGitEvent
func (g *EstafetteGitTrigger) Fires(e *EstafetteGitEvent) bool {
	// compare event as regex
	eventMatched, err := regexMatch(g.Event, e.Event)
	if err != nil || !eventMatched {
		return false
	}

	// compare repository
	repositoryMatches := strings.EqualFold(g.Repository, e.Repository)
	if !repositoryMatches {
		return false
	}

	// compare branch as regex
	branchMatched, err := regexMatch(g.Branch, e.Branch)
	if err != nil || !branchMatched {
		return false
	}

	return true
}

// Fires indicates whether EstafetteDockerTrigger fires for an EstafetteDockerEvent
func (d *EstafetteDockerTrigger) Fires(e *EstafetteDockerEvent) bool {
	return false
}

// Fires indicates whether EstafetteCronTrigger fires for an EstafetteCronEvent
func (c *EstafetteCronTrigger) Fires(e *EstafetteCronEvent) bool {

	// ParseStandard expects 5 entries representing: minute, hour, day of month, month and day of week, in that order.
	sched, err := cron.ParseStandard(c.Schedule)
	if err != nil {
		return false
	}

	// truncate event time to the minute
	eventTime := time.Date(e.Time.Year(), e.Time.Month(), e.Time.Day(), e.Time.Hour(), e.Time.Minute(), 0, 0, time.UTC)
	// subtract 1 minute, otherwise the next time is at least 1 minute later
	eventTime = eventTime.Add(time.Minute * -1)
	// get the next time the cron expression would fire
	nextTime := sched.Next(eventTime)

	return nextTime.Year() == e.Time.Year() &&
		nextTime.Month() == e.Time.Month() &&
		nextTime.Day() == e.Time.Day() &&
		nextTime.Hour() == e.Time.Hour() &&
		nextTime.Minute() == e.Time.Minute()
}

func regexMatch(pattern, value string) (bool, error) {

	// check to see if the pattern starts with any of the promql regex operators, so we can do negations
	// =~ : Select labels that regex-match the provided string (or substring).
	// !~ : Select labels that do not regex-match the provided string (or substring).

	negativeMatching := false
	if strings.HasPrefix(pattern, "=~") {
		negativeMatching = false
		pattern = strings.TrimPrefix(pattern, "=~")
	} else if strings.HasPrefix(pattern, "!~") {
		negativeMatching = true
		pattern = strings.TrimPrefix(pattern, "!~")
	}

	pattern = fmt.Sprintf("^(%v)$", strings.TrimSpace(pattern))

	match, err := regexp.MatchString(pattern, value)

	if err != nil {
		return false, err
	}

	if negativeMatching {
		return !match, nil
	}

	return match, nil
}

// Fires indicates whether EstafettePubSubTrigger fires for an EstafettePubSubEvent
func (p *EstafettePubSubTrigger) Fires(e *EstafettePubSubEvent) bool {

	projectMatch, err := regexp.MatchString(fmt.Sprintf("^(%v)$", strings.TrimSpace(p.Project)), e.Project)
	if !projectMatch || err != nil {
		return false
	}

	topicMatch, err := regexp.MatchString(fmt.Sprintf("^(%v)$", strings.TrimSpace(p.Topic)), e.Topic)
	if !topicMatch || err != nil {
		return false
	}

	return true
}

// Fires indicates whether EstafetteGithubTrigger fires for an EstafetteGithubEvent
func (p *EstafetteGithubTrigger) Fires(e *EstafetteGithubEvent) bool {
	if e.Repository != "" && e.Repository != p.Repository {
		return false
	}

	for _, ev := range p.Events {
		if ev == e.Event {
			return true
		}
	}

	return true
}

// Fires indicates whether EstafetteBitbucketTrigger fires for an EstafetteBitbucketEvent
func (p *EstafetteBitbucketTrigger) Fires(e *EstafetteBitbucketEvent) bool {
	if e.Repository != "" && e.Repository != p.Repository {
		return false
	}

	for _, ev := range p.Events {
		if ev == e.Event {
			return true
		}
	}

	return true
}
