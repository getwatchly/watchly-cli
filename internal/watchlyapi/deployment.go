package watchlyapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const WATCHLY_ENDPOINT = "https://api.watchly.dev/api/v1"

type DeploymentNotification struct {
	DeploymentID string `json:"id"`
	GitHubRunID  string `json:"github_run_id"`
	GitHubJobID  string `json:"github_job_id"`
	CommitSHA    string `json:"commit_hash"`
	CommitAuthor string `json:"commit_author"`
}

type DeploymentStartBody struct {
	CommitSHA string `json:"commit_hash"`
	Url       string `json:"url"`
}

type DeploymentStartResponse struct {
	DeploymentID string `json:"id"`
}

func StartDeployment(apiKey, githubRepoFullName, githubSha, githubRunId string) (string, error) {
	body := DeploymentStartBody{
		CommitSHA: githubSha,
		Url:       fmt.Sprintf("https://github.com/%s/actions/runs/%s", githubRepoFullName, githubRunId),
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", WATCHLY_ENDPOINT+"/webhooks/deployments/start", bytes.NewBuffer(marshalledBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	
	q := req.URL.Query()
	q.Add("api_key", apiKey)
	req.URL.RawQuery = q.Encode()

	client := NewHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to notify Watchly: %s", resp.Status)
	}

	var response DeploymentStartResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.DeploymentID, nil
}

type DeploymentFinishBody struct {
	DeploymentID string `json:"id"`
	Status       string `json:"status"`
	CompletedAt  string `json:"completed_at"`
}

func FinishDeployment(apiKey, deploymentId, status, completedAt string) error {
	body := DeploymentFinishBody{
		DeploymentID: deploymentId,
		Status:       status,
		CompletedAt:  completedAt,
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", WATCHLY_ENDPOINT+"/webhooks/deployments/finish", bytes.NewBuffer(marshalledBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := NewHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to notify Watchly: %s", resp.Status)
	}

	return nil
}
