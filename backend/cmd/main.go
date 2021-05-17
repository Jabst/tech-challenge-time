package main

import (
	"log"
	"pento/code-challenge/cmd/api"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "users [SERVICE]"}
	rootCmd.AddCommand(api.Command())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute %s", err)
	}
}
