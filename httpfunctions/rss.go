package httpfunctions

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func UrlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{Timeout: 10 * time.Second}

	response, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}
	return rssFeed, nil
}
