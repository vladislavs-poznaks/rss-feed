package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/vladislavs-poznaks/rss-feed/internal/database"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create a feed follow: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(follow))
}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	follows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Could not get feed follows: %v", err))
		return
	}

	respondWithJson(w, 200, databaseFeedFollowsToFeedFollows(follows))
}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowIdStr := chi.URLParam(r, "feedFollowId")

	feedFollowID, err := uuid.Parse(feedFollowIdStr)

	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Could not parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Could not delete feed follows: %v", err))
		return
	}

	respondWithJson(w, 204, struct{}{})
}
