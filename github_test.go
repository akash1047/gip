package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func mockGithub(t *testing.T, handler http.Handler) {
	t.Helper()
	origBase, origClient := githubAPIBase, http.DefaultClient
	githubAPIBase = "http://github.test"
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w.Result(), nil
	})}
	t.Cleanup(func() { githubAPIBase, http.DefaultClient = origBase, origClient })
}

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

	mockGithub(t, mux)

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

func TestFetchIssueWithoutToken(t *testing.T) {
	t.Setenv(tokenEnvVar, "")

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/issues/42", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Authorization"]; ok {
			t.Error("unexpected Authorization header on issue request")
		}
		json.NewEncoder(w).Encode(map[string]any{"number": 42, "title": "public issue"})
	})
	mux.HandleFunc("/repos/o/r/issues/42/comments", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Authorization"]; ok {
			t.Error("unexpected Authorization header on comments request")
		}
		json.NewEncoder(w).Encode([]any{})
	})

	mockGithub(t, mux)

	issue, err := FetchIssue("o", "r", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.Title != "public issue" {
		t.Fatalf("got issue %+v", issue)
	}
}

func TestFetchIssueInvalidToken(t *testing.T) {
	t.Setenv(tokenEnvVar, "bad-token")
	mockGithub(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer bad-token" {
			t.Errorf("Authorization header = %q", got)
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))

	_, err := FetchIssue("o", "r", 42)
	if err == nil || !strings.Contains(err.Error(), "401 Unauthorized") {
		t.Fatalf("expected API authorization error, got %v", err)
	}
}
