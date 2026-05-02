package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mcp-gateway/src/config"

	"go.uber.org/fx"
)

type FileServer struct {
	cfg   config.APIConfig
	roots map[string]string
}

type fileInfo struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size,omitempty"`
}

type fileReadResponse struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type fileWriteRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type fileWriteResponse struct {
	Path string `json:"path"`
	Size int    `json:"size"`
}

func NewFileServer(cfg config.AppConfig) *FileServer {
	return &FileServer{
		cfg: cfg.API,
		roots: map[string]string{
			"prompts":      cfg.Dirs.Prompts,
			"workflows":    cfg.Dirs.Workflows,
			"context":      cfg.Dirs.Context,
			"artifacts":    cfg.Dirs.Artifacts,
			"runs":         cfg.Dirs.Runs,
			"docs":         cfg.Dirs.Docs,
			"hugo-content": cfg.Dirs.HugoContent,
		},
	}
}

func RegisterFileServer(lc fx.Lifecycle, fs *FileServer) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/files", fs.handleFiles)
	mux.HandleFunc("/api/files/read", fs.handleRead)
	mux.HandleFunc("/api/files/write", fs.handleWrite)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", fs.cfg.Port),
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					fmt.Fprintf(os.Stderr, "file api server error: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func (fs *FileServer) handleFiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fs.listFiles(w, r)
	case http.MethodDelete:
		fs.deleteFile(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (fs *FileServer) listFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}
	if path == "." {
		fs.listEditableRoots(w)
		return
	}
	full, rel, err := fs.safeEditablePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entries, err := os.ReadDir(full)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot list files: %v", err), http.StatusInternalServerError)
		return
	}

	files := make([]fileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		childPath := entry.Name()
		if rel != "." {
			childPath = filepath.ToSlash(filepath.Join(rel, entry.Name()))
		}
		files = append(files, fileInfo{
			Path:  childPath,
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		})
	}
	writeJSON(w, files)
}

func (fs *FileServer) listEditableRoots(w http.ResponseWriter) {
	roots := make([]fileInfo, 0, len(fs.roots))
	for root := range fs.roots {
		roots = append(roots, fileInfo{
			Path:  root,
			Name:  root,
			IsDir: true,
		})
	}
	writeJSON(w, roots)
}

func (fs *FileServer) handleRead(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	full, rel, err := fs.safeEditablePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := os.ReadFile(full)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read file: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, fileReadResponse{Path: rel, Content: string(data)})
}

func (fs *FileServer) handleWrite(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req fileWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid json: %v", err), http.StatusBadRequest)
		return
	}
	full, rel, err := fs.safeEditablePath(req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		http.Error(w, fmt.Sprintf("cannot create parent directories: %v", err), http.StatusInternalServerError)
		return
	}
	if err := os.WriteFile(full, []byte(req.Content), 0644); err != nil {
		http.Error(w, fmt.Sprintf("cannot write file: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, fileWriteResponse{Path: rel, Size: len(req.Content)})
}

func (fs *FileServer) deleteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	full, _, err := fs.safeEditablePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := os.Remove(full); err != nil {
		http.Error(w, fmt.Sprintf("cannot delete file: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (fs *FileServer) safeEditablePath(path string) (string, string, error) {
	if path == "" {
		return "", "", fmt.Errorf("path is required")
	}
	clean := filepath.ToSlash(filepath.Clean(path))
	if clean == "/" {
		clean = "."
	}
	if strings.HasPrefix(clean, "/") || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", "", fmt.Errorf("invalid path %q", path)
	}
	if clean == "." {
		return "", "", fmt.Errorf("path is required")
	}

	root := strings.Split(clean, "/")[0]
	base, ok := fs.roots[root]
	if !ok {
		return "", "", fmt.Errorf("path root %q is not editable", root)
	}

	suffix := strings.TrimPrefix(clean, root)
	suffix = strings.TrimPrefix(suffix, "/")
	full := filepath.FromSlash(base)
	if suffix != "" {
		full = filepath.Join(full, filepath.FromSlash(suffix))
	}
	return full, clean, nil
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:1313" || origin == "http://127.0.0.1:1313" || origin == "http://localhost:1314" || origin == "http://127.0.0.1:1314" || origin == "http://localhost:1315" || origin == "http://127.0.0.1:1315" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}
