package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type FileSystemTools struct {
	RootDir string
}

func NewFileSystemTools() *FileSystemTools {
	root, _ := os.Getwd()
	return &FileSystemTools{RootDir: root}
}

func (f *FileSystemTools) Register(s *mcpserver.MCPServer) {
	s.AddTool(mcp.NewTool("read_file",
		mcp.WithDescription("Read the contents of a file"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file (relative to root or absolute)"),
		),
	), f.readFile)

	s.AddTool(mcp.NewTool("write_file",
		mcp.WithDescription("Write content to a file, creating it if it doesn't exist"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the file"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Content to write"),
		),
		mcp.WithBoolean("overwrite",
			mcp.Description("Overwrite if file exists (default: false)"),
		),
	), f.writeFile)

	s.AddTool(mcp.NewTool("list_directory",
		mcp.WithDescription("List files and directories in a directory"),
		mcp.WithString("path",
			mcp.Description("Directory path (default: root directory)"),
		),
	), f.listDirectory)

	s.AddTool(mcp.NewTool("create_directory",
		mcp.WithDescription("Create a directory (including all parent directories)"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Directory path to create"),
		),
	), f.createDirectory)

	s.AddTool(mcp.NewTool("delete_path",
		mcp.WithDescription("Delete a file or directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to delete"),
		),
		mcp.WithBoolean("recursive",
			mcp.Description("Delete directories recursively (default: false)"),
		),
	), f.deletePath)

	s.AddTool(mcp.NewTool("move_path",
		mcp.WithDescription("Move or rename a file or directory"),
		mcp.WithString("src",
			mcp.Required(),
			mcp.Description("Source path"),
		),
		mcp.WithString("dst",
			mcp.Required(),
			mcp.Description("Destination path"),
		),
	), f.movePath)

	s.AddTool(mcp.NewTool("file_info",
		mcp.WithDescription("Get metadata about a file or directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to inspect"),
		),
	), f.fileInfo)
}

func (f *FileSystemTools) resolve(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(f.RootDir, path)
}

func (f *FileSystemTools) readFile(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	if path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}
	data, err := os.ReadFile(f.resolve(path))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot read file: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (f *FileSystemTools) writeFile(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	content := mcp.ParseArgument(req, "content", "").(string)
	overwrite := mcp.ParseArgument(req, "overwrite", false).(bool)

	if path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}

	full := f.resolve(path)
	if !overwrite {
		if _, err := os.Stat(full); err == nil {
			return mcp.NewToolResultError("file already exists; set overwrite=true to replace it"), nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot create parent directories: %v", err)), nil
	}
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot write file: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("written %d bytes to %s", len(content), path)), nil
}

func (f *FileSystemTools) listDirectory(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	dir := f.RootDir
	if path != "" {
		dir = f.resolve(path)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot read directory: %v", err)), nil
	}

	var sb strings.Builder
	for _, e := range entries {
		kind := "file"
		if e.IsDir() {
			kind = "dir "
		}
		sb.WriteString(fmt.Sprintf("[%s] %s\n", kind, e.Name()))
	}
	return mcp.NewToolResultText(sb.String()), nil
}

func (f *FileSystemTools) createDirectory(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	if path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}
	if err := os.MkdirAll(f.resolve(path), 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot create directory: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("directory created: %s", path)), nil
}

func (f *FileSystemTools) deletePath(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	recursive := mcp.ParseArgument(req, "recursive", false).(bool)
	if path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}

	full := f.resolve(path)
	var err error
	if recursive {
		err = os.RemoveAll(full)
	} else {
		err = os.Remove(full)
	}
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot delete: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("deleted: %s", path)), nil
}

func (f *FileSystemTools) movePath(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	src := mcp.ParseArgument(req, "src", "").(string)
	dst := mcp.ParseArgument(req, "dst", "").(string)
	if src == "" || dst == "" {
		return mcp.NewToolResultError("src and dst are required"), nil
	}
	if err := os.Rename(f.resolve(src), f.resolve(dst)); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot move: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("moved %s → %s", src, dst)), nil
}

func (f *FileSystemTools) fileInfo(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseArgument(req, "path", "").(string)
	if path == "" {
		return mcp.NewToolResultError("path is required"), nil
	}
	info, err := os.Stat(f.resolve(path))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot stat: %v", err)), nil
	}
	result := fmt.Sprintf("name: %s\nsize: %d bytes\nis_dir: %v\nmod_time: %s\nmode: %s",
		info.Name(),
		info.Size(),
		info.IsDir(),
		info.ModTime().Format(time.RFC3339),
		info.Mode().String(),
	)
	return mcp.NewToolResultText(result), nil
}
