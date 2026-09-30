package main

import "testing"

func TestGithubToken(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		t.Setenv(tokenEnvVar, "")
		if _, err := githubToken(); err == nil {
			t.Fatal("expected error when token env var is unset")
		}
	})

	t.Run("set", func(t *testing.T) {
		t.Setenv(tokenEnvVar, "test-token")
		got, err := githubToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "test-token" {
			t.Fatalf("got %q, want %q", got, "test-token")
		}
	})
}
