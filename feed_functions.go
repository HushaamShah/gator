package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/HushaamShah/gator/internal/database"
	"github.com/google/uuid"
)

func handleAggregate(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	xmlBody := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("user-agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err1 := xml.Unmarshal(body, &xmlBody)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}

	xmlBody.Channel.Title = html.UnescapeString(xmlBody.Channel.Title)
	xmlBody.Channel.Description = html.UnescapeString(xmlBody.Channel.Description)

	for i := range xmlBody.Channel.Item {
		xmlBody.Channel.Item[i].Title = html.UnescapeString(xmlBody.Channel.Item[i].Title)
		xmlBody.Channel.Item[i].Description = html.UnescapeString(xmlBody.Channel.Item[i].Description)
	}

	return &xmlBody, nil
}

func handleAddFeed(s *state, cmd command) error {
	userDetails, err := s.queries.GetUser(context.Background(), s.dbconfig.User)
	if err != nil {
		return err
	}

	// feed, _ := fetchFeed(context.Background(), cmd.args[1])

	args := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    userDetails.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err1 := s.queries.CreateFeed(context.Background(), args)
	if err1 != nil {
		return err1
	}
	return nil
}

func handleFeeds(s *state, cmd command) error {
	feedsDetails, err := s.queries.GetAllFeedsWithUserName(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(feedsDetails)
	return nil
}
