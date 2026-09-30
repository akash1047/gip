package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// parseRef parses the combined "owner/repo#123" argument.
func parseRef(ref string) (owner, repo string, number int, err error) {
	repoPart, numPart, ok := strings.Cut(ref, "#")
	if !ok {
		return "", "", 0, fmt.Errorf("invalid ref %q: expected owner/repo#123", ref)
	}
	owner, repo, ok = strings.Cut(repoPart, "/")
	if !ok || owner == "" || repo == "" {
		return "", "", 0, fmt.Errorf("invalid ref %q: expected owner/repo#123", ref)
	}
	number, err = strconv.Atoi(numPart)
	if err != nil || number <= 0 {
		return "", "", 0, fmt.Errorf("invalid ref %q: issue/PR number must be a positive integer", ref)
	}
	return owner, repo, number, nil
}

// runFetch implements the `fetch` subcommand: gip fetch [-o path] owner/repo#123.
func runFetch(args []string) error {
	flags := flag.NewFlagSet("fetch", flag.ContinueOnError)
	var output string
	flags.StringVar(&output, "o", "", "write markdown to path")
	flags.StringVar(&output, "output", "", "write markdown to path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("usage: gip fetch [-o path] owner/repo#123")
	}
	var outputSet bool
	flags.Visit(func(f *flag.Flag) { outputSet = true })
	if outputSet && output == "" {
		return fmt.Errorf("output path must not be empty")
	}

	owner, repo, number, err := parseRef(flags.Arg(0))
	if err != nil {
		return err
	}

	issue, err := FetchIssue(owner, repo, number)
	if err != nil {
		return err
	}

	markdown := issueMarkdown(issue)
	if output == "" {
		_, err = fmt.Fprint(os.Stdout, markdown)
		return err
	}
	if err := os.WriteFile(output, []byte(markdown), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", output, err)
	}
	return nil
}

// issueMarkdown renders an issue/PR as frontmatter + body, per README's
// "Idea" section.
func issueMarkdown(issue *Issue) string {
	var b strings.Builder

	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %s\n", yamlString(issue.Title))
	fmt.Fprintf(&b, "number: %d\n", issue.Number)
	fmt.Fprintf(&b, "state: %s\n", yamlString(issue.State))
	fmt.Fprintf(&b, "author: %s\n", yamlString(issue.Author))
	fmt.Fprintf(&b, "url: %s\n", yamlString(issue.HTMLURL))
	fmt.Fprintf(&b, "created_at: %s\n", issue.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	fmt.Fprintf(&b, "updated_at: %s\n", issue.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	b.WriteString("---\n\n")

	b.WriteString(issue.Body)
	b.WriteString("\n")

	if len(issue.Comments) > 0 {
		b.WriteString("\n## Comments\n")
		for _, c := range issue.Comments {
			fmt.Fprintf(&b, "\n### %s — %s\n\n", c.Author, c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
			b.WriteString(c.Body)
			b.WriteString("\n")
		}
	}

	return b.String()
}

// yamlString quotes a scalar for YAML frontmatter, escaping embedded quotes
// and backslashes.
func yamlString(s string) string {
	escaped := strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(s)
	return `"` + escaped + `"`
}
