/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iseeall [directory]",
	Short: "A fast, robust CLI tool to visualize directory structures as a tree.",
	Long: `iseeall is a production-grade directory tree visualizer built for speed and reliability. 

Unlike standard tree commands, iseeall is designed with the real-world filesystem in mind:
  • Gracefully handles permission errors without crashing (warnings go to stderr).
  • Hidden files are hidden by default to reduce clutter (use -a to show them).
  • Prevents runaway recursion in massive directories with a configurable max depth (-L).
  • Follows Unix philosophy: clean tree output to stdout, errors to stderr.

Examples:
  iseeall                  # Visualize the current directory
  iseeall /var/log -L 2    # Visualize /var/log up to 2 levels deep
  iseeall -a --ignore-errors ./my-project # Show all files, silently skip unreadable dirs`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// We will implement the core logic here next!
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.iseeall.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


