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
const contextRoot = "context"
const markdownRulesFile = "markdown_rules.md"

// Runner executes a workflow in dry-run / simulation mode.
// All configuration is read from a RunConfig file — no params passed in code.
type Runner struct {
	cfg     RunConfig
	def     Def
	workDir string
	cache   map[string]string // artifact name → content
}

// NewRunner reads the run config file and the workflow YAML it references.
func NewRunner(runConfigPath string) (*Runner, error) {
	cfgData, err := os.ReadFile(runConfigPath)
	if err != nil {
		return nil, fmt.Errorf("read run config: %w", err)
	}
	var cfg RunConfig
	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		return nil, fmt.Errorf("parse run config: %w", err)
	}
	if cfg.Workflow == "" {
		return nil, fmt.Errorf("run config missing required field: workflow")
	}
	if cfg.ProjectID == "" || cfg.FeatureID == "" {
		return nil, fmt.Errorf("run config missing required fields: project_id, feature_id")
	}

	wfData, err := os.ReadFile(cfg.Workflow)
	if err != nil {
		return nil, fmt.Errorf("read workflow file %q: %w", cfg.Workflow, err)
	}
	var def Def
	if err := yaml.Unmarshal(wfData, &def); err != nil {
		return nil, fmt.Errorf("parse workflow YAML: %w", err)
	}

	return &Runner{
		cfg:     cfg,
		def:     def,
		workDir: filepath.Join(artifactsRoot, cfg.ProjectID, cfg.FeatureID),
		cache:   make(map[string]string),
	}, nil
}

func (r *Runner) Run() error {
	if err := os.MkdirAll(r.workDir, 0755); err != nil {
		return fmt.Errorf("create artifacts dir: %w", err)
	}

	fmt.Printf("Workflow  : %s\n", r.def.Name)
	fmt.Printf("Project   : %s\n", r.cfg.ProjectID)
	fmt.Printf("Feature   : %s\n", r.cfg.FeatureID)
	fmt.Printf("Artifacts : %s\n", r.workDir)
	fmt.Println(strings.Repeat("─", 56))

	for i, step := range r.def.Steps {
		fmt.Printf("\n[%d/%d] %s\n", i+1, len(r.def.Steps), step.ID)
		if err := r.runStep(step); err != nil {
			return fmt.Errorf("step %q: %w", step.ID, err)
		}
	}

	fmt.Printf("\n%s\nAll steps complete.\n", strings.Repeat("─", 56))
	return nil
}

func (r *Runner) runStep(step Step) error {
	// Load prompt template.
	promptPath := filepath.Join(promptsRoot, step.PromptFile)
	promptData, err := os.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("prompt file %q not found", promptPath)
	}

	// Build template vars: fixed context + run config args.
	vars := map[string]string{
		"project_id":     r.cfg.ProjectID,
		"feature_id":     r.cfg.FeatureID,
		"markdown_rules": readPromptInclude(markdownRulesFile),
	}
	for k, v := range r.cfg.Args {
		vars[k] = v
	}

	// Read generated artifacts from previous steps.
	for _, name := range step.Reads {
		content, err := r.readArtifact(name)
		if err != nil {
			return err
		}
		vars[name] = content
		fmt.Printf("  ← artifact: %s (%d chars)\n", name, len(content))
	}

	// Read reference/context docs: project-scoped first, fallback to global.
	for _, name := range step.Context {
		content, err := r.readContext(name)
		if err != nil {
			return err
		}
		vars[name] = content
		fmt.Printf("  ← context : %s (%d chars)\n", name, len(content))
	}

	// Render prompt.
	rendered := string(promptData)
	for k, v := range vars {
		rendered = strings.ReplaceAll(rendered, "{{"+k+"}}", v)
	}
	fmt.Printf("  prompt    : %s (%d chars rendered)\n", step.PromptFile, len(rendered))

	// Simulation: write placeholder so the next step can read it.
	if step.Writes != "" {
		placeholder := fmt.Sprintf("[simulation] output of step %q — replace with real Claude output.", step.ID)
		dest := filepath.Join(r.workDir, step.Writes+".md")
		if err := os.WriteFile(dest, []byte(placeholder), 0644); err != nil {
			return fmt.Errorf("write artifact: %w", err)
		}
		r.cache[step.Writes] = placeholder
		fmt.Printf("  → wrote   : %s → %s\n", step.Writes, dest)
	}
	return nil
}

// readContext loads context/{project_id}/{name}.md, falling back to context/global/{name}.md.
func (r *Runner) readContext(name string) (string, error) {
	projectPath := filepath.Join(contextRoot, r.cfg.ProjectID, name+".md")
	if data, err := os.ReadFile(projectPath); err == nil {
		return string(data), nil
	}

	globalPath := filepath.Join(contextRoot, "global", name+".md")
	if data, err := os.ReadFile(globalPath); err == nil {
		return string(data), nil
	}

	return "", fmt.Errorf("context %q not found at %s or %s", name, projectPath, globalPath)
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

func readPromptInclude(name string) string {
	data, err := os.ReadFile(filepath.Join(promptsRoot, name))
	if err != nil {
		return ""
	}
	return string(data)
}
