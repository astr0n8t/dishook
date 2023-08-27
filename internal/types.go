package internal

import "github.com/bwmarrin/discordgo"

// Overall struct to hold webhook command data
type WebhookSlashCommand struct {
	Name            string                         `mapstructure:"name"`
	Desc            string                         `mapstructure:"description"`
	Resp            string                         `mapstructure:"response"`
	RespCode        int                            `mapstructure:"response_code"`
	URL             string                         `mapstructure:"url"`
	Method          string                         `mapstructure:"method"`
	AuthHeaderName  string                         `mapstructure:"auth_header_name"`
	AuthHeaderValue string                         `mapstructure:"auth_header_value"`
	Headers         []WebhookHeader                `mapstructure:"headers"`
	SubCmd          map[string]WebhookSlashCommand `mapstructure:"subcommands"`
	SubCmdGrp       map[string]WebhookSlashCommand `mapstructure:"subcommand_groups"`
	Arguments       []WebhookArgument              `mapstructure:"arguments"`
	Data            map[string]interface{}         `mapstructure:"data"`
	// These fields are context specific
	CalledOptions []*discordgo.ApplicationCommandInteractionDataOption
	CalledUser    *discordgo.Member
}

// Holds a header for a webhook
type WebhookHeader struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

// Holds an argument for a webhook
type WebhookArgument struct {
	Name        string      `mapstructure:"name"`
	Desc        string      `mapstructure:"description"`
	Type        string      `mapstructure:"type"`
	Req         bool        `mapstructure:"required"`
	Default     interface{} `mapstructure:"default"`
	DiscordInfo bool        `mapstructure:"discord"`
}

// A map that shows all available argument types for a webhook command
var (
	discordCommandOption = map[string]discordgo.ApplicationCommandOptionType{
		"string": discordgo.ApplicationCommandOptionString,
		"int":    discordgo.ApplicationCommandOptionInteger,
		"float":  discordgo.ApplicationCommandOptionNumber,
		"bool":   discordgo.ApplicationCommandOptionBoolean,
	}
)
