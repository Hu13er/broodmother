package main

import (
	"log"

	"github.com/spf13/cobra"

	"gitlab.com/pirates1/broodmother"
	"gitlab.com/pirates1/broodmother/inspect"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Args:  cobra.MinimumNArgs(1),
	Run:   generate,
}

func generate(cmd *cobra.Command, args []string) {
	path := args[0]
	exec := broodmother.Executor{
		Generators: []broodmother.Generator{
			&inspect.Inspector{},
		},
	}
	if err := exec.ParseFile(path); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
