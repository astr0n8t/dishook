package internal

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/fsnotify/fsnotify"

	"github.com/astr0n8t/dishook/config"
)

// Global variables for internal
var (
	loadedSession  *discordgo.Session
	loadedCommands []*discordgo.ApplicationCommand
	loadedGuildID  string
)

// Logs in to Discord and returns a new session
// If this function fails, it exits with a fatal error
func login(token string) *discordgo.Session {
	// Try to create a new session
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	// Add a handler to log when a login is completed
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open a new webhook connection
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Return the discord session
	return s
}

// Gets a map of commands out of the config
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

// Creates the commands in discord
func createCommands(session *discordgo.Session, commands map[string]WebhookSlashCommand, guildID string) []*discordgo.ApplicationCommand {
	// Create an array of commands we created
	createdCommands := []*discordgo.ApplicationCommand{}
	// Add all of our commands
	for _, cmd := range commands {
		// Attempt to create a command
		createdCommand, err := session.ApplicationCommandCreate(session.State.User.ID, guildID, cmd.Info())
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

// Deletes the commands from discord
func deleteCommands(session *discordgo.Session, createdCommands []*discordgo.ApplicationCommand, guildID string) {
	// Try to remove our commands
	for _, cmd := range createdCommands {
		// Try to remove a command
		err := session.ApplicationCommandDelete(session.State.User.ID, guildID, cmd.ID)
		// Don't panic on failure but log it
		if err != nil {
			log.Printf("Failed to remove command: %v with error %v ", cmd.Name, err)
		} else {
			log.Printf("Removed command: %v", cmd.Name)
		}
	}
}

// Loads a new session and creates all the commands
func Load() (*discordgo.Session, []*discordgo.ApplicationCommand, string) {
	// Get our config object
	config := config.Config()

	// Read our token from the config
	token := config.GetString("token")
	// Read in the guild id to register the commands in
	guildID := config.GetString("guild_id")

	// Login to discord and get our session
	session := login(token)
	log.Printf("Loaded guild_id: %v ", guildID)

	// Get our commands map
	commands := getCommandsFromConfig(config)

	// Add a handler that maps commands to their handler functions
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := commands[i.ApplicationCommandData().Name]; ok {
			cmd.Handler(s, i)
		}
	})

	// Try to create our commands
	return session, createCommands(session, commands, guildID), guildID
}

// Kills the current sessioon and tries to start a new one
func Reload(e fsnotify.Event) {
	// Log that we're reloading the config
	log.Printf("Config file changed: %v Reloading config", e.Name)
	deleteCommands(loadedSession, loadedCommands, loadedGuildID)
	// Close our websocket with discord
	loadedSession.Close()
	// Reset our commands
	loadedSession, loadedCommands = nil, nil
	// Load the config again
	loadedSession, loadedCommands, loadedGuildID = Load()
	log.Printf("Reloaded config: %v", e.Name)
}

// Runs dishook
func Run() {

	// Make sure we can load config
	config := config.Config()
	log.Printf("Loaded config file %v", config.ConfigFileUsed())

	// Add config hot reloading
	config.OnConfigChange(Reload)
	config.WatchConfig()
	log.Printf("Watching config for changes")

	// Login and start
	loadedSession, loadedCommands, loadedGuildID = Load()

	// Don't exit until we receive stop from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+c to exit")
	<-stop

	// Try to delete out created commands
	deleteCommands(loadedSession, loadedCommands, loadedGuildID)
	// Close our websocket with discord
	loadedSession.Close()
}
