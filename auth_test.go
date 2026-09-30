package main

import "testing"

func TestGithubToken(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		t.Setenv(tokenEnvVar, "")
		t.Setenv("GITHUB_TOKEN", "")
		if _, err := githubToken(); err == nil {
			t.Fatal("expected error when token env var is unset")
		}
	})

	t.Run("set", func(t *testing.T) {
		t.Setenv(tokenEnvVar, "test-token")
		t.Setenv("GITHUB_TOKEN", "fallback-token")
		got, err := githubToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "test-token" {
			t.Fatalf("got %q, want %q", got, "test-token")
		}
	})

	t.Run("fallback", func(t *testing.T) {
		t.Setenv(tokenEnvVar, "")
		t.Setenv("GITHUB_TOKEN", "fallback-token")
		if got, err := githubToken(); err != nil || got != "fallback-token" {
			t.Fatalf("got %q, %v; want fallback-token", got, err)
		}
	})
}
