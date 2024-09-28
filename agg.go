package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rushyn/gator/internal/database"
)

func scrapeFeeds (s *state) error {

	const layout = "Mon, 02 Jan 2006 15:04:05 -0700"
	
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil{
		fmt.Println("Unable to get feed id")
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(),database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID:            nextFeed.ID,
	})

	if err != nil{
		fmt.Println("Unable to makr feed fetched")
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil{
		fmt.Println("Unable to fetch feed")
		return err
	}

	fmt.Printf("\n%s\n%s\n", feed.Channel.Title, feed.Channel.Description)

	for _, itme := range feed.Channel.Item{

		newUUID, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		tm, err := time.Parse(layout, itme.PubDate)
		
		if err != nil {
			return err
		}



		newPost := database.CreatePostParams{
			ID:          newUUID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       itme.Title,
			Url:         itme.Link,
			Desription:  itme.Description,
			PublishedAt: tm,
			FeedID:      nextFeed.ID,
		}

		post, err := s.db.CreatePost(context.Background(), newPost)
		if err != nil {
			return err
		}
		fmt.Println(post.Title)
		
	}

	return nil
}