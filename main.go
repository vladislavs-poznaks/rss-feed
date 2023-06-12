package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vladislavs-poznaks/rss-feed/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Undefined port environment variable: PORT")
	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("Undefined port environment variable: DB_URL")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	db := database.New(conn)

	apiCfg := apiConfig{DB: db}

	go scrape(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		feed, e := getRssService("https://wagslane.dev/index.xml")

		if e != nil {
			respondWithError(w, 400, fmt.Sprintf("Err: %v", e))
			return
		}

		respondWithJson(w, 200, feed)
	})

	v1Router.Get("/ready", handleReady)
	v1Router.Get("/error", handleError)
	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.handleGetUser))

	v1Router.Get("/feeds", apiCfg.handleGetFeeds)
	v1Router.Post("/users/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))

	v1Router.Get("/user/feed-follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1Router.Post("/user/feed-follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Delete("/user/feed-follows/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))
	v1Router.Get("/user/feed-follows/posts", apiCfg.middlewareAuth(apiCfg.handleGetPosts))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Listening on port: %v...", port)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
