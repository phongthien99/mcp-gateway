package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

const artifactsRoot = "artifacts"

type ArtifactTools struct{}

func NewArtifactTools() *ArtifactTools {
	return &ArtifactTools{}
}

func (a *ArtifactTools) Register(s *mcpserver.MCPServer) {
	s.AddTool(mcp.NewTool("read_artifact",
		mcp.WithDescription("Read an artifact file from the artifacts/ directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Relative path under artifacts/ (e.g. export-task-csv/discovery.md)"),
		),
	), a.readArtifact)

	s.AddTool(mcp.NewTool("write_artifact",
		mcp.WithDescription("Write content to an artifact file, creating parent directories as needed"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Relative path under artifacts/ (e.g. export-task-csv/spec.md)"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Markdown content to write"),
		),
	), a.writeArtifact)

	s.AddTool(mcp.NewTool("list_artifacts",
		mcp.WithDescription("List artifact files, optionally scoped to a workflow subdirectory"),
		mcp.WithString("workflow",
			mcp.Description("Workflow ID to scope the listing (e.g. export-task-csv)"),
		),
	), a.listArtifacts)
}

func (a *ArtifactTools) resolve(path string) string {
	return filepath.Join(artifactsRoot, path)
}

func (a *ArtifactTools) readArtifact(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	if path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}
	data, err := os.ReadFile(a.resolve(path))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot read artifact: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (a *ArtifactTools) writeArtifact(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	content := mcp.ParseArgument(req, "content", "").(string)
	if path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}
	full := a.resolve(path)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot create parent directories: %v", err)), nil
	}
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot write artifact: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("written %d bytes to artifacts/%s", len(content), path)), nil
}

func (a *ArtifactTools) listArtifacts(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workflow := mcp.ParseArgument(req, "workflow", "").(string)
	root := artifactsRoot
	if workflow != "" {
		root = filepath.Join(artifactsRoot, workflow)
	}

	var sb strings.Builder
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, _ := filepath.Rel(artifactsRoot, path)
			sb.WriteString(rel + "\n")
		}
		return nil
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot list artifacts: %v", err)), nil
	}
	if sb.Len() == 0 {
		return mcp.NewToolResultText("(no artifacts found)"), nil
	}
	return mcp.NewToolResultText(sb.String()), nil
}
