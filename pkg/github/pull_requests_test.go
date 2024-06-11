package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAPI struct {
	BaseURL string
}

func (m *MockAPI) doRequest(_ *http.Request, v interface{}) error {
	resp := httptest.NewRecorder()

	// Mock response for the test
	respBody := `[{
		"html_url": "https://github.com/owner/repo/pull/1",
		"number": 1,
		"state": "open",
		"title": "Add new feature",
		"created_at": "2023-05-01T12:34:56Z",
		"user": {
			"login": "octocat",
			"avatar_url": "https://github.com/images/error/octocat_happy.gif"
		}
	}]`
	resp.WriteString(respBody)
	resp.Code = http.StatusOK

	return json.NewDecoder(resp.Body).Decode(v)
}

func TestListPullRequests(t *testing.T) {
	api := &MockAPI{
		BaseURL: "https://api.github.com",
	}

	tests := []struct {
		owner       string
		repo        string
		opts        *PullRequestListOptions
		expectedLen int
		expectErr   bool
	}{
		{
			owner:       "owner",
			repo:        "repo",
			opts:        &PullRequestListOptions{},
			expectedLen: 1,
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%s", tt.owner, tt.repo), func(t *testing.T) {
			prs, err := api.ListPullRequestsMock(tt.owner, tt.repo, tt.opts)
			if (err != nil) != tt.expectErr {
				t.Fatalf("ListPullRequests() error = %v, expectErr %v", err, tt.expectErr)
			}
			if len(prs) != tt.expectedLen {
				t.Errorf("ListPullRequests() got = %d, want %d", len(prs), tt.expectedLen)
			}
		})
	}
}

func (a *MockAPI) ListPullRequestsMock(owner string, repo string, opts *PullRequestListOptions) ([]*PullRequest, error) {
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
	fmt.Printf("PullRequest response: %+v\n", Stringify(pulls))

	return pulls, nil
}

func addOptionsMock(s string, _ interface{}) (string, error) {
	return s, nil
}

func StringifyMock(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
