package commands

import (
	"wakeywakey/database"
	"wakeywakey/utils"

	"github.com/bwmarrin/discordgo"
)

var UnregisterWake = discordgo.ApplicationCommand{
	Name:        "unregister",
	Description: "Unregisters a device by its alias for Wake-on-LAN.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "alias",
			Description: "The alias of the device to unregister.",
			Required: true,
		},
	},
}

func HandleUnregister(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	alias := options[0].StringValue()

	err := database.RemoveWakeEntryByAlias(i.Member.User.ID, alias)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					utils.EmbedError("Unregistration Failed", "Failed to unregister device : "+err.Error()),
				},
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
			Embeds: []*discordgo.MessageEmbed{
				utils.EmbedSuccess("Unregistration Successful", "Successfully unregistered device."),
			},
		},
	})
}

func HandleUnregisterAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	entries, err := database.GetAllEntriesByUserId(i.Member.User.ID)
	if err != nil {
		return
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, entry := range entries {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  entry.Alias,
			Value: entry.Alias,
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}