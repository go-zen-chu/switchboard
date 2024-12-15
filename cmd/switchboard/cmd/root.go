package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// SwichboardRequirements is the all requirements for switchboard. Used for dependency injection.
type SwitchboardRequirements interface {
	Bluesky2XCmdRequirements
}

func NewRootCmd(req SwitchboardRequirements) *cobra.Command {
	const defaultVerbose = false
	var verbose bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "switchboard",
		Short: "A tool that connect between sns like a switchboard operator",
		Long:  `A tool that connect between sns like a switchboard operator`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLogger(verbose)
		},
	}
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose", "v",
		defaultVerbose,
		"verbose output (log level debug)",
	)
	rootCmd.AddCommand(NewBluesky2XCmd(req))
	return rootCmd
}

func setLogger(verbose bool) {
	// set logger to output to stderr because stdout is used for Generative AI response
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	if verbose {
		logHandler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	slog.SetDefault(slog.New(logHandler))
	slog.Debug("verbose debug log enabled")
}
