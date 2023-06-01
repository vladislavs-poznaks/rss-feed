package main

import "net/http"

func handleReady(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, struct{}{})
}
