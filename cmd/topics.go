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
	"github.com/ulricksennick/lcfetch/problem"
)

// topicsCmd represents the topics command
var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "List all problem topics on Leetcode.",
	Long:  `List all problem topics on Leetcode.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Number of columns to split topic list into when printing
		const columnCount = 3

		sorted := problem.GetSortedTopics()

		var outBuf bytes.Buffer
		outBuf.WriteString("Leetcode problem topics:\n")
		i := 0
		for _, topic := range sorted {
			if i%columnCount == 0 {
				outBuf.WriteByte('\n')
			}
			outBuf.WriteString(fmt.Sprintf("%-26s", topic))
			i++
		}
		outBuf.WriteByte('\n')
		fmt.Print(outBuf.String())
	},
}

func init() {
	rootCmd.AddCommand(topicsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// topicsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// topicsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
