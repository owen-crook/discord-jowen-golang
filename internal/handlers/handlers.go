package handlers

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"

	"discord-jowen-golang/internal/db"
	"discord-jowen-golang/internal/logic"
)

type Handlers struct {
	Firestore *db.FirestoreClient
}

func NewHandlers(firestore *db.FirestoreClient) *Handlers {
	return &Handlers{
		Firestore: firestore,
	}
}

// RegisterHandlers acts as the central point for registering all event and command handlers
func (h *Handlers) RegisterHandlers(s *discordgo.Session) {
	s.AddHandler(h.interactionCreate)
	log.Println("Registered InteractionCreate handler.")

	s.AddHandler(h.guildMemberRemove)
	log.Println("Registered GuildMemberRemove handler.")

	s.AddHandler(h.ready)
	log.Println("Registered Ready handler.")
}

func (h *Handlers) UnregisterSlashCommands(s *discordgo.Session) error {
	return UnregisterSlashCommands(s)
}

// ready is the handler for when the bot has successfully connected
func (h *Handlers) ready(s *discordgo.Session, r *discordgo.Ready) {
	// Register the slash commands by calling the function from slash_commands.go.
	log.Println("Bot is ready. Registering slash commands.")
	err := RegisterSlashCommands(s, h.Firestore)
	if err != nil {
		log.Fatalf("Failed to register slash commands: %v", err)
	}
}

// interactionCreate is the main entry point for all interactions
func (h *Handlers) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := SlashCommandHandlers[i.ApplicationCommandData().Name]; ok {
		handler(s, i, h.Firestore)
	}
}

// guildMemberRemove is the handler for when a user leaves the server
func (h *Handlers) guildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	log.Printf("User %s has left the server.", m.User.Username)
	err := logic.UpdateLeavingLeaderboard(context.Background(), h.Firestore, m.GuildID, m.User)
	if err != nil {
		log.Printf("Failed to update leaving leaderboard: %v", err)
	}
}
