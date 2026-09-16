package commands

import (
	"fmt"
	"wakeywakey/database"
	"wakeywakey/utils"

	"github.com/bwmarrin/discordgo"
)

var RegisterWake = discordgo.ApplicationCommand{
	Name:        "register",
	Description: "Registers a device with an IP/MAC address via an alias for remote Wake-on-LAN and sleep",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "alias",
			Description: "The alias of the device to register.",
			Required: true,
		},
		{
			Type: discordgo.ApplicationCommandOptionBoolean,
			Name: "global",
			Description: "Whether the device can be used globally by anyone",
			Required: true,
		},
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "ip-address",
			Description: "With the following format : user@ip_adress",
			Required: true,
		},
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "mac-address",
			Description: "The MAC address of the device to register.",
			Required: true,
		},
	},
}

func HandleRegister(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	alias := options[0].StringValue()
	global := options[1].BoolValue()
	ipAddress := options[2].StringValue()
	macAddress := options[3].StringValue()

	err := database.AddWakeEntry(i.Member.User.ID, global, alias, ipAddress, macAddress)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					utils.EmbedError("Registration Failed", "Failed to register device : "+err.Error()),
				},
			},
		})
		return
	}

	fmt.Println("Registered for user " + i.Member.User.Username + " (" + i.Member.User.ID + ") : " + alias + " => " + ipAddress + " - " + macAddress, global)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
			Embeds: []*discordgo.MessageEmbed{
				utils.EmbedSuccess("Registration Successful", "Successfully registered device."),
			},
		},
	})
}
