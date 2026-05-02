package tools

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

type ArtifactTools struct {
	root string
}

func NewArtifactTools(cfg config.AppConfig) *ArtifactTools {
	return &ArtifactTools{root: cfg.Dirs.Artifacts}
}

func (a *ArtifactTools) Register(s *mcpserver.MCPServer) {
	s.AddTool(mcp.NewTool("read_artifact",
		mcp.WithDescription("Read an artifact file from the artifacts/ directory. The path may point to markdown, yaml, json, or any other text artifact."),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Relative path under artifacts/ including extension when applicable (e.g. project/feature/discovery.md, project/feature/config.yaml)"),
		),
	), a.readArtifact)

	s.AddTool(mcp.NewTool("write_artifact",
		mcp.WithDescription("Write content to an artifact file, creating parent directories and Hugo _index.md files as needed. Agents may choose any suitable file extension. Use .md when the artifact is documentation; use .yaml, .json, .txt, image extensions, or another extension when that better matches the artifact."),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Relative path under artifacts/ including the chosen file extension. Documentation should use .md (e.g. project/feature/spec.md); config/data can use other extensions (e.g. project/feature/workflow.yaml)."),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("File content to write. For documentation artifacts, write markdown and use a .md path."),
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
	return filepath.Join(a.root, path)
}

func (a *ArtifactTools) readArtifact(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	if path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}
	full, err := a.safeArtifactPath(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	data, err := os.ReadFile(full)
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
	full, err := a.safeArtifactPath(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	parentDir := filepath.Dir(full)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot create parent directories: %v", err)), nil
	}
	createdIndexes, err := a.ensureArtifactIndexes(parentDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot create artifact indexes: %v", err)), nil
	}
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot write artifact: %v", err)), nil
	}
	rel, _ := filepath.Rel(a.root, full)
	msg := fmt.Sprintf("written %d bytes to artifacts/%s", len(content), rel)
	if len(createdIndexes) > 0 {
		msg += fmt.Sprintf("\ncreated indexes:\n  %s", strings.Join(createdIndexes, "\n  "))
	}
	return mcp.NewToolResultText(msg), nil
}

func (a *ArtifactTools) listArtifacts(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workflow := mcp.ParseArgument(req, "workflow", "").(string)
	root := a.root
	if workflow != "" {
		var err error
		root, err = a.safeArtifactPath(workflow)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	var sb strings.Builder
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if info.Name() == "_index.md" {
				return nil
			}
			rel, _ := filepath.Rel(a.root, path)
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

func (a *ArtifactTools) ensureArtifactIndexes(parentDir string) ([]string, error) {
	root := filepath.Clean(a.root)
	parentDir = filepath.Clean(parentDir)

	var dirs []string
	for dir := parentDir; ; dir = filepath.Dir(dir) {
		rel, err := filepath.Rel(root, dir)
		if err != nil {
			return nil, err
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return nil, fmt.Errorf("artifact path escapes %s", a.root)
		}
		dirs = append(dirs, dir)
		if dir == root {
			break
		}
	}

	var created []string
	for i := len(dirs) - 1; i >= 0; i-- {
		dir := dirs[i]
		indexPath := filepath.Join(dir, "_index.md")
		if _, err := os.Stat(indexPath); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return nil, err
		}

		content, err := a.artifactIndexContent(dir)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(indexPath, []byte(content), 0644); err != nil {
			return nil, err
		}
		rel, _ := filepath.Rel(a.root, indexPath)
		created = append(created, "artifacts/"+rel)
	}

	return created, nil
}

func (a *ArtifactTools) safeArtifactPath(path string) (string, error) {
	clean := filepath.Clean(path)
	if clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid artifact path %q", path)
	}
	return filepath.Join(a.root, clean), nil
}

func (a *ArtifactTools) artifactIndexContent(dir string) (string, error) {
	rel, err := filepath.Rel(a.root, dir)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return `---
title: "Artifacts"
weight: 10
---

# Artifacts

Generated workflow outputs.
`, nil
	}

	title := titleFromSlug(filepath.Base(dir))
	return fmt.Sprintf(`---
title: "%s"
weight: 10
---

# %s

Generated workflow outputs for %s.
`, title, title, rel), nil
}

func titleFromSlug(slug string) string {
	words := strings.Fields(strings.ReplaceAll(slug, "-", " "))
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = strings.ToUpper(word[:1]) + word[1:]
	}
	if len(words) == 0 {
		return slug
	}
	return strings.Join(words, " ")
}
