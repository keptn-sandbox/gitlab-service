package main

import (
	"fmt"
	"log"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

/**
* Here are all the handler functions for the individual event
  See https://github.com/keptn/spec/blob/0.1.3/cloudevents.md for details on the payload

  -> "sh.keptn.event.configuration.change"
  -> "sh.keptn.events.deployment-finished"
  -> "sh.keptn.events.tests-finished"
  -> "sh.keptn.event.start-evaluation"
  -> "sh.keptn.events.evaluation-done"
  -> "sh.keptn.event.problem.open"
  -> "sh.keptn.events.problem"
*/

func handleEvent(eventname string, event cloudevents.Event, keptnEvent baseKeptnEvent, logger *keptn.Logger) error {
	gitlabConfigFile, err := getGitLabConfiguration(keptnEvent, logger)
	if err != nil {
		return err
	}

	logger.Info("Loaded GitLab Config File")
	eventMapConfig, err := getEventMappingConfig(gitlabConfigFile, eventname)
	if err != nil {
		logger.Info(fmt.Sprintf("No event mapping found for %s. Therefore executing no GitLab Pipeline", eventname))
		return nil
	}

	logger.Info(fmt.Sprintf("Event Mapping %s points to action %s", eventname, eventMapConfig.Action))

	var pipelineConfig *PipelineConfig
	pipelineConfig, err = getPipelineConfig(gitlabConfigFile, eventMapConfig.Action)
	if err != nil {
		return fmt.Errorf("No Action Config found for: %s", eventMapConfig.Action)
	}

	var instanceConfig *GitLabInstanceConfig
	instanceConfig, err = getGitLabInstanceConfig(gitlabConfigFile, pipelineConfig.Instance)
	if err != nil {
		return fmt.Errorf("No Server Config found for: %s", pipelineConfig.Instance)
	}

	var success bool
	var keptnResult *KeptnResultArtifact
	success, keptnResult, err = executeGitLabPipelineAndWaitForCompletion(eventMapConfig, pipelineConfig, instanceConfig)

	success, err = sendKeptnEventForEventConfig(
		&keptnEvent, &event,
		eventMapConfig,
		success, keptnResult,
		logger)

	return err
}

//
// Handles ConfigurationChangeEventType = "sh.keptn.event.configuration.change"
// TODO: add in your handler code
//
func HandleConfigurationChangeEvent(myKeptn *keptn.Keptn, keptnEvent baseKeptnEvent, incomingEvent cloudevents.Event, data *keptn.ConfigurationChangeEventData, logger *keptn.Logger) error {
	log.Printf("Handling Configuration Changed Event: %s", incomingEvent.Context.GetID())

	return handleEvent("configuration.change", incomingEvent, keptnEvent, logger)
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleDeploymentFinishedEvent(myKeptn *keptn.Keptn, keptnEvent baseKeptnEvent, incomingEvent cloudevents.Event, data *keptn.DeploymentFinishedEventData, logger *keptn.Logger) error {
	log.Printf("Handling Deployment Finished Event: %s", incomingEvent.Context.GetID())
	return handleEvent("deployment.finished", incomingEvent, keptnEvent, logger)
	// capture start time for tests
	// startTime := time.Now()

	// run tests
	// ToDo: Implement your tests here

	// Send Test Finished Event
	// return myKeptn.SendTestsFinishedEvent(&incomingEvent, "", "", startTime, "pass", nil, "gitlab-service")
}

//
// Handles TestsFinishedEventType = "sh.keptn.events.tests-finished"
// TODO: add in your handler code
//
func HandleTestsFinishedEvent(myKeptn *keptn.Keptn, keptnEvent baseKeptnEvent, incomingEvent cloudevents.Event, data *keptn.TestsFinishedEventData, logger *keptn.Logger) error {
	log.Printf("Handling Tests Finished Event: %s", incomingEvent.Context.GetID())
	return handleEvent("test.finished", incomingEvent, keptnEvent, logger)

}

//
// Handles EvaluationDoneEventType = "sh.keptn.events.evaluation-done"
// TODO: add in your handler code
//
func HandleStartEvaluationEvent(myKeptn *keptn.Keptn, keptnEvent baseKeptnEvent, incomingEvent cloudevents.Event, data *keptn.StartEvaluationEventData, logger *keptn.Logger) error {
	log.Printf("Handling Start Evaluation Event: %s", incomingEvent.Context.GetID())
	return handleEvent("start.evaluation", incomingEvent, keptnEvent, logger)
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleEvaluationDoneEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.EvaluationDoneEventData) error {
	log.Printf("Handling Evaluation Done Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles ProblemOpenEventType = "sh.keptn.event.problem.open"
// Handles ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
//
func HandleProblemEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ProblemEventData) error {
	log.Printf("Handling Problem Event: %s", incomingEvent.Context.GetID())

	return nil
}
