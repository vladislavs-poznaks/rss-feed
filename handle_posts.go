package main

import (
	"fmt"
	"github.com/vladislavs-poznaks/rss-feed/internal/database"
	"net/http"
)

func (apiCfg *apiConfig) handleGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetUserPosts(r.Context(), database.GetUserPostsParams{
		UserID: user.ID,
		Limit:  25,
	})

	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Could not get feed follows: %v", err))
		return
	}

	respondWithJson(w, 200, databasePostsToPosts(posts))
}
