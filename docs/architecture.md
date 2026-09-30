# Architecture

GIP is a single-package Go CLI. `main` dispatches to `fetch`, `doctor`, or
`mcp`. The CLI and MCP interfaces share the same issue/PR fetching and
Markdown formatting flow.

`fetch` parses a repository issue/PR reference and requested sections.
`FetchIssue` retrieves the issue or PR and its comments from the GitHub REST
API; `issueMarkdown` turns the result into Markdown for stdout or a file.
The `mcp` command exposes that flow as `gip_fetch` over stdio or HTTP.

The auth layer reads `GIP_GITHUB_TOKEN`, falling back to `GITHUB_TOKEN`.
Fetching public repositories works without a token. `doctor` requires one
and checks it against GitHub's `/rate_limit` endpoint.

```mermaid
flowchart LR
    CLI["gip CLI"] -->|fetch| F["Fetch issue or PR"]
    CLI -->|doctor| D["Auth check"]
    CLI -->|mcp| M["MCP server"]
    M -->|gip_fetch| F
    F -->|optional| A["Environment token"]
    D --> A
    F --> G["GitHub REST API"]
    A --> G
    D -->|verify token| G
    G -->|issue and comments| R["Markdown formatting"]
    R --> O["CLI output or MCP response"]
```
