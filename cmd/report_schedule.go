package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/caelra/gitdig/internal/report"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var reportScheduleCmd = &cobra.Command{
	Use:               "report-schedule [recipients]",
	Short:             "Schedule pull requests reports",
	Args:              cobra.MinimumNArgs(1),
	RunE:              runCronReport,
	PersistentPreRunE: loadEnvironment,
}

func runCronReport(cmd *cobra.Command, args []string) error {
	crontab, err := constructCrontab(cfg.ScheduleDay, cfg.ScheduleTime)
	if err != nil {
		return err
	}

	log.Printf("Setting up cron job with schedule: %s", crontab)
	c := cron.New()

	_, err = c.AddFunc(crontab, func() {
		log.Printf("Generating report...\n")
		gh, smtp, err := initializeServices()
		if err != nil {
			log.Printf("Failed to connect to services: %s", err)
			return
		}

		r := report.New(cfg.Sender, cfg.RepoOwner, cfg.RepoName, gh, smtp)

		if err := r.Generate(args, preview); err != nil {
			log.Printf("Failed to generate report: %s", err)
			return
		}

		log.Printf("Report sent...")
	})
	if err != nil {
		log.Fatal("Error adding cron job:", err)
	}

	c.Start()
	select {}
}

func constructCrontab(day, time string) (string, error) {
	days := map[string]string{
		"sunday":    "0",
		"monday":    "1",
		"tuesday":   "2",
		"wednesday": "3",
		"thursday":  "4",
		"friday":    "5",
		"saturday":  "6",
	}

	dayOfWeek, ok := days[strings.ToLower(day)]
	if !ok {
		return "", fmt.Errorf("invalid day: %s", day)
	}

	timeParts := strings.Split(time, ":")
	if len(timeParts) != 2 {
		return "", fmt.Errorf("invalid time format: %s", time)
	}
	hour := timeParts[0]
	minute := timeParts[1]

	return fmt.Sprintf("%s %s * * %s", minute, hour, dayOfWeek), nil
}

func init() {
	reportScheduleCmd.PersistentFlags().StringVar(&cfg.RepoOwner, "repo_owner", "ziglang", "GitHub repository owner")
	reportScheduleCmd.PersistentFlags().StringVar(&cfg.RepoName, "repo_name", "zig", "GitHub repository name")
	reportScheduleCmd.PersistentFlags().StringVar(&cfg.ScheduleDay, "day", "monday", "Day of the week to run the report (e.g., monday)")
	reportScheduleCmd.PersistentFlags().StringVar(&cfg.ScheduleTime, "time", "09:00", "Time to run the report in HH:MM format")
	reportScheduleCmd.Flags().BoolVarP(&preview, "preview", "p", false, "Show report")
	rootCmd.AddCommand(reportScheduleCmd)
}
