package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getwatchly/watchly-cli/internal/watchlyapi"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "watchly-cli",
		Usage: "CLI to interact with Watchly - Docs at https://docs.watchly.dev",
		Commands: []*cli.Command{
			{
				Name:  "deployment",
				Usage: "Notify Watchly about a deployment",
				Commands: []*cli.Command{
					{
						Name:  "start",
						Usage: "Start a deployment",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "api-key",
								Aliases:  []string{"k"},
								Usage:    "Watchly API key for your project",
								Sources:  cli.EnvVars("WATCHLY_API_KEY"),
								Required: true,
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							apiKey := cmd.String("api-key")
							githubRepoFullName := os.Getenv("GITHUB_REPOSITORY")
							githubSha := os.Getenv("GITHUB_SHA")
							githubRunId := os.Getenv("GITHUB_RUN_ID")
							githubJobId := os.Getenv("GITHUB_JOB")

							if githubRepoFullName == "" || githubSha == "" || githubRunId == "" || githubJobId == "" {
								return fmt.Errorf("missing required environment variables, are you running this in a GitHub Actions environment?")
							}

							fmt.Println("watchly-cli - 🔭 Contacting Watchly ...")

							deploymentId, err := watchlyapi.StartDeployment(apiKey, githubRepoFullName, githubSha, githubRunId, githubJobId)
							if err != nil {
								return fmt.Errorf("failed to notify Watchly: %w", err)
							}

							fmt.Printf("watchly-cli - ✅ Recorded deployment: %s\n", deploymentId)

							githubEnvFile := os.Getenv("GITHUB_ENV")
							if githubEnvFile == "" {
								fmt.Println("GITHUB_ENV not set - skipping environment file update. WATCHLY_DEPLOYMENT_ID will not be available in subsequent steps.")
								return nil
							}

							f, err := os.OpenFile(githubEnvFile, os.O_APPEND|os.O_WRONLY, 0600)
							if err != nil {
								return fmt.Errorf("failed to open GITHUB_ENV file: %w", err)
							}
							defer f.Close()

							if _, err := f.WriteString(fmt.Sprintf("WATCHLY_DEPLOYMENT_ID=%s\n", deploymentId)); err != nil {
								return fmt.Errorf("failed to write to GITHUB_ENV file: %w", err)
							}

							fmt.Println("watchly-cli - ✅ Set WATCHLY_DEPLOYMENT_ID environment variable for subsequent steps")

							return nil
						},
					},
					{
						Name:  "finish",
						Usage: "Finish a deployment",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "api-key",
								Aliases:  []string{"k"},
								Usage:    "Watchly API key for your project",
								Sources:  cli.EnvVars("WATCHLY_API_KEY"),
								Required: true,
							},
							&cli.StringFlag{
								Name:    "deployment-id",
								Aliases: []string{"d"},
								Value:   "",
								Usage:   "ID of the deployment to notify Watchly about",
								Sources: cli.EnvVars("WATCHLY_DEPLOYMENT_ID"),
							},
							&cli.StringFlag{
								Name:     "status",
								Aliases:  []string{"s"},
								Required: true,
								Usage:    "Status of the deployment, one of 'successful' or 'failed'",
								Sources:  cli.EnvVars("WATCHLY_STATUS"),
							},
							&cli.StringFlag{
								Name:    "completed-at",
								Aliases: []string{"c"},
								Value:   "",
								Usage:   "Completion time of the deployment",
								Sources: cli.EnvVars("WATCHLY_COMPLETED_AT"),
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							apiKey := cmd.String("api-key")
							deploymentId := cmd.String("deployment-id")
							status := cmd.String("status")
							completedAt := cmd.String("completed-at")

							if deploymentId == "" {
								return fmt.Errorf("missing deployment-id, this should not have happened. check if env variable WATCHLY_DEPLOYMENT_ID is set")
							}

							if status != "successful" && status != "failed" {
								return fmt.Errorf("invalid status: %s", status)
							}

							if completedAt == "" {
								completedAt = time.Now().UTC().Format(time.RFC3339)
							}

							if err := watchlyapi.FinishDeployment(apiKey, deploymentId, status, completedAt); err != nil {
								return fmt.Errorf("failed to notify Watchly: %w", err)
							}

							fmt.Println("watchly-cli - ✅ Finished deployment")

							return nil
						},
					},
				},
			}},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
