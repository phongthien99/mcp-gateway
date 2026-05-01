package config

import (
	"os"
	"strconv"

	"github.com/gestgo/gest/package/extension/mcp"
)

type AppConfig struct {
	MCP MCPConfig
	API APIConfig
}

type MCPConfig = mcp.Config

type APIConfig struct {
	Port int
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

	name := os.Getenv("MCP_NAME")
	if name == "" {
		name = "mcp-workbench"
	}

	version := os.Getenv("MCP_VERSION")
	if version == "" {
		version = "1.0.0"
	}

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
	}
}
