package middleware

import (
	"net/http"

	"profile.com/context"
	"profile.com/models"
)

// RequireUserMiddleWare defines the shape of the middleware struct
type RequireUserMiddleWare struct {
	models.UserService
}

// UserMiddleWare checks for a logged in user
type UserMiddleWare struct {
	models.UserService
}

// NewRequireUserMiddleWare returns the middleware struct
func NewRequireUserMiddleWare(us models.UserService) *RequireUserMiddleWare {
	return &RequireUserMiddleWare{
		UserService: us,
	}
}

// NewUserMiddleWare returns the user middleware struct
func NewUserMiddleWare(us models.UserService) *UserMiddleWare {
	return &UserMiddleWare{
		UserService: us,
	}
}

// ApplyFn is a middleware function
func (mw *RequireUserMiddleWare) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.GetUserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "login", http.StatusFound)
			return
		}
		next(w, r)
	})
}

// ApplyFn is a middleware function
func (mw *UserMiddleWare) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}

		ctx := context.SetUserInContext(r.Context(), user)
		r = r.WithContext(ctx)

		next(w, r)
	})
}

// Apply is a middleware function
func (mw *UserMiddleWare) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
