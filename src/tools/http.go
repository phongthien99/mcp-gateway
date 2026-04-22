package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type HTTPTools struct {
	client *http.Client
}

func NewHTTPTools() *HTTPTools {
	return &HTTPTools{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (h *HTTPTools) Register(s *mcpserver.MCPServer) {
	s.AddTool(mcp.NewTool("http_get",
		mcp.WithDescription("Perform an HTTP GET request and return the response body"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to fetch"),
		),
		mcp.WithString("headers",
			mcp.Description("Optional headers in key:value format, one per line"),
		),
	), h.httpGet)

	s.AddTool(mcp.NewTool("http_post",
		mcp.WithDescription("Perform an HTTP POST request with a body"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to post to"),
		),
		mcp.WithString("body",
			mcp.Description("Request body"),
		),
		mcp.WithString("content_type",
			mcp.Description("Content-Type header (default: application/json)"),
		),
		mcp.WithString("headers",
			mcp.Description("Additional headers in key:value format, one per line"),
		),
	), h.httpPost)
}

func parseHeaders(raw string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

func (h *HTTPTools) doRequest(req *http.Request, rawHeaders string) (*mcp.CallToolResult, error) {
	for k, v := range parseHeaders(rawHeaders) {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("request failed: %v", err)), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1 MB limit
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("cannot read response: %v", err)), nil
	}

	result := fmt.Sprintf("status: %s\nbody:\n%s", resp.Status, string(body))
	return mcp.NewToolResultText(result), nil
}

func (h *HTTPTools) httpGet(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := mcp.ParseArgument(req, "url", "").(string)
	if url == "" {
		return mcp.NewToolResultError("url is required"), nil
	}
	headers := mcp.ParseArgument(req, "headers", "").(string)

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid url: %v", err)), nil
	}
	return h.doRequest(httpReq, headers)
}

func (h *HTTPTools) httpPost(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := mcp.ParseArgument(req, "url", "").(string)
	if url == "" {
		return mcp.NewToolResultError("url is required"), nil
	}
	body := mcp.ParseArgument(req, "body", "").(string)
	contentType := mcp.ParseArgument(req, "content_type", "application/json").(string)
	headers := mcp.ParseArgument(req, "headers", "").(string)

	httpReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid url: %v", err)), nil
	}
	httpReq.Header.Set("Content-Type", contentType)
	return h.doRequest(httpReq, headers)
}
