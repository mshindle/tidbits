package cmd

import (
	"math/rand"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verboseFlag bool

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tidbits",
	Short: "Run sample apps built for self education",
	Long: `
Collection of sample applications and code snippets to enable me to learn a few things.`,
	Version: "0.2.0",
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tidbits.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "turn on debug messages")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.WithError(err).Error("cannot find home directory")
			os.Exit(1)
		}

		// Search config in current & home directory with name ".tidbits" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".tidbits")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("file", viper.ConfigFileUsed()).Info("using config file")
	}
}