package introspection

import "context"

var ctxKey = &contextKey{"oauth2-introspecion"}

type contextKey struct {
	key string
}

func WithValue(ctx context.Context, result *OAuth2IntrospectionResult) context.Context {
	return context.WithValue(ctx, ctxKey, result)
}

func FromContext(ctx context.Context) *OAuth2IntrospectionResult {
	val, _ := ctx.Value(ctxKey).(*OAuth2IntrospectionResult)
	return val
}
