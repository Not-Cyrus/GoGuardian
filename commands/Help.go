package commands

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/bwmarrin/discordgo"
)

func (cmd *Commands) Help(s *discordgo.Session, message *discordgo.Message, ctx *Context) {

	if len(ctx.Fields) != 0 {
		for _, command := range cmd.Commands {
			if ctx.Fields[0] != command.Name {
				continue
			}
			var (
				count int = 0
				embed     = &discordgo.MessageEmbed{
					Title: fmt.Sprintf("Help %s", command.Name),
					Color: 0xc30101,
				}
			)
			for _, help := range command.AdvancedHelp {
				regexString = nameRegex.FindStringSubmatch(help)
				switch len(regexString) {
				case 0:
					commandName = command.Name
					commandValue = help
				default:
					commandName = regexString[0]
					commandValue = nameRegex.ReplaceAllString(help, "")
				}
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   commandName,
					Value:  commandValue,
					Inline: count%2 == 0,
				})
				count++
			}
			embed.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Sub Command Count: %d ", count)}
			s.ChannelMessageSendEmbed(message.ChannelID, embed)
		}
		return
	}

	var (
		count int = 0
		embed     = &discordgo.MessageEmbed{
			Title: "Commands",
			Color: 0xc30101,
		}
		names  = make([]string, len(cmd.Commands))
		sorted = make([]*command, len(cmd.Commands))
	)

	for index, cmds := range cmd.Commands {
		names[index] = cmds.Name
	}
	sort.Strings(names)

	// now sorted so iterate through it again
	for index, name := range names {
		for _, cmds := range cmd.Commands {
			if cmds.Name == name {
				sorted[index] = cmds
				break
			}
		}
		if sorted[index] == nil {
			s.ChannelMessageSend(message.ChannelID, "Failed to sort commands properly")
			return
		}
	}

	for _, cmds := range sorted {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   cmds.Name,
			Value:  cmds.Help,
			Inline: count%2 == 0,
		})
		count++
	}

	embed.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Command Count %d | Prefix \"%s\"", count, prefix)}
	s.ChannelMessageSendEmbed(message.ChannelID, embed)

}

var (
	commandName  string
	commandValue string
	regexString  []string
	nameRegex    = regexp.MustCompile(`\*\*[\w+]{1,30}\(.+\)\*\*`)
)
