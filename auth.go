package main

import (
	"errors"
	"os"
)

// tokenEnvVar is the only source of the GitHub PAT; no flag/config fallback (issue #6).
const tokenEnvVar = "GIP_GITHUB_TOKEN"

// githubToken reads the PAT from the environment. Callers should invoke this
// only when a token is actually needed, so a missing token doesn't fail
// subcommands that don't call the GitHub API.
func githubToken() (string, error) {
	token := os.Getenv(tokenEnvVar)
	if token == "" {
		return "", errors.New(tokenEnvVar + " environment variable is not set")
	}
	return token, nil
}
