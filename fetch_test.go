package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type testTransport func(*http.Request) (*http.Response, error)

func (f testTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

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

func TestRunFetchOutput(t *testing.T) {
	t.Setenv(tokenEnvVar, "test-token")
	origTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = testTransport(func(r *http.Request) (*http.Response, error) {
		body := `{"number":42,"title":"test issue","body":"test body","state":"open","user":{"login":"alice"}}`
		if strings.HasSuffix(r.URL.Path, "/comments") {
			body = `[]`
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	defer func() { http.DefaultClient.Transport = origTransport }()

	for _, tc := range []struct {
		name, flag string
		args       []string
		want       string
	}{
		{name: "stdout"},
		{name: "short flag", flag: "-o"},
		{name: "long flag", flag: "--output"},
		{name: "body only", args: []string{"--format=body"}, want: "test body\n"},
		{name: "title and body", args: []string{"--format=title,body"}, want: "# test issue\n\ntest body\n"},
		{name: "body and title reordered", args: []string{"--format=body,title"}, want: "test body\n\n# test issue\n"},
		{name: "metadata and body", args: []string{"--format=metadata,body"}, want: "---\n\ntest body\n"},
		{name: "metadata and title", args: []string{"--format=metadata,title"}, want: "---\n\n# test issue\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			args := append(append([]string{}, tc.args...), "o/r#42")
			if tc.flag != "" {
				args = []string{tc.flag, "issue.md", "o/r#42"}
			}
			reader, writer, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			origStdout := os.Stdout
			os.Stdout = writer
			defer func() { os.Stdout = origStdout; reader.Close(); writer.Close() }()

			if err := runFetch(args); err != nil {
				t.Fatal(err)
			}
			writer.Close()
			stdout, err := io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}
			if tc.flag == "" {
				if tc.want != "" {
					got := string(stdout)
					if strings.HasPrefix(tc.want, "---") {
						if !strings.HasPrefix(got, "---\n") || !strings.HasSuffix(got, tc.want) {
							t.Errorf("stdout = %q, want metadata and suffix %q", got, tc.want)
						}
					} else if got != tc.want {
						t.Errorf("stdout = %q, want %q", got, tc.want)
					}
				} else if !strings.Contains(string(stdout), "# test issue\n\ntest body\n") {
					t.Errorf("stdout = %q, want default title+body markdown", stdout)
				}
				if _, err := os.Stat("o-r-42.md"); !os.IsNotExist(err) {
					t.Errorf("default output file exists or stat failed: %v", err)
				}
				return
			}
			if len(stdout) != 0 {
				t.Errorf("stdout = %q, want empty", stdout)
			}
			data, err := os.ReadFile("issue.md")
			if err != nil || !strings.Contains(string(data), "test body") {
				t.Errorf("output file = %q, error = %v", data, err)
			}
		})
	}
}

func TestParseFormat(t *testing.T) {
	if got, err := parseFormat("body,title"); err != nil || strings.Join(got, ",") != "body,title" {
		t.Errorf("parseFormat(%q) = %q, %v", "body,title", got, err)
	}
	if _, err := parseFormat("body,bogus"); err == nil {
		t.Errorf("parseFormat(%q): expected error", "body,bogus")
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

	got := issueMarkdown(issue, []string{"metadata", "title", "body"})

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
	if got := issueMarkdown(issue, []string{"metadata", "title"}); strings.Contains(got, "the body") || strings.Contains(got, "a comment") {
		t.Errorf("issueMarkdown without body includes body or comments: %q", got)
	}

	if got := issueMarkdown(issue, []string{"body", "title"}); !strings.HasPrefix(got, "the body") {
		t.Errorf("issueMarkdown(body,title) = %q, want body before title", got)
	}
}
