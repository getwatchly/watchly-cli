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
	CommitSHA    string `json:"commit_sha"`
	CommitAuthor string `json:"commit_author"`
}

func NotifyDeployment(apiKey, deploymentId, githubSha, githubRunId, githubJobId string) error {
	notification := DeploymentNotification{
		DeploymentID: deploymentId,
		GitHubRunID:  githubRunId,
		GitHubJobID:  githubJobId,
		CommitSHA:    githubSha,
	}

	marshalledNotification, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", WATCHLY_ENDPOINT+"/webhooks/deployments/finish", bytes.NewBuffer(marshalledNotification))
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to notify Watchly: %s", resp.Status)
	}

	return nil
}
