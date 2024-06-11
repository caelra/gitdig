package github

func (p *PullRequest) GetHTMLURL() string {
	if p == nil || p.HTMLURL == nil {
		return ""
	}
	return *p.HTMLURL
}

func (p *PullRequest) GetNumber() int {
	if p == nil || p.Number == nil {
		return 0
	}
	return *p.Number
}

func (p *PullRequest) GetState() string {
	if p == nil || p.State == nil {
		return ""
	}
	return *p.State
}

func (p *PullRequest) GetTitle() string {
	if p == nil || p.Title == nil {
		return ""
	}
	return *p.Title
}

func (p *PullRequest) GetCreatedAt() Timestamp {
	if p == nil || p.CreatedAt == nil {
		return Timestamp{}
	}
	return *p.CreatedAt
}

func (p *PullRequest) GetClosedAt() Timestamp {
	if p == nil || p.ClosedAt == nil {
		return Timestamp{}
	}
	return *p.ClosedAt
}

func (p *PullRequest) GetMergedAt() Timestamp {
	if p == nil || p.MergedAt == nil {
		return Timestamp{}
	}
	return *p.MergedAt
}

func (p *PullRequest) GetUser() *User {
	if p == nil {
		return nil
	}
	return p.User
}

func (u *User) GetLogin() string {
	if u == nil || u.Login == nil {
		return ""
	}
	return *u.Login
}

func (u *User) GetAvatarURL() string {
	if u == nil || u.AvatarURL == nil {
		return ""
	}
	return *u.AvatarURL
}
