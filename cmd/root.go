package cmd

import (
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use: "evermosonlinestore",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Println("running root command...")
	},
}

func Run() {
	if err := rootCommand.Execute(); err != nil {
		logger.Panic(err)
	}
}
