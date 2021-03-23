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
	"os"
	"github.com/spf13/cobra"
	"github.com/stevequadros/studyplan/plan"
)

var outputAsPlan bool
var intervals []int
var startDate string
var questionsPerDay int

// genCmd represents the gen command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		d := json.NewDecoder(os.Stdin)
		d.Token()
		var questions []*plan.Question
		for d.More() {
			q := &plan.Question{}
			if err := d.Decode(q); err != nil {
				fmt.Println("error decoding; ", err)
				os.Exit(1)
			}
			questions = append(questions, q)
		}

		if len(questions) < 1 {
			fmt.Println("no questions to schedule")
			return
		}

		var startTime time.Time
		var err error
		startTime, err = time.Parse(plan.ISOLayout, startDate)
		if err != nil {
			fmt.Println("err parsing time: ", err)
			os.Exit(1)
		}
		sp := plan.New(intervals, startTime, questionsPerDay)
		sp.Schedule(questions)

		var out []byte
		if outputAsPlan {
			// plans := plan.ByDate(plan.questions)
			// out, err = json.Marshal(plans)
			// if err != nil {
			// 	panic(err)
			// }
		} else {
			out, err = json.Marshal([]*plan.StudyPlan{sp})
			if err != nil {
				fmt.Println("err marshalling plan: ", err)
				os.Exit(1)
			}
		}
		fmt.Print(string(out))
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	planCmd.Flags().BoolVarP(&outputAsPlan, "plan", "p", false, "output as plan - defaults to ouputting as questions")
	planCmd.Flags().IntSliceVarP(&intervals, "intervals", "i", []int{2,4,8,16}, "comma separated list of interval days for repeat - ex. 2,4,8,16")
	planCmd.Flags().StringVarP(&startDate, "start", "s", time.Now().Format(plan.ISOLayout), "start date for plan")
	planCmd.Flags().IntVarP(&questionsPerDay, "questions-per-day", "q", 3, "new questions per day")
}
