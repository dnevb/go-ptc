package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview <dir>",
	Short: "Run plymouthd --test with theme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		fmt.Printf("preview: dir=%s\n", dir)
		return fmt.Errorf("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)
}
