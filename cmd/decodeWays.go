package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// decodeWaysCmd represents the decodeWays command
var decodeWaysCmd = &cobra.Command{
	Use:   "decodeWays",
	Short: "decode the number of ways a number string can be converted into letters",
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		val := decode(0,args[0])
		fmt.Println("number of ways =", val)
	},
}

func init() {
	rootCmd.AddCommand(decodeWaysCmd)
}

var memo = map[int]int{}

func decode(index int, s string) int {
	if v, ok := memo[index]; ok {
		return v
	}
	if index == len(s) {
		return 1
	}
	if s[0] == '0' {
		return 0
	}
	if index == len(s)-1 {
		return 1
	}

	ans := decode(index+1,s)
	val, _ := strconv.Atoi(s[index:index+2])
	if val <= 26 {
		ans += decode(index+2,s)
	}
	memo[index] = ans

	return ans
}
