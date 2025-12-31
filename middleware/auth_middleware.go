package middleware

import (
	"net/http"
	"strconv"
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

func NewAuthMiddleware(service *service.AuthService, log *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		Service: service,
		Logger:  log,
	}
}

func (m *AuthMiddleware) SessionAuthMiddleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get session from Authorization
			authHeader := r.Header.Get("Authorization")

			// validae
			if authHeader == "" {
				utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", "token is empty please login")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", "token type invalid")
				return
			}

			sessionID, err := uuid.Parse(parts[1])
			if err != nil {
				utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", "token invalid")
				return
			}

			sess, err := m.Service.ValidateSession(sessionID)
			if err != nil {
				utils.ResponseFailed(w, http.StatusUnauthorized, "unauthorized", "token invalid or inactive")
				return
			}

			// inject to header
			r.Header.Set("X-User-ID", sess.UserID.String())
			r.Header.Set("X-Role-ID", strconv.Itoa(sess.RoleId))

			h.ServeHTTP(w, r)
		})
	}
}
