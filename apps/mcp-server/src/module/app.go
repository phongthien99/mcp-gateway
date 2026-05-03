package module

import (
	"context"
	"fmt"
	"os"

	"mcp-gateway/src/api"
	"mcp-gateway/src/config"
	"mcp-gateway/src/prompts"
	"mcp-gateway/src/resources"
	"mcp-gateway/src/scope"
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
			mcp.AsHandler(tools.NewContextTools),

			// Resources — data Claude can read by URI
			mcp.AsHandler(resources.NewArtifactResources),

			// Prompts — structured prompt templates Claude can retrieve
			mcp.AsHandler(prompts.NewWorkflowPrompts),

			api.NewFileServer,
		),

		fx.Supply(cfg.MCP),
		fx.Supply(cfg),

		fx.Provide(func(lc fx.Lifecycle, params mcp.Params) *mcpserver.MCPServer {
			srv := mcpserver.NewMCPServer(cfg.MCP.Name, cfg.MCP.Version)
			for _, h := range params.Handlers {
				h.Register(srv)
			}

			switch cfg.MCP.Transport {
			case mcp.TransportSSE:
				sseServer := mcpserver.NewSSEServer(srv,
					mcpserver.WithSSEContextFunc(scope.SSEContextFunc),
					mcpserver.WithAppendQueryToMessageEndpoint(),
				)
				lc.Append(fx.Hook{
					OnStart: func(context.Context) error {
						go sseServer.Start(fmt.Sprintf(":%d", cfg.MCP.Port)) //nolint:errcheck
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return sseServer.Shutdown(ctx)
					},
				})
			case mcp.TransportStdio:
				stdio := mcpserver.NewStdioServer(srv)
				ctx, cancel := context.WithCancel(context.Background())
				lc.Append(fx.Hook{
					OnStart: func(context.Context) error {
						go func() {
							stdio.Listen(ctx, os.Stdin, os.Stdout) //nolint:errcheck
							os.Exit(0)
						}()
						return nil
					},
					OnStop: func(context.Context) error {
						cancel()
						return nil
					},
				})
			}

			return srv
		}),

		fx.Invoke(
			func(*mcpserver.MCPServer) {},
			api.RegisterFileServer,
		),
	}

	if cfg.MCP.Transport == mcp.TransportStdio {
		opts = append([]fx.Option{fx.NopLogger}, opts...)
	}

	return fx.New(opts...)
}
