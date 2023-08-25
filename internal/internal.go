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

func Run() {
	config := config.Config()

	token := config.GetString("token")

	session := login(token)

	tmpCommands := []WebhookSlashCommand{}
	commands := map[string]WebhookSlashCommand{}

	configErr := config.UnmarshalKey("commands", &tmpCommands)
	if configErr != nil {
		log.Fatalf("Unable to read config: %v ", configErr)
	}

	for _, cmd := range tmpCommands {
		commands[cmd.Name] = cmd
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := commands[i.ApplicationCommandData().Name]; ok {
			cmd.Handler(s, i)
		}
	})

	createdCommands := []*discordgo.ApplicationCommand{}
	for _, cmd := range commands {
		createdCommand, err := session.ApplicationCommandCreate(session.State.User.ID, "", cmd.Info())
		if err != nil {
			log.Printf("Failed to register command: ", cmd.Name)
		}
		createdCommands = append(createdCommands, createdCommand)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+c to exit")
	<-stop

	for _, cmd := range createdCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Printf("Failed to remove command: %v with error %v ", cmd.Name, err)
		}
	}

	return
}
