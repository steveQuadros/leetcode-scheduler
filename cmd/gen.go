/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"time"
	"encoding/csv"
	"os"
	"io"
	"github.com/spf13/cobra"
	"github.com/stevequadros/studyplan/plan"
)

var outputAsPlan bool

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var intervals = []int{2,4,8,16}
		var startDate = time.Now()
		var questionsPerDay = 3

		c := csv.NewReader(os.Stdin)
		// skip header
		_, err := c.Read()
		if err != nil {
			panic(err)
		}
		
		var allQuestions []*plan.Question
		for line, err := c.Read(); err != io.EOF; line, err = c.Read() {
			if err != nil {
				panic(err)
			}
			allQuestions = append(allQuestions, &plan.Question{Title: line[0], Link: line[1], Difficulty: line[2]})
		}

		if len(allQuestions) < 1 {
			fmt.Println("no questions to schedule")
			return
		}

		questions := plan.GenerateDates(allQuestions, intervals, startDate, questionsPerDay)
		var out []byte
		if outputAsPlan {
			plans := plan.ByDate(questions)
			out, err = json.Marshal(plans)
			if err != nil {
				panic(err)
			}
			
		} else {
			out, err = json.Marshal(questions)
			if err != nil {
				panic(err)
			}
		}
		fmt.Print(string(out))
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	genCmd.Flags().BoolVarP(&outputAsPlan, "plan", "p", false, "output as plan - defaults to ouputting as questions")
}
