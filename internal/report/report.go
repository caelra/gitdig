package report

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caelra/gitdig/pkg/github"
	"github.com/caelra/gitdig/pkg/mailer"
)

type PRReport struct {
	GithubSvc *github.API
	MailerSvc *mailer.SMTP
	Repo      github.Repo
	Sender    string
}

func New(sender string, RepoOwner string, repoName string, githubSvc *github.API, mailerSvc *mailer.SMTP) *PRReport {
	return &PRReport{
		GithubSvc: githubSvc,
		MailerSvc: mailerSvc,
		Repo:      github.Repo{Owner: RepoOwner, Name: repoName},
		Sender:    sender,
	}
}

func (r *PRReport) Generate(recipients []string, verbose bool) error {
	listOpts := &github.ListOptions{
		PerPage: github.Int(50),
	}

	opts := &github.PullRequestListOptions{
		State:       github.String("all"),
		ListOptions: listOpts,
	}

	pulls, err := r.GithubSvc.ListPullRequests(r.Repo.Owner, r.Repo.Name, opts)
	if err != nil {
		fmt.Println("Error fetching pull requests:", err)
		return err
	}

	var (
		buffer  bytes.Buffer
		weekAgo time.Time = time.Now().AddDate(0, 0, -7)

		totalPR int = 0
	)

	for _, pull := range pulls {
		totalPR++

		if err := mailBody(&buffer, pull, weekAgo); err != nil {
			log.Printf("Error writing pull request to buffer: %v", err)
			return err
		}

		if verbose {
			previewBody(pull, weekAgo)
		}
	}

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: Weekly Pull Request Report\r\n\r\n%s\r\n\nTotal PRs: %d",
		strings.Join(recipients, ", "),
		buffer.String(),
		totalPR,
	))

	if err := r.MailerSvc.Send(r.Sender, recipients, msg); err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	return nil
}

func previewBody(pull *github.PullRequest, frecuency time.Time) {
	if pull.GetMergedAt().After(frecuency) && pull.GetState() == "closed" {
		fmt.Printf("[%s] merged #%d %s (%s)\n", pull.GetMergedAt().Format("2006-01-02"), pull.GetNumber(), pull.GetTitle(), pull.User.GetLogin())
	} else if pull.GetCreatedAt().After(frecuency) && pull.GetState() == "open" {
		fmt.Printf("[%s] %s #%d %s (%s)\n", pull.GetCreatedAt().Format("2006-01-02"), pull.GetState(), pull.GetNumber(), pull.GetTitle(), pull.User.GetLogin())
	} else if pull.GetClosedAt().After(frecuency) && pull.GetState() == "closed" {
		fmt.Printf("[%s] %s #%d %s (%s)\n", pull.GetCreatedAt().Format("2006-01-02"), pull.GetState(), pull.GetNumber(), pull.GetTitle(), pull.User.GetLogin())
	}
}

func mailBody(b *bytes.Buffer, pull *github.PullRequest, frecuency time.Time) error {
	var date github.Timestamp
	var state string

	if pull.GetMergedAt().After(frecuency) && pull.GetState() == "closed" {
		date = pull.GetMergedAt()
		state = "merged"
	} else if pull.GetCreatedAt().After(frecuency) && pull.GetState() == "open" {
		date = pull.GetCreatedAt()
		state = pull.GetState()
	} else if pull.GetClosedAt().After(frecuency) && pull.GetState() == "closed" {
		date = pull.GetClosedAt()
		state = pull.GetState()
	} else {
		return nil
	}

	_, err := b.WriteString(fmt.Sprintf(
		"[%s] %s #%d %s (%s)\n\n",
		date.Format("2009-01-02"),
		state,
		pull.GetNumber(),
		pull.GetTitle(),
		pull.User.GetLogin(),
	))
	return err
}
