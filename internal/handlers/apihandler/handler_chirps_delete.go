package apihandler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/leonlonsdale/chirpy/internal/auth"
	"github.com/leonlonsdale/chirpy/internal/database"
	"github.com/leonlonsdale/chirpy/internal/util"
)

func DeleteChirpByID(db database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			util.RespondWithError(w, http.StatusUnauthorized, "invalid jwt", nil)
			return
		}

		chirpID, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, "invalid chirp id format", err)
			return
		}

		chirp, err := db.GetChirpByID(r.Context(), chirpID)
		if err != nil {
			util.RespondWithError(w, http.StatusNotFound, "could not find chirp to delete", err)
			return
		}

		if userID != chirp.UserID {
			util.RespondWithError(w, http.StatusForbidden, "not author of chirp", nil)
			return
		}

		if err := db.DeleteChirp(r.Context(), database.DeleteChirpParams{
			ID:     chirpID,
			UserID: userID,
		}); err != nil {
			util.RespondWithError(w, http.StatusNotFound, "error deleting chirp", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
