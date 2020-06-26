package main

import "time"

const GitLabConfigFilename = "gitlab/gitlab.conf.yaml"
const GitLabConfigFilenameLOCAL = "gitlab/gitlab.conf.yaml"
const DEFAULT_TIMEOUT = 600   // 600seconds == 10 minutes
const MAX_TIMEOUT = 3600      // 3600seconds == 1h
const DEFAULT_WAIT_RETRY = 10 // 10 seconds

const DEFAULT_INVOKE_RETRY_COUNT = 3 // Try Job Invoke 3 times
const DEFAULT_INVOKE_RETRY_WAIT = 2  // 2 Seconds wait between retries

/**
 * Defines the Jenkins Configuration File structure!
 */
type GitLabConfigFile struct {
	SpecVersion    string                  `json:"spec_version" yaml:"spec_version"`
	GitLabInstance []*GitLabInstanceConfig `json:"gitlabinstance" yaml:"gitlabinstance"`
	Pipelines      []*PipelineConfig       `json:"pipelines" yaml:"pipelines"`
	Events         []*EventMappingConfig   `json:"events" yaml:"events"`
}

/**
 * Defines the Jenkins Server Configuration
 */
type GitLabInstanceConfig struct {
	Name  string `json:"name" yaml:"name"`
	URL   string `json:"url" yaml:"url"`
	User  string `json:"user" yaml:"user"`
	Token string `json:"token" yaml:"token"`
}

/**
 * Defines Action Configuration, e.g: which Jenkins Job to execute on which Jenkins Server
 */
type PipelineConfig struct {
	Name      string            `json:"name" yaml:"name"`
	Instance  string            `json:"instance" yaml:"instance"`
	ProjectID string            `json:"projectid" yaml:"projectid"`
	Ref       string            `json:"ref" yaml:"ref"`
	Variables map[string]string `json:"variables" yaml:"variables"`
}

/**
 * Defines the Event Mapping from Keptn.Event to Action
 */
type EventMappingConfig struct {
	Event      string            `json:"event" yaml:"event"`
	Action     string            `json:"action" yaml:"action"`
	Timeout    int               `json:"timeout" yaml:"timoue"`
	OnSuccess  map[string]string `json:"onsuccess" yaml:"onsuccess"`
	OnFailure  map[string]string `json:"onfailure" yaml:"onfailure"`
	startedAt  time.Time
	finishedAt time.Time
}

const KeptnResultYaml = "keptn.result.yaml"

/**
 * This is when parsing the keptn.result.yaml build artifact from the Jenkins Job
 */
type KeptnResultArtifact struct {
	Data map[string]string `json:"data" yaml:"data"`
}

type DeploymentFinishedEventData_Extended struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// Service is the name of the new service
	Service string `json:"service"`
	// TestStrategy is the testing strategy
	TestStrategy string `json:"teststrategy"`
	// DeploymentStrategy is the deployment strategy
	DeploymentStrategy string `json:"deploymentstrategy"`
	// Tag of the new deployed artifact
	Tag string `json:"tag"`
	// Image of the new deployed artifact
	Image string `json:"image"`
	// Labels contains labels
	Labels map[string]string `json:"labels"`
	// DeploymentURILocal contains the local URL
	DeploymentURILocal string `json:"deploymentURILocal,omitempty"`
	// DeploymentURIPublic contains the public URL
	DeploymentURIPublic string `json:"deploymentURIPublic,omitempty"`
	// Result can be used to specify the status of the deployment, e.g: failed or success
	Result string `json:"result,omitempty"`
}
