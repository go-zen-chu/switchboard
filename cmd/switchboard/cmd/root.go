/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd(cmdReq CommandRequirements) *cobra.Command {
	const defaultVerbose = false
	var verbose bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "switchboard",
		Short: "A tool that connect between sns",
		Long:  `A tool that connect between sns`,
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
	rootCmd.AddCommand(NewQueryCmd(cmdReq))
	return rootCmd
}
