package contextkeys

import "context"

type contextKey uint

const (
	userID contextKey = iota
)

func ContextWithUserID(ctx context.Context, userIDStr string) context.Context {
	return context.WithValue(ctx, userID, userIDStr)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userIDStr, ok := ctx.Value(userID).(string)
	return userIDStr, ok
}
