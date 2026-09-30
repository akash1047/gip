package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// githubAPIBase is a var, not a const, so tests can point it at an
// httptest server.
var githubAPIBase = "https://api.github.com"

// Comment is a single issue/PR comment, trimmed to what's needed for the
// markdown output (issue #8).
type Comment struct {
	Author    string
	Body      string
	CreatedAt time.Time
}

// Issue holds a GitHub issue or PR's title, body, comments, and metadata.
// The GitHub REST API serves PRs through the issues endpoint too, so this
// struct covers both.
type Issue struct {
	Number    int
	Title     string
	Body      string
	State     string
	Author    string
	HTMLURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Comments  []Comment
}

// userWire is the subset of GitHub's user object needed here.
type userWire struct {
	Login string `json:"login"`
}

type issueWire struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	State     string    `json:"state"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      userWire  `json:"user"`
}

type commentWire struct {
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	User      userWire  `json:"user"`
}

// FetchIssue retrieves an issue or PR (title, body, metadata) plus its
// comments from the GitHub REST API.
func FetchIssue(owner, repo string, number int) (*Issue, error) {
	token, _ := githubToken() // Public issues work without a token.

	var iw issueWire
	issueURL := fmt.Sprintf("%s/repos/%s/%s/issues/%d", githubAPIBase, owner, repo, number)
	if err := getJSON(token, issueURL, &iw); err != nil {
		return nil, fmt.Errorf("fetching issue: %w", err)
	}

	var cw []commentWire
	commentsURL := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", githubAPIBase, owner, repo, number)
	if err := getJSON(token, commentsURL, &cw); err != nil {
		return nil, fmt.Errorf("fetching comments: %w", err)
	}

	comments := make([]Comment, len(cw))
	for i, c := range cw {
		comments[i] = Comment{Author: c.User.Login, Body: c.Body, CreatedAt: c.CreatedAt}
	}

	return &Issue{
		Number:    iw.Number,
		Title:     iw.Title,
		Body:      iw.Body,
		State:     iw.State,
		Author:    iw.User.Login,
		HTMLURL:   iw.HTMLURL,
		CreatedAt: iw.CreatedAt,
		UpdatedAt: iw.UpdatedAt,
		Comments:  comments,
	}, nil
}

// getJSON performs a GET and decodes the JSON response body into out.
func getJSON(token, url string, out any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: unexpected status %s", url, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
