package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/four-cryptics/crypto-news-notifier/pkg"
)

// mainCmd represents the base command when called without any subcommands
var mainCmd = &cobra.Command{
	Use:   "go run main.go",
	Short: "Crypto news scraper",
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// action currently empty due to subcommands being used
	},
}

var scrapeData pkg.ScrapeData
var url 	string

// scrapeCmd is a Cobra command for use with the scrape and other subcommands in CLI
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Command used for running news scraping scripts. The command will scrape all known sites if no flags are provided!",
    Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
        // Initialize options and run
        scrapeData = pkg.ScrapeData{scrapeData.ReadFile(), scrapeData.OutputFileName, pkg.Site{}}
        if url != "" {
            scrapeData.EnsurePageExists(url)
        } else {
			scrapeData.ScrapeAllSites()
		}
	},
}

func init() {
	mainCmd.AddCommand(scrapeCmd)
	// Here are defined flags and configuration settings for commands
    scrapeCmd.Flags().StringVarP(&scrapeData.OutputFileName, "output", "o", "./assets/news_site_articles.json", "set your json input/output file")
	scrapeCmd.Flags().StringVarP(&url, "add-url", "a", "", "add new news source site url")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	mainCmd.Execute()
}
