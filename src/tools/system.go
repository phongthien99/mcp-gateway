package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type SystemTools struct{}

func NewSystemTools() *SystemTools {
	return &SystemTools{}
}

func (s *SystemTools) Register(srv *mcpserver.MCPServer) {
	srv.AddTool(mcp.NewTool("system_info",
		mcp.WithDescription("Get information about the host system (OS, arch, CPU count, hostname)"),
	), s.systemInfo)

	srv.AddTool(mcp.NewTool("get_env",
		mcp.WithDescription("Get the value of an environment variable"),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("Name of the environment variable"),
		),
	), s.getEnv)

	srv.AddTool(mcp.NewTool("list_env",
		mcp.WithDescription("List all environment variable names (values are hidden for security)"),
	), s.listEnv)

	srv.AddTool(mcp.NewTool("run_command",
		mcp.WithDescription("Run a shell command and return its output (stdout + stderr)"),
		mcp.WithString("command",
			mcp.Required(),
			mcp.Description("Command to execute"),
		),
		mcp.WithArray("args",
			mcp.Description("Arguments to pass to the command"),
		),
		mcp.WithString("cwd",
			mcp.Description("Working directory for the command"),
		),
	), s.runCommand)

	srv.AddTool(mcp.NewTool("current_time",
		mcp.WithDescription("Get the current server time"),
		mcp.WithString("format",
			mcp.Description("Time format: rfc3339 (default), unix, or human"),
		),
	), s.currentTime)
}

func (s *SystemTools) systemInfo(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hostname, _ := os.Hostname()
	info := fmt.Sprintf(
		"os: %s\narch: %s\ncpus: %d\nhostname: %s\ngo_version: %s",
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		hostname,
		runtime.Version(),
	)
	return mcp.NewToolResultText(info), nil
}

func (s *SystemTools) getEnv(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key := mcp.ParseArgument(req, "key", "").(string)
	if key == "" {
		return mcp.NewToolResultError("key is required"), nil
	}
	val := os.Getenv(key)
	if val == "" {
		return mcp.NewToolResultText(fmt.Sprintf("%s is not set", key)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s=%s", key, val)), nil
}

func (s *SystemTools) listEnv(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	envs := os.Environ()
	keys := make([]string, 0, len(envs))
	for _, e := range envs {
		parts := strings.SplitN(e, "=", 2)
		keys = append(keys, parts[0])
	}
	return mcp.NewToolResultText(strings.Join(keys, "\n")), nil
}

func (s *SystemTools) runCommand(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command := mcp.ParseArgument(req, "command", "").(string)
	if command == "" {
		return mcp.NewToolResultError("command is required"), nil
	}

	var args []string
	if params, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if rawArgs, ok := params["args"]; ok && rawArgs != nil {
			if argSlice, ok := rawArgs.([]interface{}); ok {
				for _, a := range argSlice {
					args = append(args, fmt.Sprintf("%v", a))
				}
			}
		}
	}

	cwd := mcp.ParseArgument(req, "cwd", "").(string)

	cmd := exec.Command(command, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("exit error: %v\noutput:\n%s", err, string(out))), nil
	}
	return mcp.NewToolResultText(string(out)), nil
}

func (s *SystemTools) currentTime(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	format := mcp.ParseArgument(req, "format", "rfc3339").(string)
	now := time.Now()
	switch format {
	case "unix":
		return mcp.NewToolResultText(fmt.Sprintf("%d", now.Unix())), nil
	case "human":
		return mcp.NewToolResultText(now.Format("Monday, 02 January 2006 15:04:05 MST")), nil
	default:
		return mcp.NewToolResultText(now.Format(time.RFC3339)), nil
	}
}
