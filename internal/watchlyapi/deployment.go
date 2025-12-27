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
	Url        string `json:"url"`
	ReleaseTag string `json:"release_tag,omitempty"`
}

type DeploymentStartResponse struct {
	DeploymentID string `json:"id"`
}

func StartDeployment(apiKey, githubSha, deploymentUrl, releaseTag string) (string, error) {
	body := DeploymentStartBody{
		Url:        deploymentUrl,
		ReleaseTag: releaseTag,
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", WATCHLY_ENDPOINT+"/webhooks/deployments/start/"+githubSha, bytes.NewBuffer(marshalledBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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
	Status      string `json:"status"`
	CompletedAt string `json:"completed_at"`
}

func FinishDeployment(apiKey, githubSha, status, completedAt string) error {
	body := DeploymentFinishBody{
		Status:      status,
		CompletedAt: completedAt,
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", WATCHLY_ENDPOINT+"/webhooks/deployments/finish/"+githubSha, bytes.NewBuffer(marshalledBody))
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

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to notify Watchly: %s", resp.Status)
	}

	return nil
}

type ProjectSettingsUpdateBody struct {
	DeploymentFreeze *bool `json:"deployment_freeze,omitempty"`
}

type ProjectSettingsUpdateResponse struct {
	Message string `json:"message"`
}

func UpdateProjectSettings(apiKey string, deploymentFreeze bool) error {
	deploymentFreezePtr := &deploymentFreeze
	body := ProjectSettingsUpdateBody{
		DeploymentFreeze: deploymentFreezePtr,
	}

	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", WATCHLY_ENDPOINT+"/webhooks/projects", bytes.NewBuffer(marshalledBody))
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

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to update project settings: %s", resp.Status)
	}

	var response ProjectSettingsUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	return nil
}

type ProjectEnabledResponse struct {
	Enabled bool `json:"enabled"`
}

func GetDeploymentsEnabled(apiKey string) (bool, error) {
	req, err := http.NewRequest("GET", WATCHLY_ENDPOINT+"/webhooks/deployments/enabled", nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := NewHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return false, fmt.Errorf("failed to get project enabled status: %s", resp.Status)
	}

	var response ProjectEnabledResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, err
	}

	return response.Enabled, nil
}
