package cmd

import (
	"math/rand"
	"time"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var verboseFlag bool

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tidbits",
	Short: "Run sample apps built for self education",
	Long: `
Collection of sample applications and code snippets to enable me to learn a few things.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verboseFlag {
			log.SetLevel(log.DebugLevel)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("exiting application")
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetLevel(log.InfoLevel)

	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "turn on debug messages")
}
