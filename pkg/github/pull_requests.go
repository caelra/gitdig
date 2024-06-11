package github

import (
	"fmt"
	"net/http"
)

type ListOptions struct {
	Page    *int `url:"page,omitempty"`
	PerPage *int `url:"per_page,omitempty"`
}

type PullRequestListOptions struct {
	State     *string `url:"state,omitempty"`
	Head      *string `url:"head,omitempty"`
	Base      *string `url:"base,omitempty"`
	Sort      *string `url:"sort,omitempty"`
	Direction *string `url:"direction,omitempty"`

	*ListOptions
}

type User struct {
	Login     *string `json:"login,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

func (u User) String() string {
	return Stringify(u)
}

type PullRequest struct {
	HTMLURL   *string    `json:"html_url,omitempty"`
	Number    *int       `json:"number,omitempty"`
	State     *string    `json:"state,omitempty"`
	Title     *string    `json:"title,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	ClosedAt  *Timestamp `json:"closed_at,omitempty"`
	MergedAt  *Timestamp `json:"merged_at,omitempty"`
	User      *User      `json:"user,omitempty"`
}

func (p PullRequest) String() string {
	return Stringify(p)
}

func (a *API) ListPullRequests(owner string, repo string, opts *PullRequestListOptions) ([]*PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", a.BaseURL, owner, repo)
	url, err := addOptions(url, opts)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	var pulls []*PullRequest
	err = a.doRequest(req, &pulls)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("PullRequest response: %+v\n", Stringify(pulls))

	return pulls, nil
}
