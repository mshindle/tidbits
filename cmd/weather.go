/*
Copyright © 2020 Michael Shindle <mshindle@gmail.com>

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
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/mshindle/tidbits/weather"
	"gitlab.com/mshindle/tidbits/weather/openweather"
)

// weatherCmd represents the weather command
var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "a simple app that tells us how we should dress",
	Long: `
An example of implementing the strategy / provider pattern which allows us to
pull data from any external service that fits our defined interface. The pattern
is especially useful if we want to do testing by exchanging sources.`,
	SilenceUsage: true,
	Args: cobra.ExactArgs(1),
	RunE: func (cmd *cobra.Command, args []string) error {
		city := args[0]
		apikey := viper.GetString("apikey")

		log.WithField("city",city).WithField("apikey",apikey).Info("looking for weather")
		provider := openweather.New(apikey)
		svc := weather.NewForecastService(provider)

		options, err := svc.HowToDress(city)
		if err != nil {
			return err
		}
		log.WithField("options", options).Info("options to choose for going out today")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(weatherCmd)
	weatherCmd.PersistentFlags().StringP("apikey","k","","apikey for provider service")

	_ = viper.BindPFlag("apikey",weatherCmd.PersistentFlags().Lookup("apikey"))
}

