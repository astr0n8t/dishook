package internal

import (
	"log"
)

func (w *WebhookSlashCommand) request() {
	log.Println("Calling request with command name: ", w.Name)
	return
}
