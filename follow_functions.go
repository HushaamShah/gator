package main

import (
	"context"
	"fmt"
	"time"

	"github.com/HushaamShah/gator/internal/database"
	"github.com/google/uuid"
)

func handleFollow(s *state, cmd command, userDetails database.User) error {

	feedId, err := s.queries.GetFeedIdByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	insertFollow, err := addFollowEntryForUser(s, feedId, userDetails)
	if err != nil {
		return err
	}
	fmt.Println(insertFollow.FeedName)
	fmt.Println(insertFollow.UserName)
	return nil
}

func handleUnfollow(s *state, cmd command, userDetails database.User) error {

	feedId, err := s.queries.GetFeedIdByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	args := database.UnfollowFeedParams{
		UserID: userDetails.ID,
		FeedID: feedId,
	}
	err1 := s.queries.UnfollowFeed(context.Background(), args)
	if err1 != nil {
		return err1
	}
	return nil
}

func addFollowEntryForUser(s *state, feedId uuid.UUID, userDetails database.User) (*database.CreateFeedFollowRow, error) {

	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feedId,
		UserID:    userDetails.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertFollow, err := s.queries.CreateFeedFollow(context.Background(), args)
	if err != nil {
		return nil, err
	}

	return &insertFollow, nil
}

func handleFollowing(s *state, cmd command, userDetails database.User) error {
	fmt.Println(s.dbconfig.User)
	userFeeds, err := s.queries.GetFeedFollowsForUser(context.Background(), userDetails.Name)
	if err != nil {
		return err
	}

	fmt.Println(userFeeds)
	return nil
}
