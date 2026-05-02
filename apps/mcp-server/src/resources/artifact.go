package resources

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mcp-gateway/src/config"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// ArtifactResources exposes workflow artifacts as MCP resources.
// URI scheme: resource://artifact/{project}/{feature}/{name}
// Maps to:    {root}/{project}/{feature}/{name}.md on disk
type ArtifactResources struct {
	root string
}

func NewArtifactResources(cfg config.AppConfig) *ArtifactResources {
	return &ArtifactResources{root: cfg.Dirs.Artifacts}
}

func (a *ArtifactResources) Register(s *mcpserver.MCPServer) {
	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"resource://artifact/{project}/{feature}/{name}",
			"Workflow Artifact",
			mcp.WithTemplateMIMEType("text/markdown"),
			mcp.WithTemplateDescription("A generated artifact file scoped by project and workflow step"),
		),
		a.read,
	)
}

func (a *ArtifactResources) read(_ context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	uri := req.Params.URI

	// resource://artifact/{project}/{feature}/{name}  →  {root}/{project}/{feature}/{name}.md
	trimmed := strings.TrimPrefix(uri, "resource://artifact/")
	parts := strings.SplitN(trimmed, "/", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid artifact URI %q — expected resource://artifact/{project}/{feature}/{name}", uri)
	}
	project, feature, name := parts[0], parts[1], parts[2]

	path := filepath.Join(a.root, project, feature, name+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("artifact not found at %s: %w", path, err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "text/markdown",
			Text:     string(data),
		},
	}, nil
}
