/*
Copyright © 2020 Nicholas Ulricksen <n.ulricksen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulricksennick/leet-cli/problem"
	"github.com/ulricksennick/leet-cli/util"
)

// Command line flags
var listTopics []string
var listDifficulty string
var listIncludePaid bool
var listFilterCompleted bool
var listFilterOutCompelted bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print a list of leetcode problems.",
	Long: `Print a list of the Leetcode problems, filtered by difficulty and/or topic.

Examples:
  'lcfetch list'
  'lcfetch list -d easy -t array,string'`,
	Run: func(cmd *cobra.Command, args []string) {
		problemSet := getFilteredProblemSet(listDifficulty, listTopics, listIncludePaid)

		if listFilterCompleted && listFilterOutCompelted {
			fmt.Println("'-i' and '-c' flags set, pick one or the other!")
			return
		}

		// Filter complete/incomplete
		if listFilterCompleted {
			problemSet = problem.FilterCompleted(problemSet)
		}
		if listFilterOutCompelted {
			problemSet = problem.FilterOutCompleted(problemSet)
		}

		if len(problemSet) == 0 {
			fmt.Println("No problems found, try widening your search.")
			return
		}

		var listBuf bytes.Buffer
		listBuf.WriteString("-----------------------------------------\n")
		listBuf.WriteString("ID\tDone?\tDifficulty\tName\t\t|\n")
		listBuf.WriteString("-----------------------------------------\n")
		for _, problem := range problemSet {
			completedCh := ' '
			if problem.Completed {
				// TODO: change this 'x' to a checkmark emoji
				completedCh = 'x'
			}

			difficulty := util.GetColoredDifficultyText(problem.Difficulty)

			listBuf.WriteString(fmt.Sprintf("%d\t%c\t%s\t%s\n",
				problem.DisplayId, completedCh, difficulty, problem.Name))
		}
		fmt.Print(listBuf.String())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	listCmd.Flags().StringSliceVarP(&listTopics, "topics", "t", []string{},
		"topic(s) of problems to list (comma-separated, no spaces)")
	listCmd.Flags().StringVarP(&listDifficulty, "difficulty", "d", "all",
		"difficulty of problems to list")
	listCmd.Flags().BoolVarP(&listIncludePaid, "paid", "p", false,
		"include paid/premium problems")
	listCmd.Flags().BoolVarP(&listFilterCompleted, "completed", "c", false,
		"list only completed problems")
	listCmd.Flags().BoolVarP(&listFilterOutCompelted, "incomplete", "i", false,
		"list only incomplete problems")
}
