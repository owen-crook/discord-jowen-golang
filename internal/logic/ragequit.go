package logic

import (
	"context"
	"fmt"
	"time"

	"discord-jowen-golang/internal/db"

	"github.com/bwmarrin/discordgo"
)

const (
	RageQuitLeaderboardCollection = "rage_quit_leaderboard"
)

type RageQuitLeaderboardEntry struct {
	Username        string    `json:"username" firestore:"username"`
	TimesLeftServer int       `json:"times_left_server" firestore:"times_left_server"`
	LatestExit      time.Time `json:"latest_exit" firestore:"latest_exit"`
}

func FetchRageQuitLeaderboard(ctx context.Context, f *db.FirestoreClient, guildId string) (map[string]RageQuitLeaderboardEntry, error) {
	// check if there is a leaderboard for the server (1 doc per server)
	exists, err := f.CheckDocumentExists(ctx, RageQuitLeaderboardCollection, guildId)
	if err != nil {
		return nil, fmt.Errorf("failed to check leaderboard existence: %w", err)
	}

	if !exists {
		return make(map[string]RageQuitLeaderboardEntry), nil
	}

	doc, err := f.FetchDocument(ctx, RageQuitLeaderboardCollection, guildId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leaderboard document: %w", err)
	}
	var leaderboard map[string]RageQuitLeaderboardEntry
	if err := doc.DataTo(&leaderboard); err != nil {
		return nil, fmt.Errorf("failed to convert leaderboard data: %w", err)
	}

	return leaderboard, nil
}

func UpdateLeavingLeaderboard(ctx context.Context, f *db.FirestoreClient, guildId string, user *discordgo.User) error {
	// check if there is a leaderboard for the server (1 doc per server)
	exists, err := f.CheckDocumentExists(ctx, RageQuitLeaderboardCollection, guildId)
	if err != nil {
		return fmt.Errorf("failed to check leaderboard existence: %w", err)
	}

	var leaderboard map[string]RageQuitLeaderboardEntry
	var entry RageQuitLeaderboardEntry

	if !exists {
		leaderboard = make(map[string]RageQuitLeaderboardEntry)
		entry = RageQuitLeaderboardEntry{
			Username:        user.Username,
			TimesLeftServer: 1,
			LatestExit:      time.Now().UTC(),
		}
		leaderboard[user.ID] = entry
	} else {
		doc, err := f.FetchDocument(ctx, RageQuitLeaderboardCollection, guildId)
		if err != nil {
			return fmt.Errorf("failed to fetch leaderboard document: %w", err)
		}
		if err := doc.DataTo(&leaderboard); err != nil {
			return fmt.Errorf("failed to convert leaderboard data: %w", err)
		}
		entry, ok := leaderboard[user.ID]
		if !ok {
			entry = RageQuitLeaderboardEntry{
				Username:        user.Username,
				TimesLeftServer: 1,
				LatestExit:      time.Now().UTC(),
			}
		} else {
			entry.TimesLeftServer++
			entry.LatestExit = time.Now().UTC()
			if entry.Username != user.Username {
				entry.Username = user.Username
			}
		}
		leaderboard[user.ID] = entry
	}

	if err := f.CreateOrOverwriteDocument(ctx, RageQuitLeaderboardCollection, guildId, leaderboard); err != nil {
		return fmt.Errorf("failed to update leaderboard document: %w", err)
	}
	return nil
}
