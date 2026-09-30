package main

import "fmt"

// runDoctor implements the `doctor` subcommand: verifies a GitHub token
// (GIP_GITHUB_TOKEN or GITHUB_TOKEN) is set and that it actually
// authenticates against the GitHub API.
func runDoctor(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: gip doctor")
	}

	token, tokenErr := githubToken()
	printCheck("GitHub token is set", tokenErr)

	var authErr error
	if tokenErr == nil {
		authErr = checkGithubAuth(token)
		printCheck("GitHub API authentication", authErr)
	}

	if tokenErr != nil || authErr != nil {
		return fmt.Errorf("doctor found problems")
	}
	return nil
}

func printCheck(name string, err error) {
	if err != nil {
		fmt.Printf("[FAIL] %s: %v\n", name, err)
		return
	}
	fmt.Printf("[PASS] %s\n", name)
}

// checkGithubAuth makes one lightweight authenticated call to confirm the
// token works, without requiring any particular token scope.
func checkGithubAuth(token string) error {
	var out map[string]any
	return getJSON(token, githubAPIBase+"/rate_limit", &out)
}
