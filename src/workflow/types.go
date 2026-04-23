package workflow

// Def is the top-level workflow definition parsed from YAML.
type Def struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

// Step describes one stage of the workflow.
type Step struct {
	ID          string `yaml:"id"`
	PromptFile  string `yaml:"prompt_file"`  // file name under prompts/ dir
	Description string `yaml:"description"`
	ExtraArgs   []Arg  `yaml:"extra_args"`   // args beyond project_id + feature_id
	Reads       []string `yaml:"reads"`      // artifact names to load as template vars
	Writes      string   `yaml:"writes"`     // artifact name this step produces
}

// Arg defines an extra argument for a step's prompt.
type Arg struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
}
