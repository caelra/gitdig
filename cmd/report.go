package cmd

import (
	"fmt"

	"github.com/caelra/gitdig/internal/report"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:               "report [recipients]",
	Short:             "Create a pull requests report",
	Args:              cobra.MinimumNArgs(1),
	RunE:              runReport,
	PersistentPreRunE: loadEnvironment,
}

func runReport(cmd *cobra.Command, args []string) error {
	gh, smtp, err := initializeServices()
	if err != nil {
		return err
	}

	r := report.New(cfg.Sender, cfg.RepoOwner, cfg.RepoName, gh, smtp)

	if err := r.Generate(args, preview); err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	return nil
}

func init() {
	reportCmd.PersistentFlags().StringVar(&cfg.RepoOwner, "repo_owner", "ziglang", "GitHub repository owner")
	reportCmd.PersistentFlags().StringVar(&cfg.RepoName, "repo_name", "zig", "GitHub repository name")
	reportCmd.Flags().BoolVarP(&preview, "preview", "p", false, "show report when it generates")
	rootCmd.AddCommand(reportCmd)
}
