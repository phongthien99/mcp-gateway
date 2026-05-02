package prompts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mcp-gateway/src/config"
	"mcp-gateway/src/workflow"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"
)

const markdownRulesFile = "markdown_rules.md"

// WorkflowPrompts scans workflows/*.yaml at startup and dynamically registers
// one MCP prompt per step. Adding a new step or workflow requires only YAML +
// a .md file — no Go changes needed.
type WorkflowPrompts struct {
	artifacts string
	prompts   string
	context   string
	workflows string
}

func NewWorkflowPrompts(cfg config.AppConfig) *WorkflowPrompts {
	return &WorkflowPrompts{
		artifacts: cfg.Dirs.Artifacts,
		prompts:   cfg.Dirs.Prompts,
		context:   cfg.Dirs.Context,
		workflows: cfg.Dirs.Workflows,
	}
}

func (w *WorkflowPrompts) Register(s *mcpserver.MCPServer) {
	files, err := filepath.Glob(filepath.Join(w.workflows, "*.yaml"))
	if err != nil || len(files) == 0 {
		return
	}
	for _, f := range files {
		w.registerWorkflow(s, f)
	}
}

func (w *WorkflowPrompts) registerWorkflow(s *mcpserver.MCPServer, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var def workflow.Def
	if err := yaml.Unmarshal(data, &def); err != nil {
		return
	}

	for _, step := range def.Steps {
		step := step // capture loop var

		opts := []mcp.PromptOption{
			mcp.WithPromptDescription(step.Description),
			mcp.WithArgument("project_id",
				mcp.ArgumentDescription("ID project để namespace artifacts (vd: my-app)"),
				mcp.RequiredArgument(),
			),
			mcp.WithArgument("feature_id",
				mcp.ArgumentDescription("Tên feature đang phát triển (vd: export-task-csv)"),
				mcp.RequiredArgument(),
			),
		}
		for _, arg := range step.ExtraArgs {
			argOpts := []mcp.ArgumentOption{mcp.ArgumentDescription(arg.Description)}
			if arg.Required {
				argOpts = append(argOpts, mcp.RequiredArgument())
			}
			opts = append(opts, mcp.WithArgument(arg.Name, argOpts...))
		}

		s.AddPrompt(mcp.NewPrompt(step.ID, opts...), w.buildHandler(step))
	}
}

func (w *WorkflowPrompts) buildHandler(step workflow.Step) mcpserver.PromptHandlerFunc {
	return func(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := req.Params.Arguments
		projectID := args["project_id"]
		featureID := args["feature_id"]

		vars := map[string]string{
			"project_id":     projectID,
			"feature_id":     featureID,
			"markdown_rules": w.readPromptInclude(markdownRulesFile),
		}

		// Extra args forwarded as template vars.
		for _, arg := range step.ExtraArgs {
			vars[arg.Name] = args[arg.Name]
		}

		// Load generated artifacts from previous steps.
		for _, name := range step.Reads {
			content, err := w.readArtifact(projectID, featureID, name)
			if err != nil {
				return nil, err
			}
			vars[name] = content
		}

		// Load reference/context docs: project-scoped first, fallback to global.
		for _, name := range step.Context {
			content, err := w.readContext(projectID, name)
			if err != nil {
				return nil, err
			}
			vars[name] = content
		}

		text, err := w.renderPrompt(step.PromptFile, vars)
		if err != nil {
			return nil, err
		}

		return &mcp.GetPromptResult{
			Description: step.Description,
			Messages: []mcp.PromptMessage{
				{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: text}},
			},
		}, nil
	}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func (w *WorkflowPrompts) renderPrompt(promptFile string, vars map[string]string) (string, error) {
	path := filepath.Join(w.prompts, promptFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("prompt file %q not found: %w", path, err)
	}
	text := string(data)
	for k, v := range vars {
		text = strings.ReplaceAll(text, "{{"+k+"}}", v)
	}
	return text, nil
}

func (w *WorkflowPrompts) readPromptInclude(name string) string {
	data, err := os.ReadFile(filepath.Join(w.prompts, name))
	if err != nil {
		return ""
	}
	return string(data)
}

func (w *WorkflowPrompts) readArtifact(projectID, featureID, name string) (string, error) {
	path := filepath.Join(w.artifacts, projectID, featureID, name+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("artifact %q not found (project=%s, feature=%s) — run previous steps first",
			name, projectID, featureID)
	}
	return string(data), nil
}

// readContext loads a reference doc, preferring project-scoped over global.
//
//	{context}/{project_id}/{name}.md  →  project-specific version
//	{context}/global/{name}.md        →  fallback shared version
func (w *WorkflowPrompts) readContext(projectID, name string) (string, error) {
	projectPath := filepath.Join(w.context, projectID, name+".md")
	if data, err := os.ReadFile(projectPath); err == nil {
		return string(data), nil
	}

	globalPath := filepath.Join(w.context, "global", name+".md")
	if data, err := os.ReadFile(globalPath); err == nil {
		return string(data), nil
	}

	return "", fmt.Errorf("context %q not found at %s or %s", name, projectPath, globalPath)
}
