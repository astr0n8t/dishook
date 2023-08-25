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

	//session :=
	login(token)

	commands := []WebhookSlashCommand{}

	configErr := config.UnmarshalKey("commands", &commands)
	if configErr != nil {
		log.Fatalf("Unable to read config: %v ", configErr)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+c to exit")
	<-stop

	return
}
