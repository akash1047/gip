package main

import (
	"strings"
	"testing"
	"time"
)

func TestParseRef(t *testing.T) {
	cases := []struct {
		ref         string
		owner, repo string
		number      int
		wantErr     bool
	}{
		{ref: "octocat/hello-world#123", owner: "octocat", repo: "hello-world", number: 123},
		{ref: "no-hash", wantErr: true},
		{ref: "no-slash#1", wantErr: true},
		{ref: "owner/repo#abc", wantErr: true},
		{ref: "owner/repo#0", wantErr: true},
	}

	for _, tc := range cases {
		owner, repo, number, err := parseRef(tc.ref)
		if tc.wantErr {
			if err == nil {
				t.Errorf("parseRef(%q): expected error", tc.ref)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseRef(%q): unexpected error: %v", tc.ref, err)
			continue
		}
		if owner != tc.owner || repo != tc.repo || number != tc.number {
			t.Errorf("parseRef(%q) = %q, %q, %d; want %q, %q, %d",
				tc.ref, owner, repo, number, tc.owner, tc.repo, tc.number)
		}
	}
}

func TestIssueMarkdown(t *testing.T) {
	issue := &Issue{
		Number:    42,
		Title:     `Title with "quotes"`,
		Body:      "the body",
		State:     "open",
		Author:    "alice",
		HTMLURL:   "https://github.com/o/r/issues/42",
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		Comments: []Comment{
			{Author: "bob", Body: "a comment", CreatedAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)},
		},
	}

	got := issueMarkdown(issue)

	for _, want := range []string{
		`title: "Title with \"quotes\""`,
		"number: 42",
		"the body",
		"## Comments",
		"### bob — 2024-01-03T00:00:00Z",
		"a comment",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("issueMarkdown() missing %q, got:\n%s", want, got)
		}
	}
}
