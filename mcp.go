package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type fetchInput struct {
	Ref    string `json:"ref" jsonschema:"GitHub issue or PR as owner/repo#123"`
	Format string `json:"format,omitempty" jsonschema:"Comma-separated sections: metadata, title, body; default title,body"`
}

func fetchTool(_ context.Context, _ *mcp.CallToolRequest, in fetchInput) (*mcp.CallToolResult, any, error) {
	if in.Format == "" {
		in.Format = defaultFormat
	}
	sections, err := parseFormat(in.Format)
	if err != nil {
		return nil, nil, err
	}
	owner, repo, number, err := parseRef(in.Ref)
	if err != nil {
		return nil, nil, err
	}
	issue, err := FetchIssue(owner, repo, number)
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{
		&mcp.TextContent{Text: issueMarkdown(issue, sections)},
	}}, nil, nil
}

func newMCPServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "gip", Version: "0.1.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "gip_fetch", Description: "Fetch a GitHub issue or PR as Markdown"}, fetchTool)
	return server
}

func mcpHTTPHandler(server *mcp.Server) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, nil))
	return mux
}

func runMCP(args []string) error {
	flags := flag.NewFlagSet("mcp", flag.ContinueOnError)
	transport := flags.String("transport", "stdio", "transport: stdio or http")
	addr := flags.String("addr", ":8080", "HTTP listen address")
	httpAddr := flags.String("http", "", "shorthand for --transport=http --addr=ADDRESS")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("usage: gip mcp [--transport=stdio|http] [--addr=:8080]")
	}
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "http" {
			*transport, *addr = "http", *httpAddr
		}
	})
	if *transport != "stdio" && *transport != "http" {
		return fmt.Errorf("invalid MCP transport %q: want stdio or http", *transport)
	}
	if *transport == "http" && *addr == "" {
		return fmt.Errorf("HTTP listen address must not be empty")
	}

	server := newMCPServer()
	if *transport == "stdio" {
		return server.Run(context.Background(), &mcp.StdioTransport{})
	}
	return http.ListenAndServe(*addr, mcpHTTPHandler(server))
}
