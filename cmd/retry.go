/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/mshindle/tidbits/retry"
)

// retryCmd represents the retry command
var retryCmd = &cobra.Command{
	Use:   "retry",
	Short: "a simple circuit breaker example",
	Long: `
for this tidbit, we create two http servers where one server (port 8080) always returns a fail,
and the second server (port 8081) always returns some json successfully. We use a client which 
calls the bad host first and hopefully fails over to the good one.`,
	Run: func(cmd *cobra.Command, args []string) {
		hosts := []string{"localhost:8080", "localhost:8081"}
		mobyDick := &Book{Id: 1, Title: "Moby Dick", Author: "Herman Melville"}

		// set up the bad server
		go listenAndServe(hosts[0], http.StatusInternalServerError, mobyDick)
		// set up the good server
		go listenAndServe(hosts[1], http.StatusOK, mobyDick)

		// try and get a resource
		client := retry.NewClient(hosts...)

		resp, err := client.Get("/1")
		if err != nil {
			panic(err)
		}
		_ = resp.Write(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(retryCmd)
}

// Book is a simple object which we transform into a JSON string
type Book struct {
	Id     int
	Title  string
	Author string
}

// listenAndServe sets up a mini web server that serves a predetermined response.
func listenAndServe(addr string, statusCode int, book *Book) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			output, _ := json.MarshalIndent(book, "", "  ")
			_, _ = w.Write(output)
		}
		return
	})
	_ = http.ListenAndServe(addr, mux)
}
