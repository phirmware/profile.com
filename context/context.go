package context

import (
	"context"

	"profile.com/models"
)

type userCtx string

var u userCtx = "user"

// SetUserInContext sets the user in the request context object
func SetUserInContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, u, user)
}

// GetUserFromContext gets the user from the context
func GetUserFromContext(ctx context.Context) *models.User {
	if user := ctx.Value(u); user != nil {
		if user, t := user.(*models.User); t {
			return user
		}
	}
	return nil
}
