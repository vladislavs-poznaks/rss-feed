package main

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/vladislavs-poznaks/rss-feed/internal/database"
	"log"
	"sync"
	"time"
)

func scrape(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf("Scraping on %v go routines every %s duration", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFetchFeeds(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Println("Error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFetchedFeed(context.Background(), feed.ID)

	if err != nil {
		log.Println("Error marking fetched feed: ", err)
		return
	}

	rssFeed, err := getRssService(feed.Url)

	if err != nil {
		log.Println("Error fetching feed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}

		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			log.Println("Could not parse publish date", err)
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			FeedID:      feed.ID,
			PublishedAt: publishedAt.UTC(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		})

		if err != nil {
			log.Println("Error creating post", err)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
