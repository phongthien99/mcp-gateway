package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const artifactsRoot = "artifacts"
const promptsRoot = "prompts"

// Runner executes a workflow definition step by step (dry-run / simulation mode).
// It proves the artifact chain: each step reads its inputs and reports what
// it would send to Claude, without making any API calls.
type Runner struct {
	def      Def
	workDir  string
	cache    map[string]string // artifact name → content
}

func NewRunner(workflowPath, projectID, featureID string) (*Runner, error) {
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return nil, fmt.Errorf("read workflow file: %w", err)
	}
	var def Def
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("parse workflow YAML: %w", err)
	}
	return &Runner{
		def:     def,
		workDir: filepath.Join(artifactsRoot, projectID, featureID),
		cache:   make(map[string]string),
	}, nil
}

func (r *Runner) Run(extraArgs map[string]string) error {
	if err := os.MkdirAll(r.workDir, 0755); err != nil {
		return fmt.Errorf("create artifacts dir: %w", err)
	}

	fmt.Printf("Workflow : %s\n", r.def.Name)
	fmt.Printf("Artifacts: %s\n", r.workDir)
	fmt.Println(strings.Repeat("─", 56))

	for i, step := range r.def.Steps {
		fmt.Printf("\n[%d/%d] %s\n", i+1, len(r.def.Steps), step.ID)
		if err := r.runStep(step, extraArgs); err != nil {
			return fmt.Errorf("step %q: %w", step.ID, err)
		}
	}

	fmt.Printf("\n%s\nAll steps complete.\n", strings.Repeat("─", 56))
	return nil
}

func (r *Runner) runStep(step Step, extraArgs map[string]string) error {
	// Load prompt template.
	promptPath := filepath.Join(promptsRoot, step.PromptFile)
	promptData, err := os.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("prompt file %q not found", promptPath)
	}

	// Build template vars.
	vars := make(map[string]string)
	for k, v := range extraArgs {
		vars[k] = v
	}

	// Read input artifacts.
	for _, name := range step.Reads {
		content, err := r.readArtifact(name)
		if err != nil {
			return err
		}
		vars[name] = content
		fmt.Printf("  ← read  : %s (%d chars)\n", name, len(content))
	}

	// Render prompt (replace placeholders).
	rendered := string(promptData)
	for k, v := range vars {
		rendered = strings.ReplaceAll(rendered, "{{"+k+"}}", v)
	}

	fmt.Printf("  prompt  : %s (%d chars rendered)\n", step.PromptFile, len(rendered))

	// In simulation mode: write a placeholder artifact so the next step can read it.
	if step.Writes != "" {
		placeholder := fmt.Sprintf("[simulation] output of step %q — replace with real Claude output.", step.ID)
		dest := filepath.Join(r.workDir, step.Writes+".md")
		if err := os.WriteFile(dest, []byte(placeholder), 0644); err != nil {
			return fmt.Errorf("write artifact: %w", err)
		}
		r.cache[step.Writes] = placeholder
		fmt.Printf("  → wrote : %s → %s\n", step.Writes, dest)
	}
	return nil
}

func (r *Runner) readArtifact(name string) (string, error) {
	if content, ok := r.cache[name]; ok {
		return content, nil
	}
	path := filepath.Join(r.workDir, name+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("artifact %q not found at %s — run previous steps first", name, path)
	}
	content := string(data)
	r.cache[name] = content
	return content, nil
}
