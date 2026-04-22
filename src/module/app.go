package module

import (
	"mcp-gateway/src/config"
	"mcp-gateway/src/tools"

	"github.com/gestgo/gest/package/extension/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	cfg := config.Load()

	opts := []fx.Option{
		fx.Provide(
			mcp.AsHandler(tools.NewFileSystemTools),
			mcp.AsHandler(tools.NewSystemTools),
			mcp.AsHandler(tools.NewHTTPTools),
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
