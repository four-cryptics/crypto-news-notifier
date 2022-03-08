package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// // Counts how many of the pre-defined buzzwords are included in the article title
// func (o Options) FindCryptoBuzz() {
//     for _, site := range o.Sites {

//     }
//     for word, count := range o.Buzzwords {
//         if strings.Contains(o.Title, word) {
//             a.Buzzwords[word] = count+1
//         }
//     }
// }

var evaluateCmd = &cobra.Command{
	Use:   "evaluate",
	Short: "Command used for running article relevance evaluation scripts",
    Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
        config, err := LoadConfig(".")
        if err != nil {
            fmt.Println("cannot load config:", err)
        }
		fmt.Println(config.OutDir)
	},
}

func init() {
    rootCmd.AddCommand(evaluateCmd)
	// Here are defined flags and configuration settings for commands
}