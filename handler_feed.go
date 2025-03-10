package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Ultimace1314/blogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feedName> <feedURL>", cmd.Name)
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currentUser.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, currentUser)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Feeds:")
	for _, feed := range feeds {
		userName, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("* ID:            %s\n", feed.ID)
		fmt.Printf("* Created:       %v\n", feed.CreatedAt)
		fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
		fmt.Printf("* Name:          %s\n", feed.Name)
		fmt.Printf("* URL:           %s\n", feed.Url)
		fmt.Printf("* UserID:        %s\n", userName)
	}
	fmt.Println("=====================================")
	return nil
}

func handlerFollowFeed(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feedURL>", cmd.Name)
	}
	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return err
	}
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    currentUser.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Feed: %s\n", feedFollow.FeedName)
	fmt.Printf("Followed by: %s\n", s.cfg.CurrentUserName)
	return nil
}

func handlerGetFeedFollows(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err

	}
	follows, err := s.db.GetFeedFollows(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(follows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, ff := range follows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
