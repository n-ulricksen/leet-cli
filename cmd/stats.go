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
	"github.com/ulricksennick/lcfetch/db"
	"github.com/ulricksennick/lcfetch/problem"
	"github.com/ulricksennick/lcfetch/util"
)

var statsIncludePaid bool

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Print details about completed questions per category and difficulty.",
	Long:  `Print details about completed questions per category and difficulty.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Number of columns wide to print stats to screen
		const columnCount = 2

		database, err := db.CreateDB()
		must(err)

		problemSet, err := database.GetAllProblems()
		must(err)

		// Filter paid problems unless flag is set
		if !statsIncludePaid {
			problemSet = problem.FilterOutPaid(problemSet)
		}

		// Create a map of stats categories, each containing a slice [x,y],
		// where x is the # of complete problems and y is the number of total
		// problems per category.
		completedCount := make(map[string][]int)
		completedCount["all"] = make([]int, 2)
		completedCount["easy"] = make([]int, 2)
		completedCount["medium"] = make([]int, 2)
		completedCount["hard"] = make([]int, 2)
		for _, topic := range problem.GetSortedTopics() {
			completedCount[topic] = make([]int, 2)
		}

		for _, prob := range problemSet {
			// Total (all problems)
			if prob.Completed {
				completedCount["all"][0]++
			}
			completedCount["all"][1]++

			// Difficulty
			var diff string
			switch prob.Difficulty {
			case 1:
				diff = "easy"
			case 2:
				diff = "medium"
			case 3:
				diff = "hard"
			}
			if prob.Completed {
				completedCount[diff][0]++
			}
			completedCount[diff][1]++

			// Topic
			for _, topic := range prob.Topics {
				if prob.Completed {
					completedCount[topic][0]++
				}
				completedCount[topic][1]++
			}
		}

		var outBuf bytes.Buffer
		outBuf.WriteString("Leetcode problem statistics:\n\n")

		// Print difficulties
		outBuf.WriteString(fmt.Sprintf("Easy: %d/%d\n", completedCount["easy"][0],
			completedCount["easy"][1]))
		outBuf.WriteString(fmt.Sprintf("Medium: %d/%d\n", completedCount["medium"][0],
			completedCount["medium"][1]))
		outBuf.WriteString(fmt.Sprintf("Hard: %d/%d\n", completedCount["hard"][0],
			completedCount["hard"][1]))

		// Print topics
		sorted := problem.GetSortedTopics()
		i := 0
		for _, topic := range sorted {
			if i%columnCount == 0 {
				outBuf.WriteByte('\n')
			}
			outBuf.WriteString(fmt.Sprintf("%-22s%d/%d\t", util.KebabToCapital(topic),
				completedCount[topic][0], completedCount[topic][1]))
			i++
		}
		fmt.Print(outBuf.String())
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	statsCmd.Flags().BoolVarP(&statsIncludePaid, "paid", "p", false,
		"include paid/premium questions")
}
