package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/ghodss/yaml"

	configutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnutils "github.com/keptn/go-utils/pkg/lib"

	"github.com/xanzy/go-gitlab"
)

func makeJson(data interface{}) string {
	str, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(json.RawMessage(str))
}

//
// Loads gitlab.conf for the current service
//
func getGitLabConfiguration(keptnEvent baseKeptnEvent, logger *keptnutils.Logger) (*GitLabConfigFile, error) {

	// if we run in a runlocal mode we are just getting the file from the local disk
	var fileContent []byte
	var err error
	if runlocal {
		fileContent, err = ioutil.ReadFile(GitLabConfigFilenameLOCAL)
		if err != nil {
			logMessage := fmt.Sprintf("No %s file found LOCALLY for service %s in stage %s in project %s", GitLabConfigFilenameLOCAL, keptnEvent.service, keptnEvent.stage, keptnEvent.project)
			logger.Info(logMessage)
			return nil, errors.New(logMessage)
		}
	} else {
		resourceHandler := configutils.NewResourceHandler("configuration-service:8080")
		//resourceHandler := configutils.NewAuthenticatedResourceHandler(endPoint.Host+"/configuration-service", apiToken, "x-token", nil, "http")
		keptnResourceContent, err := resourceHandler.GetServiceResource(keptnEvent.project, keptnEvent.stage, keptnEvent.service, GitLabConfigFilename)
		if err != nil {
			logMessage := fmt.Sprintf("No %s file found for service %s in stage %s in project %s", GitLabConfigFilename, keptnEvent.service, keptnEvent.stage, keptnEvent.project)
			fmt.Println(err)
			logger.Info(logMessage)
			return nil, errors.New(logMessage)
		}
		fileContent = []byte(keptnResourceContent.ResourceContent)
	}

	// replace the placeholders
	fileAsString := string(fileContent)
	fileAsStringReplaced := replaceKeptnPlaceholders(fileAsString, keptnEvent)
	fileContent = []byte(fileAsStringReplaced)

	// Print file content
	// fmt.Println(fileAsStringReplaced)

	// unmarshal the file
	var gitLabConfFile *GitLabConfigFile
	gitLabConfFile, err = parseGitLabConfigFile(fileContent)

	if err != nil {
		logMessage := fmt.Sprintf("Couldn't parse %s file found for service %s in stage %s in project %s. Error: %s", GitLabConfigFilename, keptnEvent.service, keptnEvent.stage, keptnEvent.project, err.Error())
		logger.Error(logMessage)
		return nil, errors.New(logMessage)
	}

	return gitLabConfFile, nil
}

func parseGitLabConfigFile(input []byte) (*GitLabConfigFile, error) {
	gitlabConfFile := &GitLabConfigFile{}
	err := yaml.Unmarshal([]byte(input), &gitlabConfFile)

	if err != nil {
		return nil, err
	}

	return gitlabConfFile, nil
}

/**
 * Iterates through the events and returns the first matching the event name
 */
func getEventMappingConfig(gitlabConfigFile *GitLabConfigFile, eventname string) (*EventMappingConfig, error) {
	// get the entry for the passed name
	if gitlabConfigFile != nil && gitlabConfigFile.Events != nil {
		for _, eventMapConfig := range gitlabConfigFile.Events {
			if eventMapConfig.Event == eventname {
				return eventMapConfig, nil
			}
		}
	}

	return nil, errors.New("No Event Map Configuration found for " + eventname)
}

/**
 * Iterates through the GitLabConfigFile and returns the pipeline configuration by name
 */
func getPipelineConfig(gitlabConfigFile *GitLabConfigFile, eventname string) (*PipelineConfig, error) {
	// get the entry for the passed name
	if gitlabConfigFile != nil && gitlabConfigFile.Pipelines != nil {
		for _, pipelineConfig := range gitlabConfigFile.Pipelines {
			if pipelineConfig.Name == eventname {
				return pipelineConfig, nil
			}
		}
	}

	return nil, errors.New("No GitLab Pipeline Configuration found for " + eventname)
}

/**
 * Iterates through the GitLabConfigFile and returns the gitlab instance configuration by name
 */
func getGitLabInstanceConfig(gitlabConfigFile *GitLabConfigFile, instanceName string) (*GitLabInstanceConfig, error) {
	// get the entry for the passed name
	if gitlabConfigFile != nil && gitlabConfigFile.GitLabInstance != nil {
		for _, instanceConfig := range gitlabConfigFile.GitLabInstance {
			if instanceConfig.Name == instanceName {
				return instanceConfig, nil
			}
		}
	}

	return nil, errors.New("No GitLab Instance Configuration found for " + instanceName)
}

/**
 * Executes the job and waits until completition if specified in the configuration
 */
func executeGitLabPipelineAndWaitForCompletion(eventMapConfig *EventMappingConfig, pipelineConfig *PipelineConfig, gitlabInstanceConfig *GitLabInstanceConfig) (bool, *KeptnResultArtifact, error) {

	// before we execute the job we save current time in the eventMap
	eventMapConfig.startedAt = time.Now()
	eventMapConfig.finishedAt = time.Now()

	git, err := gitlab.NewClient(
		gitlabInstanceConfig.Token,
		gitlab.WithBaseURL("https://"+gitlabInstanceConfig.URL),
	)
	if err != nil {
		log.Fatal(err)
	}

	/* List all projects
	projects, _, err := git.Projects.ListProjects(nil)
	if err != nil {
		log.Fatal(err)
	}*/

	variables := []*gitlab.PipelineVariable{}

	for key, value := range pipelineConfig.Variables {
		pvar := gitlab.PipelineVariable{
			Key:          key,
			Value:        value,
			VariableType: "env_var",
		}

		variables = append(variables, &pvar)
	}

	options := &gitlab.CreatePipelineOptions{
		Ref:       &pipelineConfig.Ref,
		Variables: variables,
	}

	glpipeline, _, err := git.Pipelines.CreatePipeline(pipelineConfig.ProjectID, options)

	if err != nil {
		log.Fatal(err)
		return false, nil, err
	}

	for {

		glstatus, _, err := git.Pipelines.GetPipeline(pipelineConfig.ProjectID, glpipeline.ID)
		if err != nil {
			log.Fatal(err)
			return false, nil, err
		}

		fmt.Println(glstatus.Status)
		time.Sleep(5 * time.Second)
		if glstatus.Status != "pending" && glstatus.Status != "running" {
			if glstatus.Status == "failed" {
				return false, nil, nil
			} else {
				return true, nil, nil
			}
			break
		}
	}

	logMessage := fmt.Sprintf("Pipeline %s did not finish within %d seconds", pipelineConfig.Name, 3600)
	return false, nil, errors.New(logMessage)
}

/**
 * Based on the actionSuccess sends the onSuccess or onFailure Event definition
 */
func sendKeptnEventForEventConfig(incomingBaseEvent *baseKeptnEvent, incomingEvent *cloudevents.Event, eventMappingConfig *EventMappingConfig, actionSuccess bool, keptnResult *KeptnResultArtifact, logger *keptnutils.Logger) (bool, error) {

	// first lets get the correct OnXX data set
	var eventData map[string]string
	if actionSuccess {
		eventData = eventMappingConfig.OnSuccess
	} else {
		eventData = eventMappingConfig.OnFailure
	}

	// lets merge whats in keptnResult with eventData as this is the data that came back from the pipeline
	if keptnResult != nil {
		for name, value := range keptnResult.Data {
			log.Printf("KeptnResult: %s=%s", name, value)
			eventData[name] = value
		}
	}

	// double check if we have to do anything at all, e.g: if there is no OnXX data set we are done
	if len(eventData) == 0 {
		return true, nil
	}

	// we have to have at least the event name
	eventName, eventNameExists := eventData["event"]
	if eventNameExists == false {
		return false, errors.New("No event was specified to send back to Keptn")
	} else {
		log.Println(fmt.Sprintf("Processing OnXX configuration with eventData: %s", eventData))
	}

	// now lets validate that we have at least the event name, e.g: deployment.finished
	switch eventName {
	case "deployment.finished":
		deploymentURILocal, _ := eventData["deploymentURILocal"]
		deploymentURIPublic, _ := eventData["deploymentURIPublic"]
		result, _ := eventData["result"]

		sendDeploymentFinishedEvent(incomingBaseEvent.context, incomingEvent, incomingBaseEvent.project, incomingBaseEvent.service, incomingBaseEvent.stage, incomingBaseEvent.testStrategy, incomingBaseEvent.deployment, "", "",
			deploymentURILocal,
			deploymentURIPublic,
			result,
			incomingBaseEvent.labels, logger)
	case "tests.finished":
		result, _ := eventData["result"]
		sendTestsFinishedEvent(incomingBaseEvent.context, incomingEvent, incomingBaseEvent.project, incomingBaseEvent.service, incomingBaseEvent.stage, incomingBaseEvent.testStrategy, incomingBaseEvent.deployment,
			eventMappingConfig.startedAt,
			eventMappingConfig.finishedAt,
			result,
			incomingBaseEvent.labels, logger)
	default:
		return false, errors.New("Event not supported: " + eventName)
	}

	return true, nil
}
