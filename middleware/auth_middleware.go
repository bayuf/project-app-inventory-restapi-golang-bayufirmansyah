package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/service"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	Service *service.AuthService
	Logger  *zap.Logger
}

type contextKey string

const authUserKey contextKey = "authUser"

type AuthUser struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
	Role   string
}

func NewAuthMiddleware(service *service.AuthService, log *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		Service: service,
		Logger:  log,
	}
}

func GetAuthUser(r *http.Request) (*AuthUser, bool) {
	user, ok := r.Context().Value(authUserKey).(*AuthUser)
	return user, ok
}

func (m *AuthMiddleware) SessionAuthMiddleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// get session from Authorization
			authHeader := r.Header.Get("Authorization")

			// validae
			if authHeader == "" {
				utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", "token is empty please login")
				return
			}

			// validate token type
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", "token type invalid")
				return
			}

			// parse string to uuid
			sessionID, err := uuid.Parse(parts[1])
			if err != nil {
				utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", "token invalid")
				return
			}

			// validate token
			sess, err := m.Service.ValidateSession(ctx, sessionID)
			if err != nil {
				utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", "token invalid or inactive")
				return
			}

			authUser := &AuthUser{
				ID:     sess.ID,
				UserID: sess.UserID,
				Name:   sess.Username,
				Role:   sess.RoleName,
			}

			ctxv := context.WithValue(r.Context(), authUserKey, authUser)

			h.ServeHTTP(w, r.WithContext(ctxv))
		})
	}
}

func (m *AuthMiddleware) RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetAuthUser(r)
			if !ok {
				utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}

			if _, exists := allowed[user.Role]; !exists {
				utils.ResponseFailed(w, http.StatusForbidden, "forbidden", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
