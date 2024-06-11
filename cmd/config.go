package cmd

type Config struct {
	GithubToken   string `env:"GITHUB_TOKEN,required"`
	GoogleMail    string `env:"GOOGLE_MAIL,required"`
	GoogleAppPass string `env:"GOOGLE_APP_PASS,required"`
	Recipient     string `env:"RECIPIENT" envDefault:"caelra99œgmail.com"`
	Sender        string `env:"SENDER" envDefault:"caelra99œgmail.com"`
	RepoOwner     string `env:"REPO_OWNER" envDefault:"ziglang"`
	RepoName      string `env:"REPO_NAME" envDefault:"zig"`
}
