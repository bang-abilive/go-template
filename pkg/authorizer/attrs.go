package authorizer

import "context"

// UserAttr holds the subject attributes used in ABAC policy evaluation.
// It is stored on the echo context and passed as the `sub` argument to Enforce.
type UserAttr struct {
	ID    int64
	Role  string // role slug, e.g. "admin"
	Level int    // role level, e.g. 100 for admin
}

type contextKey struct{}

// ContextKey is the key used to store UserAttr in a context.
var ContextKey = contextKey{}

// FromContext retrieves UserAttr from a context, returning a zero-value and false if not set.
func FromContext(ctx context.Context) (UserAttr, bool) {
	v, ok := ctx.Value(ContextKey).(UserAttr)
	return v, ok
}

// WithContext returns a new context carrying the given UserAttr.
func WithContext(ctx context.Context, attr UserAttr) context.Context {
	return context.WithValue(ctx, ContextKey, attr)
}
