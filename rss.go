package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
