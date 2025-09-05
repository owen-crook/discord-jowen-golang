package handlers

import (
	"context"
	"discord-jowen-golang/internal/db"
	"discord-jowen-golang/internal/logic"
	"fmt"
	"log"
	"sort"

	"github.com/bwmarrin/discordgo"
)

// Commands holds all the top-level slash commands for the bot
var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "hellowen",
		Description: "Says hello back!",
	},
	{
		Name:        "jowen-how",
		Description: "Provides the Github Link that powers jowen",
	},
	{
		Name:        "altf4",
		Description: "View the servers rage quitters",
	},
}

// SlashCommandHandlers maps command names to their handler functions
var SlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, db *db.FirestoreClient){
	"hellowen": func(s *discordgo.Session, i *discordgo.InteractionCreate, db *db.FirestoreClient) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "hellowen!",
			},
		})
	},
	"jowen-how": func(s *discordgo.Session, i *discordgo.InteractionCreate, db *db.FirestoreClient) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "[discord-jowen-golang](https://github.com/owen-crook/discord-jowen-golang)",
			},
		})
	},
	"altf4": func(s *discordgo.Session, i *discordgo.InteractionCreate, db *db.FirestoreClient) {
		leaderboard, err := logic.FetchRageQuitLeaderboard(context.Background(), db, i.GuildID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Failed to fetch leaderboard: %v", err),
				},
			})
			return
		}

		var userIDs []string
		for id := range leaderboard {
			userIDs = append(userIDs, id)
		}

		sort.Slice(userIDs, func(i, j int) bool {
			return leaderboard[userIDs[i]].TimesLeftServer > leaderboard[userIDs[j]].TimesLeftServer
		})

		medals := []string{"🥇", "🥈", "🥉"}
		formatted := "🏆 **Rage Quit Leaderboard** 🏆\n\n"
		if len(userIDs) == 0 {
			formatted += "No one has rage quit yet!"
		} else {
			topN := 3
			if len(userIDs) < topN {
				topN = len(userIDs)
			}
			formatted += "**Top Rage Quitters:**\n"
			for i := 0; i < topN; i++ {
				medal := medals[i]
				entry := leaderboard[userIDs[i]]
				formatted += fmt.Sprintf("%s <@%s> — %d times\n", medal, userIDs[i], entry.TimesLeftServer)
			}

			sort.Slice(userIDs, func(i, j int) bool {
				return leaderboard[userIDs[i]].LatestExit.After(leaderboard[userIDs[j]].LatestExit)
			})
			mostRecentID := userIDs[0]
			mostRecentEntry := leaderboard[mostRecentID]
			formatted += fmt.Sprintf("\n**Most Recent Rage Quitter:** <@%s> (%s)\n", mostRecentID, mostRecentEntry.LatestExit.Format("Jan 2, 2006 15:04 MST"))
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: formatted,
			},
		})
	},
}

// RegisterSlashCommands registers all slash command handlers
func RegisterSlashCommands(s *discordgo.Session, db *db.FirestoreClient) error {
	log.Println("Registering slash commands...")
	registered, err := s.ApplicationCommands(s.State.User.ID, s.State.Application.GuildID)
	if err != nil {
		return fmt.Errorf("cannot fetch registered commands: %v", err)
	}

	registeredCommandNames := make(map[string]bool)
	registeredCommandIDs := make(map[string]string)
	for _, cmd := range registered {
		registeredCommandNames[cmd.Name] = true
		registeredCommandIDs[cmd.Name] = cmd.ID
	}

	// Build a set of desired command names
	desiredCommandNames := make(map[string]bool)
	for _, v := range Commands {
		desiredCommandNames[v.Name] = true
	}

	// Unregister commands that exist but are not in our Commands list
	for name, id := range registeredCommandIDs {
		if !desiredCommandNames[name] {
			log.Printf("Unregistering command not in list: %s", name)
			err := s.ApplicationCommandDelete(s.State.User.ID, s.State.Application.GuildID, id)
			if err != nil {
				log.Printf("Failed to unregister command %s: %v", name, err)
			} else {
				log.Printf("Unregistered command: %s", name)
			}
		}
	}

	// Register only new commands
	for _, v := range Commands {
		if registeredCommandNames[v.Name] {
			log.Printf("Command already registered: %s", v.Name)
			continue
		}
		_, err := s.ApplicationCommandCreate(s.State.User.ID, s.State.Application.GuildID, v)
		if err != nil {
			return fmt.Errorf("cannot create slash command '%v': %v", v.Name, err)
		}
		log.Printf("Registered command: %s", v.Name)
	}
	log.Println("Slash commands registered!")
	return nil
}

// UnregisterSlashCommands removes all slash commands
func UnregisterSlashCommands(s *discordgo.Session) error {
	log.Println("Unregistering commands...")
	commands, err := s.ApplicationCommands(s.State.User.ID, s.State.Application.GuildID)
	if err != nil {
		log.Printf("Could not retrieve commands: %v", err)
		return err
	}

	var failed []string
	for _, command := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, s.State.Application.GuildID, command.ID)
		if err != nil {
			log.Printf("Could not delete command %s: %v", command.Name, err)
			failed = append(failed, command.Name)
		} else {
			log.Printf("Unregistered command: %s", command.Name)
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("failed to delete commands: %v", failed)
	}
	log.Println("All commands have been unregistered.")
	return nil
}
