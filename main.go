package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/bitrise-io/go-utils/log"
)

const (
	triggeredBuildSlug   = "TRIGGERED_BUILD_SLUG"
	triggeredBuildNumber = "TRIGGERED_BUILD_NUMBER"
	triggeredBuildURL    = "TRIGGERED_BUILD_URL"
	triggeredWorkflowID  = "TRIGGERED_WORKFLOW_ID"
)

func main() {
	configs := createConfigsModelFromEnvs()
	configs.dump()
	if err := configs.validate(); err != nil {
		log.Errorf("Issue with input: %s", err)
		os.Exit(1)
	}

	requestBody, err := createRequestBodyFromConfigs(configs)
	if err != nil {
		log.Errorf("Could not create request body, error: %s", err)
		os.Exit(2)
	}

	request, err := createRequest(configs.AppSlug, requestBody)
	if err != nil {
		log.Errorf("Could not create request, error: %s", err)
		os.Exit(2)
	}

	responseModel, err := performRequest(request)
	if err != nil {
		log.Errorf("Could not send request, error: %s", err)
		os.Exit(3)
	}

	log.Infof("Build Trigger status: %s", responseModel.Status)

	if responseModel.Message != "ok" {
		log.Errorf("Build not triggered, status: %s", responseModel.Message)
		os.Exit(4)
	}

	fmt.Println()
	log.Infof("Triggered build slug: %s", responseModel.BuildSlug)
	log.Infof("Triggered build number: %d", responseModel.BuildNumber)
	log.Infof("Triggered build URL: %s", responseModel.BuildURL)
	log.Infof("Triggered workflow ID: %s", responseModel.TriggeredWorkflow)

	if err := exportEnvironmentWithEnvman(triggeredBuildSlug, responseModel.BuildSlug); err != nil {
		log.Errorf("Could not export triggered build slug: %s", err)
		os.Exit(5)
	}

	if err := exportEnvironmentWithEnvman(triggeredBuildNumber, strconv.Itoa(responseModel.BuildNumber)); err != nil {
		log.Errorf("Could not export triggered build number: %s", err)
		os.Exit(5)
	}

	if err := exportEnvironmentWithEnvman(triggeredBuildURL, responseModel.BuildURL); err != nil {
		log.Errorf("Could not export triggered build URL: %s", err)
		os.Exit(5)
	}

	if err := exportEnvironmentWithEnvman(triggeredWorkflowID, responseModel.TriggeredWorkflow); err != nil {
		log.Errorf("Could not export triggered workflow: %s", err)
		os.Exit(5)
	}
}

func createRequestBodyFromConfigs(configs ConfigsModel) ([]byte, error) {
	requestModel := RequestModel{
		HookInfo: HookInfoModel{
			Type:     "bitrise",
			APIToken: configs.APIToken,
		},
		BuildParams: BuildParamsModel{
			Branch:                   configs.Branch,
			Tag:                      configs.Tag,
			CommitHash:               configs.CommitHash,
			CommitMessage:            configs.CommitMessage,
			WorkflowID:               configs.WorkflowID,
			BranchDest:               configs.BranchDest,
			PullRequestID:            configs.PullRequestID,
			PullRequestRepositoryURL: configs.PullRequestRepositoryURL,
			PullRequestMergeBranch:   configs.PullRequestMergeBranch,
			PullRequestHeadBranch:    configs.PullRequestHeadBranch,
			Environments:             createExportedEnvironment(configs.ExportedVariableNames),
			BranchRepoOwner:          configs.BranchRepoOwner,
			BranchDestRepoOwner:      configs.BranchDestRepoOwner,
		},
	}

	return json.Marshal(requestModel)
}

func createRequest(appSlug string, body []byte) (*http.Request, error) {
	requestURL := fmt.Sprintf("https://app.bitrise.io/app/%s/build/start.json", appSlug)
	request, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(body))
	if request != nil {
		request.Header.Add("Content-Type", "application/json")
	}
	return request, err
}

func performRequest(request *http.Request) (ResponseModel, error) {
	client := http.Client{}
	response, err := client.Do(request)
	var responseModel ResponseModel

	if err != nil {
		return responseModel, err
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Warnf("Failed to close response body, error: %s", err)
		}
	}()

	log.Infof("Build Trigger API HTTP response status: %s", response.Status)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return responseModel, err
	}

	err = json.Unmarshal(contents, &responseModel)
	return responseModel, err
}
