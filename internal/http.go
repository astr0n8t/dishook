package internal

import (
	"fmt"
)

func request(command WebhookSlashCommand) {
	fmt.Println("Calling request with command name: ", command.Name)
	return
}
