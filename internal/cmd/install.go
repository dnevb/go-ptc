package cmd

import (
	"fmt"

	"github.com/dan/plymouth-theme-creator/internal/install"
	"github.com/spf13/cobra"
)

var installSystemDir string

var installCmd = &cobra.Command{
	Use:   "install <dir>",
	Short: "Copy theme to system themes directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		if err := install.Install(dir, installSystemDir); err != nil {
			return err
		}
		fmt.Printf("installed %s → %s\n", dir, installSystemDir)
		return nil
	},
}

func init() {
	installCmd.Flags().StringVar(&installSystemDir, "system-dir", "/usr/share/plymouth/themes", "System themes directory")
	rootCmd.AddCommand(installCmd)
}
