package api

import (
	api "pento/code-challenge/application"

	"github.com/spf13/cobra"
)

// Command creates cobra command.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tracker",
		Short: "Start Tracker API",
		RunE:  Run(),
	}

	return cmd
}

func Run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		api.SetupAPI()

		return nil
	}
}
