package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Hu13er/broodmother"
	"github.com/Hu13er/broodmother/httpgen"
	"github.com/Hu13er/broodmother/inspect"
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
			&httpgen.HttpGen{},
		},
	}
	if err := exec.ParseFile(path); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
