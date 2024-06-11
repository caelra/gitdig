package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const baseURL = "https://api.github.com"

type Repo struct {
	Owner string
	Name  string
}

type API struct {
	BaseURL string
	Token   string
}

func New(token string) (*API, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}

	return &API{
		BaseURL: baseURL,
		Token:   token,
	}, nil
}

func (a *API) doRequest(req *http.Request, v interface{}) error {
	a.addToken(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if isError(resp) {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github: %s", string(b))
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func isError(resp *http.Response) bool {
	return resp.StatusCode < 200 || resp.StatusCode > 299
}

func (a *API) addToken(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+a.Token)
}

func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func Bool(v bool) *bool {
	return &v
}

func Int(v int) *int {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func String(v string) *string {
	return &v
}
