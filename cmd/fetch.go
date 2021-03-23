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
	"bufio"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/stevequadros/studyplan/plan"
)

var limit = -1

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		var questionNames []string
		var count int
		for scanner.Scan() {
			if limit != -1 && count >= limit {
				break
			}
			count++
			name := scanner.Text()
			questionNames = append(questionNames, name)
		}
		
		if err := scanner.Err(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		questions, err := plan.Fetch(questionNames)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		out, err := json.Marshal(questions)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	fetchCmd.Flags().IntVarP(&limit, "limit", "l", -1, "limit number of fetches, defaults to -1 (no limit)")
}
