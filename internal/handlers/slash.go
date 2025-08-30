package handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Commands holds all the top-level slash commands for the bot
var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "hello",
		Description: "Says hello back!",
	},
	{
		Name:        "jowen-how",
		Description: "Provides the Github Link that powers jowen",
	},
}

// SlashCommandHandlers maps command names to their handler functions
var SlashCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"hello": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello!",
			},
		})
	},
	"jowen-how": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "[discord-jowen-golang](https://github.com/owen-crook/discord-jowen-golang)",
			},
		})
	},
}

// RegisterSlashCommands registers all slash command handlers
func RegisterSlashCommands(s *discordgo.Session) error {
	log.Println("Registering slash commands...")
	for _, v := range Commands {
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
