package config

import (
	"os"
	"strconv"

	"github.com/gestgo/gest/package/extension/mcp"
)

type AppConfig struct {
	MCP mcp.Config
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
		name = "mcp-gateway"
	}

	version := os.Getenv("MCP_VERSION")
	if version == "" {
		version = "1.0.0"
	}

	return AppConfig{
		MCP: mcp.Config{
			Name:      name,
			Version:   version,
			Transport: transport,
			Port:      port,
		},
	}
}
