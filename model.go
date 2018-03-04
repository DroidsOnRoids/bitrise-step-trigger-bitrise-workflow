package main


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
	ExportedVariableNames    string
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
	Environments             []EnvironmentVariableModel `json:"environments"`
}

// EnvironmentVariableModel ...
type EnvironmentVariableModel struct {
	MappedTo string `json:"mapped_to"`
	Value    string `json:"value"`
	IsExpand bool `json:"is_expand"`
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