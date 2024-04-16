package httpfunctions

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/DreamyMemories/blog-aggregator/internal/database"
)

func StartScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v gorountines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	// Ensures the first one runs
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Println("Error fetching feeds", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(wg, db, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done() // Tell WG to decrease count by one once done
	rssFeed, err := UrlToFeed(feed.Url)
	if err != nil {
		log.Println("Something went wrong in fetching feed", err.Error())
		return
	}
	_, err = db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Soemthing went wrong in updating the marked feed", err.Error())
		return
	}

	// Test
	for _, item := range rssFeed.Channel.Item {
		log.Println("Found Post", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
