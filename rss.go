package main

import (
	"context"
	"encoding/xml"
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
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	request.Header.Add("User-Agent", "gator")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return &RSSFeed{}, err
	}
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	rss := RSSFeed{}
	err = xml.Unmarshal(resBody, &rss)
	if err != nil {
		return &RSSFeed{}, err
	}

	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	//rss.Channel.Link = html.UnescapeString(rss.Channel.Link)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)

	for i := range rss.Channel.Item{
		rss.Channel.Item[i].Title = html.UnescapeString(rss.Channel.Item[i].Title)
		//rss.Channel.Item[i].Link = html.UnescapeString(rss.Channel.Item[i].Link )
		rss.Channel.Item[i].Description = html.UnescapeString(rss.Channel.Item[i].Description)
		//rss.Channel.Item[i].PubDate = html.UnescapeString(rss.Channel.Item[i].PubDate)
	}

	return &rss, nil
}