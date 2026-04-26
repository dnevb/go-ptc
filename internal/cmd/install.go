package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installSystemDir string

var installCmd = &cobra.Command{
	Use:   "install <dir>",
	Short: "Copy theme to system themes directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		fmt.Printf("install: dir=%s system-dir=%s\n", dir, installSystemDir)
		return fmt.Errorf("not implemented")
	},
}

func init() {
	installCmd.Flags().StringVar(&installSystemDir, "system-dir", "/usr/share/plymouth/themes", "System themes directory")
	rootCmd.AddCommand(installCmd)
}
