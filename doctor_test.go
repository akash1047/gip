package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunDoctor(t *testing.T) {
	t.Run("no token", func(t *testing.T) {
		t.Setenv(tokenEnvVar, "")
		if err := runDoctor(nil); err == nil {
			t.Fatal("expected error when token is unset")
		}
	})

	t.Run("token set, auth ok", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{}`))
		}))
		defer srv.Close()

		orig := githubAPIBase
		githubAPIBase = srv.URL
		defer func() { githubAPIBase = orig }()

		t.Setenv(tokenEnvVar, "test-token")
		if err := runDoctor(nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("token set, auth fails", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		orig := githubAPIBase
		githubAPIBase = srv.URL
		defer func() { githubAPIBase = orig }()

		t.Setenv(tokenEnvVar, "test-token")
		if err := runDoctor(nil); err == nil {
			t.Fatal("expected error when auth check fails")
		}
	})
}
