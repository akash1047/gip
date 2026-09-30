package main

import (
	"net/http"
	"testing"
)

func TestRunDoctor(t *testing.T) {
	t.Run("no token", func(t *testing.T) {
		t.Setenv(tokenEnvVar, "")
		t.Setenv("GITHUB_TOKEN", "")
		if err := runDoctor(nil); err == nil {
			t.Fatal("expected error when token is unset")
		}
	})

	t.Run("token set, auth ok", func(t *testing.T) {
		mockGithub(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{}`))
		}))

		t.Setenv(tokenEnvVar, "test-token")
		if err := runDoctor(nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("token set, auth fails", func(t *testing.T) {
		mockGithub(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		t.Setenv(tokenEnvVar, "test-token")
		if err := runDoctor(nil); err == nil {
			t.Fatal("expected error when auth check fails")
		}
	})
}
