package internal

import "github.com/bwmarrin/discordgo"

type WebhookSlashCommand struct {
	Name          string                         `mapstructure:"name"`
	Desc          string                         `mapstructure:"description"`
	Resp          string                         `mapstructure:"response"`
	RespCode      int                            `mapstructure:"response_code"`
	URL           string                         `mapstructure:"url"`
	Method        string                         `mapstructure:"method"`
	Headers       []WebhookHeader                `mapstructure:"headers"`
	SubCmd        map[string]WebhookSlashCommand `mapstructure:"subcommands"`
	SubCmdGrp     map[string]WebhookSlashCommand `mapstructure:"subcommand_groups"`
	Arguments     []WebhookArgument              `mapstructure:"arguments"`
	Data          map[string]interface{}         `mapstructure:"data"`
	CalledOptions []*discordgo.ApplicationCommandInteractionDataOption
	CalledUser    *discordgo.Member
}

type WebhookHeader struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

type WebhookArgument struct {
	Name        string      `mapstructure:"name"`
	Desc        string      `mapstructure:"description"`
	Type        string      `mapstructure:"type"`
	Req         bool        `mapstructure:"required"`
	Default     interface{} `mapstructure:"default"`
	DiscordInfo bool        `mapstructure:"discord"`
}

var (
	discordCommandOption = map[string]discordgo.ApplicationCommandOptionType{
		"string": discordgo.ApplicationCommandOptionString,
		"int":    discordgo.ApplicationCommandOptionInteger,
		"float":  discordgo.ApplicationCommandOptionNumber,
		"bool":   discordgo.ApplicationCommandOptionBoolean,
	}
)
