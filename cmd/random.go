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
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"gitlab.com/mshindle/tidbits/toy"
)

var n uint64

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random n",
	Short: "generate a random number between 0 and n-1",
	Long: `
random is a solution to an interview question that was posed to me. Given a function flip()
which returns true or false, generate a random number between 0 and n-1 with an equal
probability of each number occurring.`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(1)(cmd, args)
		if err != nil {
			return err
		}
		n, err = strconv.ParseUint(args[0], 10, 64)
		if err != nil || n < 1 {
			return fmt.Errorf("argument must be an integer greater than 0, received %s", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		r := toy.RandRange(uint(n))
		fmt.Printf("random number generated: %d", r)
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
}
