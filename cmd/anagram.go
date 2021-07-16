package cmd

import (
	"fmt"
	"reflect"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

// anagramCmd represents the anagram command
var anagramCmd = &cobra.Command{
	Use:   "anagram",
	Short: "find the anagrams of a string",
	Long: `
Given two strings s and p, return an array of all the start indices of p's anagrams in s.
You may return the answer in any order.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ans := findAnagrams(args[0], args[1])
		fmt.Println(ans)
	},
}

func init() {
	rootCmd.AddCommand(anagramCmd)
}

func findAnagrams(s string, p string) []int {
	output := make([]int, 0, 5)
	pLen := len(p)

	if len(s) < pLen {
		return output
	}

	var check, window map[rune]int
	check = make(map[rune]int)
	window = make(map[rune]int)

	// fill out the map for p
	for _, ch := range p {
		check[ch] = check[ch] + 1
	}

	// initialize the window for s
	var head int = 0
	for i, ch := range s {
		window[ch] = window[ch] + 1
		if i < pLen-1 {
			continue
		}
		cmp := reflect.DeepEqual(window, check)
		if cmp {
			output = append(output, head)
		}

		log.WithFields(log.Fields{
			"check":   check,
			"window":  window,
			"compare": cmp,
			"index":   i,
			"rune":    string(ch),
			"ptr":    head,
		}).Debug("map result")

		// need to remove head key otherwise DeepEqual fails...
		r := rune(s[head])
		window[r] = window[r] - 1
		if window[r] <= 0 {
			delete(window,r)
		}
		head++
	}
	return output
}
