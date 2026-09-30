package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchIssue(t *testing.T) {
	t.Setenv(tokenEnvVar, "test-token")

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/issues/42", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization header = %q", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"number": 42, "title": "t", "body": "b", "state": "open",
			"html_url": "https://example/42",
			"user":     map[string]string{"login": "alice"},
		})
	})
	mux.HandleFunc("/repos/o/r/issues/42/comments", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{
			{"body": "hi", "user": map[string]string{"login": "bob"}},
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	orig := githubAPIBase
	githubAPIBase = srv.URL
	defer func() { githubAPIBase = orig }()

	issue, err := FetchIssue("o", "r", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.Title != "t" || issue.Author != "alice" {
		t.Fatalf("got %+v", issue)
	}
	if len(issue.Comments) != 1 || issue.Comments[0].Author != "bob" {
		t.Fatalf("got comments %+v", issue.Comments)
	}
}
