package scope

import (
	"context"
	"net/http"
)

type contextKey struct{}

// FromContext returns the project scope stored in ctx, or "" if not set.
func FromContext(ctx context.Context) string {
	if v, ok := ctx.Value(contextKey{}).(string); ok {
		return v
	}
	return ""
}

// WithProject returns a new context with the given project scope.
func WithProject(ctx context.Context, projectID string) context.Context {
	return context.WithValue(ctx, contextKey{}, projectID)
}

// SSEContextFunc is a WithSSEContextFunc-compatible function that extracts
// the ?project= query parameter from the SSE connection request.
func SSEContextFunc(ctx context.Context, r *http.Request) context.Context {
	if p := r.URL.Query().Get("project"); p != "" {
		return WithProject(ctx, p)
	}
	return ctx
}
