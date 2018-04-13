package main

import (
	"github.com/mshindle/tidbits/cmd"
)

func main() {
	cmd.Execute()

	//app := cli.NewApp()
	//app.Name = "tidbits"
	//app.Usage = "execute samples of learning code"
	//app.Author = "Mike Shindle"
	//app.Email = "mshindle@gmail.com"
	//app.Version = "0.0.1"
	//app.Commands = []cli.Command{
	//
	//	{
	//		Name:      "breaker",
	//		Aliases:   []string{"b"},
	//		Usage:     "run a circuit breaker example",
	//		ArgsUsage: "host1 [host2...]",
	//		Action: func(c *cli.Context) error {
	//			fmt.Println("Running breaker =>")
	//			retry.RunBreaker(c.Args()...)
	//			return nil
	//		},
	//	},
	//	{
	//		Name:    "limit",
	//		Aliases: []string{"l"},
	//		Usage:   "run a request rate limiter example",
	//		Action: func(c *cli.Context) error {
	//			fmt.Println("Running limiter =>")
	//			limit.RunRequest()
	//			return nil
	//		},
	//	},
	//}
}
