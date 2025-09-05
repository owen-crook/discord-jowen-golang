package config

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken                 string
	GuildID                  string
	UnregisterCommandsOnExit bool
	FirestoreProjectId       string
	FirestoreDatabaseId      string
}

func LoadConfig() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is not set")
	}

	guildID := os.Getenv("GUILD_ID")
	if guildID == "" {
		return nil, fmt.Errorf("GUILD_ID is not set")
	}

	FirestoreProjectId := os.Getenv("FIRESTORE_PROJECT_ID")
	if FirestoreProjectId == "" {
		return nil, fmt.Errorf("FIRESTORE_PROJECT_ID is not set")
	}

	FirestoreDatabaseId := os.Getenv("FIRESTORE_DATABASE_ID")
	if FirestoreDatabaseId == "" {
		return nil, fmt.Errorf("FIRESTORE_DATABASE_ID is not set")
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		return nil, fmt.Errorf("ENVIRONMENT is not set")
	}
	if environment == "LOCAL" && os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS required for local development")
	}

	unregisterCommandsOnExit := os.Getenv("UNREGISTER_COMMANDS_ON_EXIST") == "true"

	return &Config{
		BotToken:                 botToken,
		GuildID:                  guildID,
		UnregisterCommandsOnExit: unregisterCommandsOnExit,
		FirestoreProjectId:       FirestoreProjectId,
		FirestoreDatabaseId:      FirestoreDatabaseId,
	}, nil
}
