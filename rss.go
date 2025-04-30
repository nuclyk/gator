package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/nuclyk/gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	scrapeFeed(s.db, feed)

	return nil
}

func parseTime(t string) (time.Time, error) {
	layouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
	}

	var err error
	var date time.Time
	for _, layout := range layouts {
		date, err = time.Parse(layout, t)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, err
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	markedFeed, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Could not mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), markedFeed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		publicationDate, err := parseTime(item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse time: %s: %v", item.PubDate, err)
			return
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: publicationDate,
			FeedID:      feed.ID,
		}

		db.CreatePost(context.Background(), params)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Set up http client, request and header
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error with request: %w\n", err)
	}
	req.Header.Set("User-Agent", "gator")

	// Send a request
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error with request: %w\n", err)
	}

	// Unmarshal the data
	var rssFeed RSSFeed
	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("Error around io.ReadAll: %w", err)
	}
	xml.Unmarshal(data, &rssFeed)

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i, item := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &rssFeed, nil
}
