package module

import (
	"mcp-gateway/src/config"
	"mcp-gateway/src/prompts"
	"mcp-gateway/src/resources"
	"mcp-gateway/src/tools"

	"github.com/gestgo/gest/package/extension/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	cfg := config.Load()

	opts := []fx.Option{
		fx.Provide(
			// Tools — actions Claude can invoke
			mcp.AsHandler(tools.NewArtifactTools),

			// Resources — data Claude can read by URI
			mcp.AsHandler(resources.NewArtifactResources),

			// Prompts — structured prompt templates Claude can retrieve
			mcp.AsHandler(prompts.NewWorkflowPrompts),
		),

		fx.Supply(cfg.MCP),

		mcp.Module(),

		fx.Invoke(func(*mcpserver.MCPServer) {}),
	}

	// Stdio transport writes to stdout — suppress all Fx logs to avoid protocol corruption.
	if cfg.MCP.Transport == mcp.TransportStdio {
		opts = append([]fx.Option{fx.NopLogger}, opts...)
	}

	return fx.New(opts...)
}
