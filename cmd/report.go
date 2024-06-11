package cmd

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/caelra/gitdig/internal/report"
	"github.com/caelra/gitdig/pkg/github"
	"github.com/caelra/gitdig/pkg/mailer"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var cfg Config
var preview bool

var reportCmd = &cobra.Command{
	Use:               "report [recipients]",
	Short:             "create a pull requests report",
	Args:              cobra.MinimumNArgs(1),
	RunE:              runReport,
	PersistentPreRunE: loadEnvironment,
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

func runReport(cmd *cobra.Command, args []string) error {
	if args == nil {
		return fmt.Errorf("empty recipient parameter")
	}

	gh, err := github.New(cfg.GithubToken)
	if err != nil {
		return err
	}

	smtp := mailer.New(cfg.GoogleMail, cfg.GoogleAppPass)

	r := report.New(cfg.Sender, cfg.RepoOwner, cfg.RepoName, gh, smtp)

	if err := r.Generate(args, preview); err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	return nil
}

func init() {
	reportCmd.PersistentFlags().StringVar(&cfg.RepoOwner, "repo_owner", "ziglang", "GitHub repository owner")
	reportCmd.PersistentFlags().StringVar(&cfg.RepoName, "repo_name", "zig", "GitHub repository name")
	reportCmd.Flags().BoolVarP(&preview, "preview", "p", false, "show report")

	rootCmd.AddCommand(reportCmd)
}
