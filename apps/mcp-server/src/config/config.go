package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/gestgo/gest/package/extension/mcp"
)

type AppConfig struct {
	MCP  MCPConfig
	API  APIConfig
	Dirs DirsConfig
}

type MCPConfig = mcp.Config

type APIConfig struct {
	Port int
}

type DirsConfig struct {
	ProjectsRoot string // base dir for all project context + artifacts
	Prompts      string
	Workflows    string
	Runs         string
	Docs         string
	HugoContent  string
}

func Load() AppConfig {
	transport := mcp.TransportSSE
	if os.Getenv("MCP_TRANSPORT") == "stdio" {
		transport = mcp.TransportStdio
	}

	port := 8099
	if p := os.Getenv("MCP_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}

	name := envOr("MCP_NAME", "mcp-workbench")
	version := envOr("MCP_VERSION", "1.0.0")

	apiPort := 8110
	if p := os.Getenv("API_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			apiPort = v
		}
	}

	return AppConfig{
		MCP: MCPConfig{
			Name:      name,
			Version:   version,
			Transport: transport,
			Port:      port,
		},
		API: APIConfig{
			Port: apiPort,
		},
		Dirs: loadDirs(),
	}
}

// loadDirs builds directory paths.
// DOCS_ROOT (e.g. /hugo-src/content/docs) is the base for prompts and workflows.
// PROJECTS_ROOT (e.g. /hugo-src/content/docs/projects) is the single root for all
// project context and artifacts, organised as {project_id}/context/ and
// {project_id}/artifacts/. Falls back to {DOCS_ROOT}/projects when not set.
func loadDirs() DirsConfig {
	docsRoot     := os.Getenv("DOCS_ROOT")
	projectsRoot := os.Getenv("PROJECTS_ROOT")
	if projectsRoot == "" {
		if docsRoot != "" {
			projectsRoot = filepath.Join(docsRoot, "projects")
		} else {
			projectsRoot = filepath.Join("data", "projects")
		}
	}
	return DirsConfig{
		ProjectsRoot: projectsRoot,
		Prompts:      docsDir(docsRoot, "prompts"),
		Workflows:    workflowsDir(docsRoot),
		Runs:         "runs",
		Docs:         "docs",
		HugoContent:  hugoContentDir(docsRoot),
	}
}

func workflowsDir(docsRoot string) string {
	if docsRoot == "" {
		return filepath.Join("data", "workflows")
	}
	return filepath.Join(docsRoot, "workflows")
}

func docsDir(root, name string) string {
	if root == "" {
		// local dev: files live directly under data/{name}/
		return filepath.Join("data", name)
	}
	// Docker: volume is mounted at {DOCS_ROOT}/{name}/source
	return filepath.Join(root, name, "source")
}

func hugoContentDir(docsRoot string) string {
	if docsRoot == "" {
		return "hugo-content"
	}
	return filepath.Dir(docsRoot) // /hugo-src/content/docs → /hugo-src/content
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
