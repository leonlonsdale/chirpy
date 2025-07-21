package auth

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/util"
)

type contextKey string

const UserIDKey contextKey = "userID"

func (a *Auth) JWTProtect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := a.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, "Unauthorized: missing or invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := a.ValidateJWT(tokenString, a.cfg.Secret)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, "invalid token", nil)

			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}

	return userID, true
}
