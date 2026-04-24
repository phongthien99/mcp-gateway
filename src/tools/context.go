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

const contextRoot = "context"

// contextTemplates are the standard context doc names with their default content.
var contextTemplates = map[string]string{
	"architecture": `# Architecture Guidelines

## Stack
<!-- Mô tả stack công nghệ: backend, frontend, database, cache -->

## Layer structure
<!-- Mô tả cấu trúc tầng: handler → service → repository -->

## API conventions
<!-- RESTful, error format, pagination, auth... -->

## Database
<!-- ORM/query builder, transaction policy, index strategy -->
`,
	"coding-standards": `# Coding Standards

## Language / Framework
<!-- Naming conventions, error handling, testing rules -->

## Git
<!-- Branch naming, commit message format, PR requirements -->

## Testing
<!-- Coverage targets, integration test policy -->
`,
	"compliance-rules": `# Compliance Rules

<!-- Quy tắc bảo mật, data privacy, logging policy, ... -->
`,
}

type ContextTools struct{}

func NewContextTools() *ContextTools { return &ContextTools{} }

func (c *ContextTools) Register(s *mcpserver.MCPServer) {
	s.AddTool(mcp.NewTool("init_project",
		mcp.WithDescription("Initialise a project context directory at context/{project_id}/ and write template files for any standard context docs that do not yet exist (architecture, coding-standards, compliance-rules). Returns the list of files created."),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("Unique project identifier, used as the directory name under context/"),
		),
	), c.initProject)

	s.AddTool(mcp.NewTool("set_context",
		mcp.WithDescription("Write or overwrite a context document for a project at context/{project_id}/{name}.md"),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("Project identifier"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Context doc name without extension (e.g. architecture, coding-standards, compliance-rules, or any custom name)"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Markdown content to write"),
		),
	), c.setContext)

	s.AddTool(mcp.NewTool("get_context",
		mcp.WithDescription("Read a context document. Looks in context/{project_id}/{name}.md first, then context/global/{name}.md"),
		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("Project identifier"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Context doc name without extension"),
		),
	), c.getContext)

	s.AddTool(mcp.NewTool("list_context",
		mcp.WithDescription("List all context documents available to a project (project-scoped and global, deduplicated)"),
		mcp.WithString("project_id",
			mcp.Description("Project identifier — omit to list only global context docs"),
		),
	), c.listContext)
}

func (c *ContextTools) initProject(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := mcp.ParseArgument(req, "project_id", "").(string)
	if projectID == "" {
		return mcp.NewToolResultError("project_id is required"), nil
	}

	dir := filepath.Join(contextRoot, projectID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot create context dir: %v", err)), nil
	}

	var created []string
	for name, tmpl := range contextTemplates {
		path := filepath.Join(dir, name+".md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte(tmpl), 0644); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("cannot write %s: %v", name, err)), nil
			}
			created = append(created, "context/"+projectID+"/"+name+".md")
		}
	}

	if len(created) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("project %q already initialised — no new files created", projectID)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("initialised project %q\ncreated:\n  %s", projectID, strings.Join(created, "\n  "))), nil
}

func (c *ContextTools) setContext(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := mcp.ParseArgument(req, "project_id", "").(string)
	name := mcp.ParseArgument(req, "name", "").(string)
	content := mcp.ParseArgument(req, "content", "").(string)
	if projectID == "" || name == "" {
		return mcp.NewToolResultError("project_id and name are required"), nil
	}

	dir := filepath.Join(contextRoot, projectID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot create context dir: %v", err)), nil
	}

	path := filepath.Join(dir, name+".md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot write context: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("written %d bytes to context/%s/%s.md", len(content), projectID, name)), nil
}

func (c *ContextTools) getContext(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := mcp.ParseArgument(req, "project_id", "").(string)
	name := mcp.ParseArgument(req, "name", "").(string)
	if projectID == "" || name == "" {
		return mcp.NewToolResultError("project_id and name are required"), nil
	}

	projectPath := filepath.Join(contextRoot, projectID, name+".md")
	if data, err := os.ReadFile(projectPath); err == nil {
		return mcp.NewToolResultText(string(data)), nil
	}

	globalPath := filepath.Join(contextRoot, "global", name+".md")
	if data, err := os.ReadFile(globalPath); err == nil {
		return mcp.NewToolResultText(string(data)), nil
	}

	return mcp.NewToolResultError(fmt.Sprintf("context %q not found for project %q or in global", name, projectID)), nil
}

func (c *ContextTools) listContext(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := mcp.ParseArgument(req, "project_id", "").(string)

	seen := map[string]string{} // name → source label

	if projectID != "" {
		dir := filepath.Join(contextRoot, projectID)
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				name := strings.TrimSuffix(e.Name(), ".md")
				seen[name] = "project"
			}
		}
	}

	globalDir := filepath.Join(contextRoot, "global")
	entries, _ := os.ReadDir(globalDir)
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			name := strings.TrimSuffix(e.Name(), ".md")
			if _, ok := seen[name]; !ok {
				seen[name] = "global"
			}
		}
	}

	if len(seen) == 0 {
		return mcp.NewToolResultText("(no context docs found)"), nil
	}

	var sb strings.Builder
	for name, src := range seen {
		sb.WriteString(fmt.Sprintf("%s  [%s]\n", name, src))
	}
	return mcp.NewToolResultText(strings.TrimRight(sb.String(), "\n")), nil
}
