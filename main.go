package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"discord-jowen-golang/internal/config"
	"discord-jowen-golang/internal/db"
	"discord-jowen-golang/internal/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	d, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}
	defer d.Close()

	// intitialze clients
	firestoreClient, err := db.NewFirestoreClient(context.Background(), cfg.FirestoreProjectId, cfg.FirestoreDatabaseId)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	d.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages // TODO: figure these out actualy

	h := handlers.NewHandlers(firestoreClient)
	h.RegisterHandlers(d)

	err = d.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if cfg.UnregisterCommandsOnExit {
		h.UnregisterSlashCommands(d)
	}
}
