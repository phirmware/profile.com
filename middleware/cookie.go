package middleware

import (
	"net/http"

	"profile.com/context"
	"profile.com/models"
)

// MiddleWare defines the shape of the middleware struct
type MiddleWare struct {
	models.UserService
}

// NewMiddleWare returns the middleware struct
func NewMiddleWare(us models.UserService) *MiddleWare {
	return &MiddleWare{
		UserService: us,
	}
}

// ApplyFn is a middleware function
func (mw *MiddleWare) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.SetUserInContext(r.Context(), user)
		r = r.WithContext(ctx)

		next(w, r)
	})
}
