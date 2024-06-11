package cmd

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/caelra/gitdig/pkg/github"
	"github.com/caelra/gitdig/pkg/mailer"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	cfg     Config
	preview bool
)

var rootCmd = &cobra.Command{
	Use:   "gitdig",
	Short: "creates github reports",
}

func Execute() error {
	return rootCmd.Execute()
}

func loadEnvironment(cmd *cobra.Command, args []string) error {
	log.Print("Starting Pre-run")
	if err := godotenv.Load(); err != nil {
		return err
	}

	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("unable to load environment: %w", err)
	}

	return nil
}

func initializeServices() (*github.API, *mailer.SMTP, error) {
	gh, err := github.New(cfg.GithubToken)
	if err != nil {
		return nil, nil, err
	}

	smtp := mailer.New(cfg.GoogleMail, cfg.GoogleAppPass)
	return gh, smtp, nil
}
