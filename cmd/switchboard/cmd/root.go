package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-zen-chu/switchboard"
	"github.com/spf13/cobra"
)

// SwichboardRequirements is the all requirements for switchboard. Used for dependency injection.
type SwitchboardRequirements struct {
	Ctx           context.Context
	BlueskyClient switchboard.BlueskyClient
	XClient       switchboard.XClient
}

func NewRootCmd(switchboardReq *SwitchboardRequirements) *cobra.Command {
	const defaultVerbose = false
	var verbose bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "switchboard",
		Short: "A tool that connect between sns like a switchboard operator",
		Long:  `A tool that connect between sns like a switchboard operator`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
		},
	}
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose", "v",
		defaultVerbose,
		"verbose output (log level debug)",
	)
	rootCmd.AddCommand(NewBluesky2XCmd(switchboardReq.Ctx, switchboardReq.BlueskyClient, switchboardReq.XClient))
	return rootCmd
}
