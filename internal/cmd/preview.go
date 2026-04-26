package cmd

import (
	"github.com/dan/plymouth-theme-creator/internal/preview"
	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview <dir>",
	Short: "Run plymouthd --test with theme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		return preview.Preview(dir)
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)
}
