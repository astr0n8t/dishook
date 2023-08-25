package internal

import "github.com/bwmarrin/discordgo"

type WebhookSlashCommand struct {
	Name      string                 `mapstructure:"name"`
	Desc      string                 `mapstructure:"description"`
	Resp      string                 `mapstructure:"response"`
	URL       string                 `mapstructure:"url"`
	Headers   []WebhookHeader        `mapstructure:"headers"`
	Arguments []WebhookArgument      `mapstructure:"arguments"`
	Data      map[string]interface{} `mapstructure:"data"`
}

type WebhookHeader struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

type WebhookArgument struct {
	Name string `mapstructure:"name"`
	Desc string `mapstructure:"description"`
	Type string `mapstructure:"type"`
	Req  bool   `mapstructure:"required"`
}

var (
	discordCommandOption = map[string]discordgo.ApplicationCommandOptionType{
		"string": discordgo.ApplicationCommandOptionString,
		"int":    discordgo.ApplicationCommandOptionInteger,
		"float":  discordgo.ApplicationCommandOptionNumber,
		"bool":   discordgo.ApplicationCommandOptionBoolean,
	}
)
