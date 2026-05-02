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
	Artifacts   string
	Prompts     string
	Context     string
	Workflows   string
	Runs        string
	Docs        string
	HugoContent string
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
// If DOCS_ROOT is set (e.g. /hugo-src/content/docs), all four content dirs are
// derived as {DOCS_ROOT}/{name}/source so Go and Hugo share the same paths.
// Otherwise each dir defaults to a simple relative name.
func loadDirs() DirsConfig {
	docsRoot := os.Getenv("DOCS_ROOT")
	return DirsConfig{
		Artifacts:   docsDir(docsRoot, "artifacts"),
		Prompts:     docsDir(docsRoot, "prompts"),
		Context:     docsDir(docsRoot, "context"),
		Workflows:   docsDir(docsRoot, "workflows"),
		Runs:        "runs",
		Docs:        "docs",
		HugoContent: hugoContentDir(docsRoot),
	}
}

func docsDir(root, name string) string {
	if root == "" {
		return name
	}
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
