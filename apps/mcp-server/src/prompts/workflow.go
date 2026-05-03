package prompts

import (
	"context"
	"fmt"
	"log"
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
	prompts   string
	workflows string
}

func NewWorkflowPrompts(cfg config.AppConfig) *WorkflowPrompts {
	return &WorkflowPrompts{
		prompts:   cfg.Dirs.Prompts,
		workflows: cfg.Dirs.Workflows,
	}
}

func (w *WorkflowPrompts) Register(s *mcpserver.MCPServer) {
	// workflows are stored as {workflows_dir}/{name}/{name}.yaml
	files, err := filepath.Glob(filepath.Join(w.workflows, "*", "*.yaml"))
	if err != nil {
		log.Printf("[workflows] glob %q: %v", w.workflows, err)
		return
	}
	if len(files) == 0 {
		log.Printf("[workflows] no workflow yaml files found under %q", w.workflows)
		return
	}
	log.Printf("[workflows] found %d workflow file(s) under %q", len(files), w.workflows)
	for _, f := range files {
		w.registerWorkflow(s, f)
	}
}

func (w *WorkflowPrompts) registerWorkflow(s *mcpserver.MCPServer, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[workflows] read %q: %v", path, err)
		return
	}
	var def workflow.Def
	if err := yaml.Unmarshal(data, &def); err != nil {
		log.Printf("[workflows] parse %q: %v", path, err)
		return
	}

	for i, step := range def.Steps {
		step := step // capture loop var
		nextStepID := ""
		if i+1 < len(def.Steps) {
			nextStepID = def.Steps[i+1].ID
		}

		promptPath := filepath.Join(w.prompts, step.PromptFile)
		if _, err := os.Stat(promptPath); err != nil {
			log.Printf("[prompts] MISSING %q (step %q in workflow %q)", promptPath, step.ID, def.ID)
		} else {
			log.Printf("[prompts] ok %q (step %q)", promptPath, step.ID)
		}

		opts := []mcp.PromptOption{
			mcp.WithPromptDescription(step.Description),
		}
		for _, arg := range step.ExtraArgs {
			argOpts := []mcp.ArgumentOption{mcp.ArgumentDescription(arg.Description)}
			if arg.Required {
				argOpts = append(argOpts, mcp.RequiredArgument())
			}
			opts = append(opts, mcp.WithArgument(arg.Name, argOpts...))
		}

		s.AddPrompt(mcp.NewPrompt(step.ID, opts...), w.buildHandler(step, nextStepID))
	}
	log.Printf("[workflows] registered %q: %d step(s)", def.ID, len(def.Steps))
}

func (w *WorkflowPrompts) buildHandler(step workflow.Step, nextStepID string) mcpserver.PromptHandlerFunc {
	return func(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := req.Params.Arguments

		vars := map[string]string{
			"markdown_rules": w.readPromptInclude(markdownRulesFile),
			"step_id":        step.ID,
			"reads":          strings.Join(step.Reads, ", "),
			"writes":         step.Writes,
			"context_docs":   strings.Join(step.Context, ", "),
			"next_step":      nextStepID,
		}
		for _, arg := range step.ExtraArgs {
			vars[arg.Name] = args[arg.Name]
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
	log.Printf("[prompts] loading %q", path)
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

