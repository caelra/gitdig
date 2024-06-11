package github

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	api, err := New("test_token")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if api.Token != "test_token" {
		t.Errorf("Expected token to be 'test_token', got %s", api.Token)
	}

	_, err = New("")
	if err == nil {
		t.Fatal("Expected error for empty token, got none")
	}
}

func TestAddToken(t *testing.T) {
	api, _ := New("test_token")
	req, _ := http.NewRequest("GET", "https://api.github.com", nil)
	api.addToken(req)

	if got := req.Header.Get("Authorization"); got != "Bearer test_token" {
		t.Errorf("Expected 'Bearer test_token', got %s", got)
	}
}

func TestIsError(t *testing.T) {
	tests := []struct {
		statusCode int
		isError    bool
	}{
		{statusCode: 200, isError: false},
		{statusCode: 299, isError: false},
		{statusCode: 300, isError: true},
		{statusCode: 400, isError: true},
	}

	for _, test := range tests {
		resp := &http.Response{StatusCode: test.statusCode}
		if got := isError(resp); got != test.isError {
			t.Errorf("isError(%d) = %v; want %v", test.statusCode, got, test.isError)
		}
	}
}

func TestAddOptions(t *testing.T) {
	type Options struct {
		Page int `url:"page,omitempty"`
	}

	s, err := addOptions("https://api.github.com/repos", Options{Page: 2})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "https://api.github.com/repos?page=2"
	if s != expected {
		t.Errorf("Expected %s, got %s", expected, s)
	}
}

func TestDoRequest(t *testing.T) {
	api, _ := New("test_token")

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	})

	req, _ := http.NewRequest("GET", server.URL, nil)

	var result map[string]interface{}
	err := api.doRequest(req, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["success"] != true {
		t.Errorf("Expected success to be true, got %v", result["success"])
	}
}
