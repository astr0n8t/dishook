package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/astr0n8t/dishook/internal"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dishook",
	Short: "Converts slash commands to webhooks",
	Long: `Dishook allows you to take a definition through templating and YAML
and create very simple Discord commands which correlate to web requests`,
	Run: func(cmd *cobra.Command, args []string) { 
		internal.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
