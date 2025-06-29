package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, strings.NewReader(""))
	if err != nil{
		return nil, fmt.Errorf("Error generating request: %w", err)
	}
	req.Header.Set("User-Agent", "hbgator")
	res, err := client.Do(req)
	if err != nil{
		return nil, fmt.Errorf("Error requesting RSSFeed: %w", err)
	}
	defer res.Body.Close()
	rssData, err := io.ReadAll(res.Body)
	if err != nil{
		return nil, fmt.Errorf("Error reading RSSFeed data: %w", err)
	}
	feed := &RSSFeed{}
	if err = xml.Unmarshal(rssData, feed); err != nil{
		return nil, fmt.Errorf("Error unmarshaling RSSFeed data: %w", err)
	}
	feed.unescapeFeed()
	return feed, nil
}

func (f *RSSFeed) unescapeFeed(){
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)
	for _, item := range f.Channel.Item{
		item.unescapeItem()
	}
}

func (i *RSSItem) unescapeItem(){
	i.Title = html.UnescapeString(i.Title)
	i.Description = html.UnescapeString(i.Description)
}