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
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ulricksennick/lcfetch/db"
	"github.com/ulricksennick/lcfetch/parser"
	"github.com/ulricksennick/lcfetch/problem"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootDifficulty string
var rootTopics []string
var rootIncludePaid bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lcfetch",
	Short: "Program used to retrieve Leetcode problem URLs.",
	Long: `Get a random problem from Leetcode based on difficulty and/or topic.

Examples:
  'lcfetch -d hard -t dynamic-programming'
  'lcfetch -d medium -t array,two-pointers'`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		database, err := db.CreateDB()
		must(err)

		problemSet, err := database.GetAllProblems()
		must(err)

		// TODO: create a command to update problems/topics

		if len(problemSet) == 0 {
			fmt.Println("Updating problems...")
			parser.FetchAndStoreTopics()
			parser.FetchAndStoreProblems()
			fmt.Println("Update complete.")
			fmt.Println()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Fetch all of the problems, and print a random one after applying
		// the appropriate filters.
		problemSet := getFilteredProblemSet(rootDifficulty, rootTopics, rootIncludePaid)
		// Filter completed problems
		problemSet = problem.FilterOutCompleted(problemSet)
		if len(problemSet) == 0 {
			fmt.Println("No problems found with the provided topics/difficulty...")
			return
		}

		// Pick and print a random problem to the screen
		rand.Seed(time.Now().UnixNano())
		selected := problemSet[rand.Intn(len(problemSet))]
		fmt.Printf("#%d - %s\n", selected.DisplayId, selected.Name)
		fmt.Println(selected.Url)

		// Ask to save problem to file
		var getCmd *cobra.Command
		for _, c := range cmd.Commands() {
			if c.Name() == "get" {
				getCmd = c
				break
			}
		}

		fmt.Println("\nSave problem to file? (y/n)")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Bytes()

		if len(input) == 0 || input[0] == ' ' || input[0] == 'y' || input[0] == 'Y' {
			getCmd.Run(getCmd, []string{strconv.Itoa(selected.DisplayId)})
		}
	},
}

func getFilteredProblemSet(difficulty string, topics []string, includePaid bool) []*problem.Problem {
	database, err := db.CreateDB()
	must(err)

	problemSet, err := database.GetAllProblems()
	must(err)

	// Filter paid problems
	if !includePaid {
		problemSet = problem.FilterOutPaid(problemSet)
	}

	// Apply topic filters
	if len(topics) > 0 {
		for _, topic := range topics {
			if topic == "" {
				continue
			}
			problemSet = problem.FilterByTopic(problemSet, topic)
		}
	}
	if len(problemSet) == 0 {
		fmt.Println("No problems found with the provided topic...")
		fmt.Println("Run 'lcfetch topics' to list all topics.")
		return []*problem.Problem{}
	}

	// Apply difficulty filter
	if difficulty != "all" {
		var difficultyRating int
		switch strings.ToLower(difficulty) {
		case "easy":
			difficultyRating = problem.EASY
			break
		case "medium":
			difficultyRating = problem.MEDIUM
			break
		case "hard":
			difficultyRating = problem.HARD
			break
		default:
			fmt.Println("Invalid difficulty rating... easy, medium, or hard")
			return []*problem.Problem{}
		}
		problemSet = problem.FilterByDifficulty(problemSet, difficultyRating)
	}
	if len(problemSet) == 0 {
		fmt.Println("No problems found with the provided topics/difficulty...")
		return []*problem.Problem{}
	}

	return problemSet
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lcfetch.yaml)")

	// Local command flags
	rootCmd.Flags().StringVarP(&rootDifficulty, "difficulty", "d", "all",
		"difficulty of problem to select")
	rootCmd.Flags().StringSliceVarP(&rootTopics, "topics", "t", []string{},
		"topic(s) to select problem from (comma-separated, no spaces)")
	rootCmd.Flags().BoolVarP(&rootIncludePaid, "paid", "p", false,
		"include paid/premium problems")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".lcfetch" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lcfetch")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
