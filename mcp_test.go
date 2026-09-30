package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPFetchTool(t *testing.T) {
	t.Setenv(tokenEnvVar, "test-token")
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/issues/42", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing GitHub token")
		}
		json.NewEncoder(w).Encode(map[string]any{"number": 42, "title": "hello", "body": "world"})
	})
	mux.HandleFunc("/repos/o/r/issues/42/comments", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]any{})
	})
	mockGithub(t, mux)

	ctx := context.Background()
	serverSide, clientSide := mcp.NewInMemoryTransports()
	ss, err := newMCPServer().Connect(ctx, serverSide, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ss.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0"}, nil)
	cs, err := client.Connect(ctx, clientSide, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()
	tools, err := cs.ListTools(ctx, nil)
	if err != nil || len(tools.Tools) != 1 || tools.Tools[0].Name != "gip_fetch" {
		t.Fatalf("tools = %+v, %v", tools, err)
	}
	for _, tc := range []struct {
		name string
		args map[string]any
		want string
		err  bool
	}{
		{"default", map[string]any{"ref": "o/r#42"}, "# hello\n\nworld\n", false},
		{"body", map[string]any{"ref": "o/r#42", "format": "body"}, "world\n", false},
		{"invalid ref", map[string]any{"ref": "bad"}, "", true},
		{"invalid format", map[string]any{"ref": "o/r#42", "format": "bad"}, "", true},
		{"missing ref", map[string]any{}, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "gip_fetch", Arguments: tc.args})
			if err != nil {
				t.Fatal(err)
			}
			if result.IsError != tc.err {
				t.Fatalf("IsError = %v, want %v", result.IsError, tc.err)
			}
			if !tc.err {
				if len(result.Content) != 1 || result.Content[0].(*mcp.TextContent).Text != tc.want {
					t.Errorf("content = %+v, want %q", result.Content, tc.want)
				}
			}
		})
	}
}

func TestMCPFlags(t *testing.T) {
	for _, args := range [][]string{
		{"--transport=bogus"}, {"--transport=http", "--addr="}, {"--http="}, {"extra"},
	} {
		if err := runMCP(args); err == nil {
			t.Errorf("runMCP(%q): expected error", strings.Join(args, " "))
		}
	}
}

func TestMCPHTTP(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json, text/event-stream")
	response := httptest.NewRecorder()
	mcpHTTPHandler(newMCPServer()).ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"name":"gip"`) {
		t.Fatalf("HTTP status = %d, body = %s", response.Code, response.Body.String())
	}
}
