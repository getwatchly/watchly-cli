package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/getwatchly/watchly-cli/internal/watchlyapi"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "deployment",
		Usage: "Notify Watchly about a deployment",
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			apiKey := cmd.String("api-key")
			deploymentId := cmd.String("deployment-id")
			githubSha := os.Getenv("GITHUB_SHA")
			githubRunId := os.Getenv("GITHUB_RUN_ID")
			githubJobId := os.Getenv("GITHUB_JOB")

			if githubSha == "" || githubRunId == "" || githubJobId == "" {
				return fmt.Errorf("missing required environment variables, are you running this in a GitHub Actions environment?")
			}

			fmt.Println("watchly-cli - 🔭 Contacting Watchly ...")

			if err := watchlyapi.NotifyDeployment(apiKey, deploymentId, githubSha, githubRunId, githubJobId); err != nil {
				return fmt.Errorf("failed to notify Watchly: %w", err)
			}

			fmt.Println("watchly-cli - ✅ Recorded deployment")

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
