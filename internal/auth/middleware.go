package auth

import (
	"context"
	"github.com/muhammadnajie/kp/internal/pkg/jwt"
	"github.com/muhammadnajie/kp/internal/users"
	"net/http"
	"strconv"
	"strings"
)

func Middleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			if !strings.Contains(authorizationHeader, "Bearer") {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
			username, err := jwt.ParseToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			user := users.User{Username: username}
			id, err := users.GetUserIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user.ID = strconv.Itoa(id)
			ctx := context.WithValue(r.Context(), "user", &user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ExtractUserContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value("user").(*users.User)
	return raw
}
