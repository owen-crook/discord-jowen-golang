package config

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken string
	GuildID  string
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

	return &Config{
		BotToken: botToken,
		GuildID:  guildID,
	}, nil
}
