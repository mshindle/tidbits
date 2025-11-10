package cmd

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"gitlab.com/mshindle/tidbits/foserver"
)

// foCmd executes the foserver fan-out example.
var timeout int64
var foCmd = &cobra.Command{
	Use:   "foserver",
	Short: "run a web server with endpoint that has multiple fan-out requests",
	Long:  ``,
	RunE:  runFanOutServer,
}

func init() {
	rootCmd.AddCommand(foCmd)
	foCmd.Flags().Int64VarP(&timeout, "timeout", "t", 2, "timeout (in seconds) for requests")
}

func runFanOutServer(cmd *cobra.Command, _ []string) error {
	log.WithField("timeout", timeout).Info("creating http client")
	var client = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return foserver.Execute(cmd.Context(), client)
}
