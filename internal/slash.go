package internal

import (
	"github.com/bwmarrin/discordgo"
)

// Returns the information needed for discord to register the commands
func (w *WebhookSlashCommand) Info() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        w.Name,
		Description: w.Desc,
		Options:     []*discordgo.ApplicationCommandOption{},
	}

	if len(w.SubCmd) == 0 && len(w.SubCmdGrp) == 0 {
		cmd.Options = w.optionsInfo()
	} else {
		cmd.Options = append(w.subCmdInfo(), w.subCmdGrpInfo()...)
	}

	return cmd
}

func (w *WebhookSlashCommand) subCmdGrpInfo() []*discordgo.ApplicationCommandOption {
	subCmdGrp := []*discordgo.ApplicationCommandOption{}
	for name, cmd := range w.SubCmdGrp {
		subCmdGrp = append(subCmdGrp, &discordgo.ApplicationCommandOption{
			Name:        name,
			Description: cmd.Desc,
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Options:     cmd.subCmdInfo(),
		})
	}

	return subCmdGrp
}

func (w *WebhookSlashCommand) subCmdInfo() []*discordgo.ApplicationCommandOption {
	subCmd := []*discordgo.ApplicationCommandOption{}
	for name, cmd := range w.SubCmd {
		subCmd = append(subCmd, &discordgo.ApplicationCommandOption{
			Name:        name,
			Description: cmd.Desc,
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options:     cmd.optionsInfo(),
		})
	}

	return subCmd
}

func (w *WebhookSlashCommand) optionsInfo() []*discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{}
	for _, opt := range w.Arguments {
		options = append(options, &discordgo.ApplicationCommandOption{
			Type:        discordCommandOption[opt.Type],
			Name:        opt.Name,
			Description: opt.Desc,
			Required:    opt.Req,
		})
	}
	return options
}

// The actual function that gets called when a command is run
func (w *WebhookSlashCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if we have initialized our options yet
	if w.CalledOptions == nil {
		// If we haven't get the options from the interaction
		w.CalledOptions = i.ApplicationCommandData().Options
	}
	// Check if this is a subcommand call
	if len(w.SubCmdGrp) != 0 || len(w.SubCmd) != 0 {
		// If it is, get the name
		name := w.CalledOptions[0].Name
		// Set the command to the current parent command
		cmd := w
		// Check if its a subgroup command
		subCmdGrp, isSubCmdGrp := cmd.SubCmdGrp[name]
		if isSubCmdGrp {
			// Get the suboptions for the subcommand
			w.CalledOptions = w.CalledOptions[0].Options
			// Set the name to the inner subcommand
			name = w.CalledOptions[0].Name
			// Set the current cmd to the sub command group
			cmd = &subCmdGrp
		}
		// Check if the given name is a subcommand of the current
		// command
		subCmd, isSubCmd := cmd.SubCmd[name]
		if isSubCmd {
			// Get the suboptions for the subcommand
			subCmd.CalledOptions = w.CalledOptions[0].Options
			// Call the subcommand handler now
			subCmd.Handler(s, i)
		}
		// Otherwise this is the correct command for the handler
	} else {

		// Process our webhook request
		w.request()

		// Return the given response for this command
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: w.Resp,
			},
		})
	}
	// Reset our called options
	w.CalledOptions = nil
	return
}
