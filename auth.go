package main

import (
	"errors"
	"os"
)

// tokenEnvVar is the preferred source of the GitHub PAT.
const tokenEnvVar = "GIP_GITHUB_TOKEN"

// githubToken reads the PAT from the environment. Callers should invoke this
// only when a token is actually needed, so a missing token doesn't fail
// subcommands that don't call the GitHub API.
func githubToken() (string, error) {
	token := os.Getenv(tokenEnvVar)
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return "", errors.New(tokenEnvVar + " or GITHUB_TOKEN environment variable is not set")
	}
	return token, nil
}
