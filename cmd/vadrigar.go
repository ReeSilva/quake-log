/*
Copyright © 2020 Renato Biancalana da Silva <reesilva@pm.me>

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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/reesilva/quake-log/pkg/output"
	"github.com/reesilva/quake-log/pkg/parser"
	"github.com/spf13/cobra"
)

var (
	meanOfDeath bool
	logFile     string
	outputFile  string
)

// vadrigarCmd represents the vadrigar command
var vadrigarCmd = &cobra.Command{
	Use:   "vadrigar",
	Short: "Vadrigar will parse a file that contains logs for a specific Quake 3 Arena Server",
	Long: `With vadrigar command you will parse an entire Quake 3 Arena servers and receive,
in stdout or in a file, the logs for each game structured in JSON. You will also be able to 
activate an option to show to you the number of deaths by mean in each game.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(logFile)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		matches := []parser.Match{}
		for scanner.Scan() {
			matched, err := regexp.MatchString(`\d+:\d+ (\w+): (.*)`, scanner.Text())
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			if !matched {
				continue
			}
			parser.ParseLine(len(matches)-1, &matches, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		report, err := output.CreateMatchReport(matches, meanOfDeath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		j, err := json.MarshalIndent(report, "", "\t")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if outputFile != "" {
			ioutil.WriteFile(outputFile, j, 0644)
			os.Exit(0)
		}
		fmt.Println(string(j))
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(vadrigarCmd)

	vadrigarCmd.Flags().BoolVarP(&meanOfDeath, "mean-of-death", "m", false, "Enable or disable logs of deaths by mean")
	vadrigarCmd.Flags().StringVarP(&logFile, "log-file", "f", "", "Path for the Quake 3 Arena Server logs file")
	vadrigarCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Output file. If not set, will print as JSON in stdout")
	vadrigarCmd.MarkFlagRequired("log-file")
}
