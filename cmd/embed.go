// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/spf13/cobra"

	"gitlab.com/mshindle/tidbits/embed"
)

// errorhCmd represents the gcd command
var embedCmd = &cobra.Command{
	Use:   "embed",
	Short: "understand how embedded fs works in conjunction with http server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetHandler(text.Default)
		embed.Execute(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(embedCmd)
}
