/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"time"

	"github.com/apex/log"
	"github.com/spf13/cobra"

	"gitlab.com/mshindle/tidbits/montecarlo"
)

var (
	points int64
	numWorkers int
)

// piCmd represents the pi command
var piCmd = &cobra.Command{
	Use:   "pi",
	Short: "calculate the value of PI",
	Long: `
Calculate PI using the Monte Carlo method. This method consists of drawing on a canvas a square with an inner circle. 
We then generate a large number of random points within the square and count how many fall in the enclosed circle.

The area of the circle is πr^2,
The area of the square is width2 = (2r)^2 = 4r^2.

If we divide the area of the circle, by the area of the square we get π/4.

The same ratio can be used between the number of points within the square and the number of points within the circle.

Hence we can use the following formula to estimate Pi: π ~ 4 * (number of points within circle / number of points total)`,
	RunE: pi,
}

func init() {
	rootCmd.AddCommand(piCmd)
	piCmd.Flags().Int64VarP(&points,"points", "p", 50000000, "number of points to use for calculation")
	piCmd.Flags().IntVarP(&numWorkers,"workers", "w", 3, "number of workers to calculate if points are in or out of circle")
}

func pi(cmd *cobra.Command, args []string) error {
	log.WithField("points",points).WithField("numWorkers", numWorkers).Info("pi computation started")
	start := time.Now()
	p := montecarlo.NewPI(points, numWorkers)
	err := p.Compute()
	log.WithField("duration", time.Now().Sub(start)).Info("elapsed time")
	return err
}