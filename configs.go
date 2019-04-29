package main

import (
	"errors"
	"fmt"
	"github.com/bitrise-io/go-utils/log"
	"os"
	"strings"
)

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
		ExportedVariableNames:    os.Getenv("exported_environment_variable_names"),
		BranchRepoOwner:          os.Getenv("branch_repo_owner"),
		BranchDestRepoOwner:      os.Getenv("branch_dest_repo_owner"),
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
	log.Printf(" - ExportedVariableNames: %s", configs.ExportedVariableNames)
	log.Printf(" - BranchRepoOwner: %s", configs.BranchRepoOwner)
	log.Printf(" - BranchDestRepoOwner: %s", configs.BranchDestRepoOwner)
}

func (configs ConfigsModel) validate() error {
	if configs.AppSlug == "" {
		return errors.New("empty App slug specified")
	}

	if configs.APIToken == "" {
		return errors.New("empty Build Trigger API token specified")
	}

	for _, environmentVariableName := range splitPipeSeparatedStringArray(configs.ExportedVariableNames) {
		if environmentVariableName == "" {
			return errors.New("empty environment variable name specified")
		} else if strings.Contains(environmentVariableName, "=") {
			return fmt.Errorf("environment variable with '=' character specified: %s", environmentVariableName)
		}
	}

	return nil
}
