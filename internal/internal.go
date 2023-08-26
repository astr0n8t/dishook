package internal

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	"github.com/astr0n8t/dishook/config"
)

func init() {
	return
}

func login(token string) *discordgo.Session {
	s, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer s.Close()
	return s
}

func getCommandsFromConfig(config config.Provider) map[string]WebhookSlashCommand {
	// Store the commands in a map
	commands := map[string]WebhookSlashCommand{}

	// Try to read in the commands from the config
	configErr := config.UnmarshalKey("commands", &commands)
	if configErr != nil {
		log.Fatalf("Unable to read config: %v ", configErr)
	}

	// Add the names to the commands
	for name, cmd := range commands {
		cmd.Name = name
		// For sub command groups, it needs to add names to the
		// command group itself as well as its sub commands
		for subCmdGrpName, subCmdGrp := range cmd.SubCmdGrp {
			subCmdGrp.Name = subCmdGrpName
			// Add the name to the sub commands for this group
			for subCmdName, subCmd := range subCmdGrp.SubCmd {
				subCmd.Name = subCmdName
				subCmdGrp.SubCmd[subCmdName] = subCmd
			}
			cmd.SubCmdGrp[subCmdGrpName] = subCmdGrp
		}
		// For sub commands, it just needs to add names to the
		// sub commands themselves
		for subCmdName, subCmd := range cmd.SubCmd {
			subCmd.Name = subCmdName
			cmd.SubCmd[subCmdName] = subCmd
		}
		commands[name] = cmd
	}
	return commands
}

func createCommands(session *discordgo.Session, commands map[string]WebhookSlashCommand) []*discordgo.ApplicationCommand {
	// Create an array of commands we created
	createdCommands := []*discordgo.ApplicationCommand{}
	// Add all of our commands
	for _, cmd := range commands {
		// Attempt to create a command
		createdCommand, err := session.ApplicationCommandCreate(session.State.User.ID, "", cmd.Info())
		if err != nil {
			log.Printf("Failed to register command: %v", cmd.Name)
		} else {
			log.Printf("Registered command: %v", cmd.Name)
		}
		// Store that command in our array
		createdCommands = append(createdCommands, createdCommand)
	}
	return createdCommands
}

func deleteCommands(session *discordgo.Session, createdCommands []*discordgo.ApplicationCommand) {
	// Try to remove our commands
	for _, cmd := range createdCommands {
		// Try to remove a command
		err := session.ApplicationCommandDelete(session.State.User.ID, "", cmd.ID)
		// Don't panic on failure but log it
		if err != nil {
			log.Printf("Failed to remove command: %v with error %v ", cmd.Name, err)
		}
	}

}

func Run() {
	// Get our config object
	config := config.Config()

	// Read our token from the config
	token := config.GetString("token")

	// Login to discord and get our session
	session := login(token)

	// Get our commands map
	commands := getCommandsFromConfig(config)

	// Add a handler that maps commands to their handler functions
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := commands[i.ApplicationCommandData().Name]; ok {
			cmd.Handler(s, i)
		}
	})

	// Try to create our commands
	createdCommands := createCommands(session, commands)

	// Don't exit until we receive stop from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+c to exit")
	<-stop

	// Try to delete out created commands
	deleteCommands(session, createdCommands)

	return
}
