package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mcp-gateway/src/config"
	"mcp-gateway/src/scope"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type ArtifactTools struct {
	projectsRoot string
}

func NewArtifactTools(cfg config.AppConfig) *ArtifactTools {
	return &ArtifactTools{projectsRoot: cfg.Dirs.ProjectsRoot}
}

func (a *ArtifactTools) Register(s *mcpserver.MCPServer) {
	s.AddTool(mcp.NewTool("read_artifact",
		mcp.WithDescription("Read an artifact file. Path format: {project_id}/{feature}/{file} (e.g. react-mui/user-management/spec.md)"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Relative path as {project_id}/{feature}/{file} including extension"),
		),
	), a.readArtifact)

	s.AddTool(mcp.NewTool("write_artifact",
		mcp.WithDescription("Write content to an artifact file, creating parent directories and Hugo _index.md files as needed. Path format: {project_id}/{feature}/{file}"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Relative path as {project_id}/{feature}/{file} including extension (e.g. react-mui/user-management/spec.md)"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("File content to write. Use .md for documentation artifacts."),
		),
	), a.writeArtifact)

	s.AddTool(mcp.NewTool("list_artifacts",
		mcp.WithDescription("List artifact files, optionally scoped to a project"),
		mcp.WithString("workflow",
			mcp.Description("Project ID to scope the listing (e.g. react-mui)"),
		),
	), a.listArtifacts)
}

// artifactsDir returns the artifacts subdirectory for a project.
func (a *ArtifactTools) artifactsDir(projectID string) string {
	return filepath.Join(a.projectsRoot, projectID, "artifacts")
}

// resolve maps {project_id}/{feature}/{file} → {projectsRoot}/{project_id}/artifacts/{feature}/{file}.
func (a *ArtifactTools) resolve(path string) (string, error) {
	clean := filepath.Clean(path)
	if clean == "." || filepath.IsAbs(clean) || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid artifact path %q", path)
	}
	parts := strings.SplitN(clean, string(os.PathSeparator), 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("artifact path must be {project_id}/{feature}/{file}, got %q", path)
	}
	projectID := parts[0]
	rest := parts[1]
	return filepath.Join(a.artifactsDir(projectID), rest), nil
}

func (a *ArtifactTools) readArtifact(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	if path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}
	full, err := a.resolve(path)
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
	full, err := a.resolve(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	parentDir := filepath.Dir(full)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot create parent directories: %v", err)), nil
	}

	// extract projectID to scope index creation
	parts := strings.SplitN(filepath.Clean(path), string(os.PathSeparator), 2)
	projectID := parts[0]
	createdIndexes, err := a.ensureArtifactIndexes(projectID, parentDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot create artifact indexes: %v", err)), nil
	}
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot write artifact: %v", err)), nil
	}
	artifactsRoot := a.artifactsDir(projectID)
	rel, _ := filepath.Rel(artifactsRoot, full)
	msg := fmt.Sprintf("written %d bytes to %s/artifacts/%s", len(content), projectID, rel)
	if len(createdIndexes) > 0 {
		msg += fmt.Sprintf("\ncreated indexes:\n  %s", strings.Join(createdIndexes, "\n  "))
	}
	return mcp.NewToolResultText(msg), nil
}

func (a *ArtifactTools) listArtifacts(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workflow := mcp.ParseArgument(req, "workflow", "").(string)
	// If no explicit workflow and a project scope is set on the connection, use it.
	if workflow == "" {
		workflow = scope.FromContext(ctx)
	}

	// Returns paths as {project_id}/{feature}/{file} — no "artifacts/" segment —
	// so the result can be passed directly to read_artifact.
	var sb strings.Builder
	if workflow != "" {
		root := a.artifactsDir(workflow)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Name() != "_index.md" {
				rel, _ := filepath.Rel(root, path)
				sb.WriteString(filepath.ToSlash(filepath.Join(workflow, rel)) + "\n")
			}
			return nil
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot list artifacts: %v", err)), nil
		}
	} else {
		entries, err := os.ReadDir(a.projectsRoot)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot read projects root: %v", err)), nil
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			projectID := e.Name()
			root := a.artifactsDir(projectID)
			_ = filepath.Walk(root, func(path string, info os.FileInfo, werr error) error {
				if werr != nil || info.IsDir() || info.Name() == "_index.md" {
					return nil
				}
				rel, _ := filepath.Rel(root, path)
				sb.WriteString(filepath.ToSlash(filepath.Join(projectID, rel)) + "\n")
				return nil
			})
		}
	}

	if sb.Len() == 0 {
		return mcp.NewToolResultText("(no artifacts found)"), nil
	}
	return mcp.NewToolResultText(sb.String()), nil
}

// ensureArtifactIndexes creates missing _index.md files from parentDir up to
// and including the project's artifacts/ dir.
func (a *ArtifactTools) ensureArtifactIndexes(projectID, parentDir string) ([]string, error) {
	artifactsRoot := filepath.Clean(a.artifactsDir(projectID))
	parentDir = filepath.Clean(parentDir)

	// ensure project _index.md
	projectDir := filepath.Join(a.projectsRoot, projectID)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, err
	}
	projectIndex := filepath.Join(projectDir, "_index.md")
	if _, err := os.Stat(projectIndex); os.IsNotExist(err) {
		if err := os.WriteFile(projectIndex, []byte(projectIndexTemplate(projectID)), 0644); err != nil {
			return nil, err
		}
	}

	var dirs []string
	for dir := parentDir; ; dir = filepath.Dir(dir) {
		rel, err := filepath.Rel(artifactsRoot, dir)
		if err != nil {
			return nil, err
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			break
		}
		dirs = append(dirs, dir)
		if dir == artifactsRoot {
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
		content, err := a.artifactIndexContent(projectID, dir)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(indexPath, []byte(content), 0644); err != nil {
			return nil, err
		}
		rel, _ := filepath.Rel(a.artifactsDir(projectID), indexPath)
		created = append(created, projectID+"/artifacts/"+rel)
	}
	return created, nil
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

func (a *ArtifactTools) artifactIndexContent(projectID, dir string) (string, error) {
	artifactsRoot := a.artifactsDir(projectID)
	rel, err := filepath.Rel(artifactsRoot, dir)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return `---
title: "Artifacts"
weight: 20
bookCollapseSection: true
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
