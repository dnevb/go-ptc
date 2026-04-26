package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	createFps         int
	createRes         string
	createLoop        bool
	createTransition  string
	createOutputDir   string
)

var createCmd = &cobra.Command{
	Use:   "create <name> <media...>",
	Short: "Generate theme from media files",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		media := args[1:]
		fmt.Printf("create: name=%s media=%v fps=%d res=%s loop=%v transition=%s out=%s\n",
			name, media, createFps, createRes, createLoop, createTransition, createOutputDir)
		return fmt.Errorf("not implemented")
	},
}

func init() {
	createCmd.Flags().IntVar(&createFps, "fps", 30, "Frames per second for video/GIF extraction")
	createCmd.Flags().StringVar(&createRes, "res", "1920x1080", "Max resolution WxH")
	createCmd.Flags().BoolVar(&createLoop, "loop", false, "Loop animation")
	createCmd.Flags().StringVar(&createTransition, "transition", "none", "Transition type: fade|none")
	createCmd.Flags().StringVar(&createOutputDir, "output-dir", ".", "Output directory for theme")
	rootCmd.AddCommand(createCmd)
}
