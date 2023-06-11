package main

import (
	"fmt"
	"github.com/vladislavs-poznaks/rss-feed/internal/auth"
	"github.com/vladislavs-poznaks/rss-feed/internal/database"
	"net/http"
)

type handleAuth func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handle handleAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 401, fmt.Sprintf("Unauthorized: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("Not found: %v", err))
			return
		}

		handle(w, r, user)
	}
}
