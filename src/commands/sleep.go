package commands

// To work properly and be wake-able, the device must not enter a "deep-sleep" mode.
//
// For example, you'll need to disable Hibernation in Windows : powercfg /h off
//
// For the SSH command to work properly, you'll need to connect at least once from the container 
// to your target device, to properly exchange needed keys/fingerprint

import (
	"fmt"
	"os/exec"

	"wakeywakey/database"
	"wakeywakey/utils"

	"github.com/bwmarrin/discordgo"
)

var SleepDevice = discordgo.ApplicationCommand{
	Name:        "sleep",
	Description: "Puts a registered device to sleep",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "alias",
			Description: "The Alias of the device to put to sleep.",
			Required:    true,
			Autocomplete: true,
		},
	},
}

func HandleSleepDevice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	alias := options[0].StringValue()

	var err error
	ipAddress, err := database.GetIpByAlias(i.Member.User.ID, alias)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
			Content: "Putting `" + alias + "` to sleep...",
		},
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction: " + err.Error())
		return
	}

	err = runSSHCommand(ipAddress)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				utils.EmbedError("Failed to put device to sleep", "Error for `" + alias + "`: " + err.Error()),
			},
		})
		return
	}

	fmt.Println("Put " + alias + " to sleep for user " + i.Member.User.Username + " (" + i.Member.User.ID + ")")

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			utils.EmbedSuccess("Device going to sleep", "`" + alias + "` is entering sleep mode."),
		},
	})
}

func runSSHCommand(ipAddress string) (error) {
	command := exec.Command("/usr/bin/ssh", ipAddress, "runDLL32.exe powrprof.dll,SetSuspendState 0,1,0")
	_, err := command.CombinedOutput()
	return err
}

func HandleSleepAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
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