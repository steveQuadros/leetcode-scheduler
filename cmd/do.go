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
	"os"
	"encoding/json"
	"strings"
	// "time"

	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
	"github.com/stevequadros/studyplan/plan"
)

var filePath string
var doneAt string

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err  := os.Open(filePath)
		if err != nil {
			fmt.Println("error openign file: ", err.Error())
			os.Exit(1)
		}
		defer file.Close()

		dec := json.NewDecoder(file)
		// expecting an array of questions -> [
		dec.Token()

		// var questions []*plan.Question
		var sp = &plan.StudyPlan{}
		for dec.More() {
			if err := dec.Decode(sp); err != nil {
				fmt.Println("error decoding -> ", err.Error())
				os.Exit(1)
			}
		}
		questions := sp.Questions

		templates := &promptui.SelectTemplates{
			Label:    "Select A Question By Title",
			Active:   "\U0001F336 {{ .Title | cyan }} ({{ .Difficulty | red }})",
			Inactive: "  {{ .Title | cyan }} ({{ .Difficulty | red }})",
			Selected: "\U0001F336 {{ .Title | red | cyan }}",
			Details: `
--------- Question ----------
{{ "Title:" }}	{{ .Title }}
{{ "Difficulty:" }}	{{ .Difficulty }}
{{ "Link:" }}	{{ .Link }}`,
		}

		searcher := func(input string, index int) bool {
			question := questions[index]
			return strings.Contains(question.Title, input)
		}

		prompt := promptui.Select{
			Label: "Select Question",
			Items: questions,
			StartInSearchMode: true,
			Templates: templates,
			Size: 10,
			Searcher: searcher,
		}
	
		idx, result, err := prompt.Run()	
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You have marked as done %q\n", result)
		question := questions[idx]
		if startDate == "" {
			question.Do(sp.Intervals)
		} else {
			parsed, err := time.Parse(plan.ISOLayout, doneAt)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			question.DoAt(parsed, sp.Intervals)	
		}
			
		out, err := json.Marshal(sp)
		if err != nil {
			fmt.Println("error marshalling questions", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	doCmd.Flags().StringVarP(&filePath, "plan", "p", "examples/plan.json", "path to the study plan json file")
}
