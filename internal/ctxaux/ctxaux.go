package ctxaux

import (
	"context"

	middlewares "github.com/patrick-devel/shorturl/internal/middlwares"
)

func GetUserIDFromContext(ctx context.Context) string {
	contextData := ctx.Value(string(middlewares.ContextUserID))
	if contextData == nil {
		return ""
	}

	userID, ok := contextData.(string)
	if !ok {
		return ""
	}

	return userID
}
