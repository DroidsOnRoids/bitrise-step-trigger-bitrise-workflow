package main

import (
	"os"
	"errors"
	"github.com/bitrise-io/go-utils/log"
	"fmt"
	"strings"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/bitrise-io/go-utils/command"
)

const (
	triggeredBuildSlug = "TRIGGERED_BUILD_SLUG"
	triggeredBuildNumber = "TRIGGERED_BUILD_NUMBER"
	triggeredBuildURL = "TRIGGERED_BUILD_URL"
	triggeredWorkflow = "TRIGGERED_WORKFLOW"
)

// ConfigsModel ...
type ConfigsModel struct {
	AppSlug                  string
	APIToken                 string
	Branch                   string
	Tag                      string
	CommitHash               string
	CommitMessage            string
	WorkflowID               string
	BranchDest               string
	PullRequestID            string
	PullRequestRepositoryURL string
	PullRequestMergeBranch   string
	PullRequestHeadBranch    string
}

// RequestModel ...
type RequestModel struct {
	HookInfo    HookInfoModel `json:"hook_info"`
	BuildParams BuildParamsModel `json:"build_params"`
}

// HookInfoModel ...
type HookInfoModel struct {
	Type     string `json:"type"`
	APIToken string `json:"api_token"`
}

// BuildParamsModel ...
type BuildParamsModel struct {
	Branch                   string `json:"branch"`
	Tag                      string `json:"tag"`
	CommitHash               string `json:"commit_hash"`
	CommitMessage            string `json:"commit_message"`
	WorkflowID               string `json:"workflow_id"`
	BranchDest               string `json:"branch_dest"`
	PullRequestID            string `json:"pull_request_id"`
	PullRequestRepositoryURL string `json:"pull_request_repository_url"`
	PullRequestMergeBranch   string `json:"pull_request_merge_branch"`
	PullRequestHeadBranch    string `json:"pull_request_head_branch"`
}

// ResponseModel ...
type ResponseModel struct {
	Status            string `json:"message"`
	Message           string `json:"status"`
	BuildSlug         string `json:"build_slug"`
	BuildNumber       int `json:"build_number"`
	BuildURL          string `json:"build_url"`
	TriggeredWorkflow string `json:"triggered_workflow"`
}

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

	if err := exportEnvironmentWithEnvman(triggeredBuildNumber, string(responseModel.BuildNumber)); err != nil {
		log.Errorf("Could not export triggered build number: %s", err)
		os.Exit(5)
	}

	if err := exportEnvironmentWithEnvman(triggeredBuildURL, responseModel.BuildURL); err != nil {
		log.Errorf("Could not export triggered build URL: %s", err)
		os.Exit(5)
	}

	if err := exportEnvironmentWithEnvman(triggeredWorkflow, responseModel.TriggeredWorkflow); err != nil {
		log.Errorf("Could not export triggered workflow: %s", err)
		os.Exit(5)
	}
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		AppSlug:                  os.Getenv("app_slug"),
		APIToken:                 os.Getenv("api_token"),
		Branch:                   os.Getenv("branch"),
		Tag:                      os.Getenv("tag"),
		CommitHash:               os.Getenv("commit_hash"),
		CommitMessage:            os.Getenv("commit_message"),
		WorkflowID:               os.Getenv("workflow_id"),
		BranchDest:               os.Getenv("branch_dest"),
		PullRequestID:            os.Getenv("pull_request_id"),
		PullRequestRepositoryURL: os.Getenv("pull_request_repository_url"),
		PullRequestMergeBranch:   os.Getenv("pull_request_merge_branch"),
		PullRequestHeadBranch:    os.Getenv("pull_request_head_branch"),
	}
}

func (configs ConfigsModel) dump() {
	fmt.Println()
	log.Infof("Configs:")
	log.Printf(" - AppSlug (hidden): %s", configs.AppSlug)
	log.Printf(" - ApiToken (hidden): %s", strings.Repeat("*", 5))
	log.Printf(" - Branch: %s", configs.Branch)
	log.Printf(" - Tag: %s", configs.Tag)
	log.Printf(" - CommitHash: %s", configs.CommitHash)
	log.Printf(" - CommitMessage: %s", configs.CommitMessage)
	log.Printf(" - WorkflowID: %s", configs.WorkflowID)
	log.Printf(" - BranchDest: %s", configs.BranchDest)
	log.Printf(" - PullRequestID: %s", configs.PullRequestID)
	log.Printf(" - PullRequestRepositoryURL: %s", configs.PullRequestRepositoryURL)
	log.Printf(" - PullRequestMergeBranch: %s", configs.PullRequestMergeBranch)
	log.Printf(" - PullRequestHeadBranch: %s", configs.PullRequestHeadBranch)
}

func (configs ConfigsModel) validate() error {
	if configs.AppSlug == "" {
		return errors.New("empty App slug specified")
	}
	if configs.APIToken == "" {
		return errors.New("empty Build Trigger API token specified")
	}
	return nil
}

func createRequestBodyFromConfigs(configs ConfigsModel) ([]byte, error) {
	requestModel := RequestModel{
		HookInfo:HookInfoModel{
			Type:"bitrise",
			APIToken:configs.APIToken,
		},
		BuildParams:BuildParamsModel{
			Branch:configs.Branch,
			Tag:configs.Tag,
			CommitHash:configs.CommitHash,
			CommitMessage:configs.CommitMessage,
			WorkflowID:configs.WorkflowID,
			BranchDest:configs.BranchDest,
			PullRequestID:configs.PullRequestID,
			PullRequestRepositoryURL:configs.PullRequestRepositoryURL,
			PullRequestHeadBranch:configs.PullRequestHeadBranch,
		},
	}

	return json.Marshal(requestModel)
}

func createRequest(appSlug string, body []byte) (*http.Request, error) {
	requestURL := fmt.Sprintf("https://www.bitrise.io/app/%s/build/start.json", appSlug)
	request, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(body))
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

func exportEnvironmentWithEnvman(keyStr, valueStr string) error {
	cmd := command.New("envman", "add", "--key", keyStr)
	cmd.SetStdin(strings.NewReader(valueStr))
	return cmd.Run()
}