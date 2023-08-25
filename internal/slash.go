package internal

import (
	"github.com/bwmarrin/discordgo"
)

func (w *WebhookSlashCommand) Info() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        w.Name,
		Description: w.Desc,
		Options:     []*discordgo.ApplicationCommandOption{},
	}

	for _, opt := range w.Arguments {
		cmd.Options = append(cmd.Options, &discordgo.ApplicationCommandOption{
			Type:        discordCommandOption[opt.Type],
			Name:        opt.Name,
			Description: opt.Desc,
			Required:    opt.Req,
		})
	}

	return cmd
}

// The actual function that gets called when a command is run
func (w *WebhookSlashCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	request(*w)
	return
}
