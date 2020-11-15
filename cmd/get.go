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
	"io/ioutil"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ulricksennick/lcfetch/db"
	"github.com/ulricksennick/lcfetch/parser"
	"github.com/ulricksennick/lcfetch/util"
)

var codeLanguage string

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Lookup one or more problems by Leetcode ID.",
	Long: `Lookup one or more problems by Leetcode ID.

Examples:
  'lcfetch get 521'
  'lcfetch get 72 1262 980'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		database, err := db.CreateDB()
		must(err)

		displayIds := make([]int, len(args))
		for i, arg := range args {
			problemId, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Invalid problem ID: %v\n", arg)
				continue
			}
			displayIds[i] = problemId
		}

		problemSet, err := database.GetProblemsByDisplayId(displayIds)
		if err != nil {
			fmt.Println(err)
			return
		}

		var outBuf bytes.Buffer
		for _, prob := range problemSet {
			outBuf.WriteString(fmt.Sprintf("#%d\t%s", prob.DisplayId, prob.Name))
			if prob.Paid {
				outBuf.WriteString(" (Premium)")
			}

			var difficulty string
			switch prob.Difficulty {
			case 1:
				difficulty = "Easy"
			case 2:
				difficulty = "Medium"
			case 3:
				difficulty = "Hard"
			}
			outBuf.WriteString(fmt.Sprintf("\n%s\t%s\n\n", difficulty, prob.Url))

			problemDetails := parser.GetProblemDetails(prob.Slug)
			sourceCode := problemDetails.CodeDefinitions[codeLanguage]
			if len(sourceCode) == 0 {
				var buf bytes.Buffer
				buf.WriteString(fmt.Sprintf("Invalid language: %v.\n\n",
					codeLanguage))
				buf.WriteString("Available languages:\n")
				for lang := range problemDetails.CodeDefinitions {
					buf.WriteString(lang)
					buf.WriteByte('\n')
				}
				fmt.Println(buf.String())
				os.Exit(1)
			}

			filename := fmt.Sprintf("%s.%s", prob.Slug, util.FileExt[codeLanguage])
			ioutil.WriteFile(filename, []byte(sourceCode), 0664)
			outBuf.WriteString(fmt.Sprintln("Problem code defintion stored at: " +
				filename))
		}
		fmt.Print(outBuf.String())
	},
}

func getMapKeys(m map[string]string) []string {
	ret := make([]string, len(m))
	i := 0
	for k := range m {
		ret[i] = k
		i++
	}
	return ret
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCmd.Flags().StringVarP(&codeLanguage, "language", "l", "javascript",
		"programming language for code definition")
}
