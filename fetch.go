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

// defaultFormat is the section list used when --format is not passed.
const defaultFormat = "title,body"

// runFetch implements the `fetch` subcommand.
func runFetch(args []string) error {
	flags := flag.NewFlagSet("fetch", flag.ContinueOnError)
	var output, format string
	flags.StringVar(&output, "o", "", "write markdown to path")
	flags.StringVar(&output, "output", "", "write markdown to path")
	flags.StringVar(&format, "format", defaultFormat, "comma-separated sections to print, in order: metadata, title, body")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("usage: gip fetch [-o path] [--format=metadata,title,body] owner/repo#123")
	}
	var outputSet bool
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "o" || f.Name == "output" {
			outputSet = true
		}
	})
	if outputSet && output == "" {
		return fmt.Errorf("output path must not be empty")
	}

	sections, err := parseFormat(format)
	if err != nil {
		return err
	}

	owner, repo, number, err := parseRef(flags.Arg(0))
	if err != nil {
		return err
	}

	issue, err := FetchIssue(owner, repo, number)
	if err != nil {
		return err
	}

	markdown := issueMarkdown(issue, sections)
	if output == "" {
		_, err = fmt.Fprint(os.Stdout, markdown)
		return err
	}
	if err := os.WriteFile(output, []byte(markdown), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", output, err)
	}
	return nil
}

// parseFormat splits and validates a comma-separated --format value, e.g.
// "title,body" or "metadata,title,body".
func parseFormat(format string) ([]string, error) {
	sections := strings.Split(format, ",")
	for _, s := range sections {
		switch s {
		case "metadata", "title", "body":
		default:
			return nil, fmt.Errorf("invalid --format section %q: want metadata, title, or body", s)
		}
	}
	return sections, nil
}

// issueMarkdown renders the requested sections, concatenated in the given
// order.
func issueMarkdown(issue *Issue, sections []string) string {
	parts := make([]string, len(sections))
	for i, s := range sections {
		switch s {
		case "metadata":
			parts[i] = metadataSection(issue)
		case "title":
			parts[i] = fmt.Sprintf("# %s", issue.Title)
		case "body":
			parts[i] = bodySection(issue)
		}
	}
	return strings.Join(parts, "\n\n") + "\n"
}

// metadataSection renders the YAML front matter block.
func metadataSection(issue *Issue) string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %s\n", yamlString(issue.Title))
	fmt.Fprintf(&b, "number: %d\n", issue.Number)
	fmt.Fprintf(&b, "state: %s\n", yamlString(issue.State))
	fmt.Fprintf(&b, "author: %s\n", yamlString(issue.Author))
	fmt.Fprintf(&b, "url: %s\n", yamlString(issue.HTMLURL))
	fmt.Fprintf(&b, "created_at: %s\n", issue.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	fmt.Fprintf(&b, "updated_at: %s\n", issue.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	b.WriteString("---")
	return b.String()
}

// bodySection renders the issue/PR body plus any comments.
func bodySection(issue *Issue) string {
	var b strings.Builder
	b.WriteString(issue.Body)
	if len(issue.Comments) > 0 {
		b.WriteString("\n\n## Comments\n")
		for _, c := range issue.Comments {
			fmt.Fprintf(&b, "\n### %s — %s\n\n", c.Author, c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
			b.WriteString(c.Body)
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
