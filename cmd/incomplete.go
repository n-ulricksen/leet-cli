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
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ulricksennick/leet-cli/db"
)

// incompleteCmd represents the incomplete command
var incompleteCmd = &cobra.Command{
	Use:   "incomplete",
	Short: "Mark one or more problems incomplete.",
	Long: `Mark more ore more problems, allowing them to show up when requesting a random problem.

Example:
  'lcfetch incomplete 1337'
  'lcfetch incomplete 628 12 52'`,
	Run: func(cmd *cobra.Command, args []string) {
		database, err := db.CreateDB()
		must(err)

		for _, arg := range args {
			problemId, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Invalid problem ID: %v\n", arg)
				continue
			}
			database.SetProblemIncomplete(problemId)
			fmt.Printf("Problem #%v marked incomplete.\n", problemId)
		}
	},
}

func init() {
	rootCmd.AddCommand(incompleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// incompleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// incompleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
