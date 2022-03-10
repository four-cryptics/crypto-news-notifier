package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var cfgFile string

type Config struct {
    OutDir              string          `mapstructure:"OUT_DIR"`
    CryptoBuzz          []string        `mapstructure:"CRYPTO_BUZZ"`
    WarPoliticsBuzz     []string        `mapstructure:"WAR_POLITICS_BUZZ"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go run main.go",
	Short: "Crypto-News-Notifier",
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// action currently empty due to subcommands being used
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