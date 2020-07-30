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
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ulricksennick/lcfetch/db"
	"github.com/ulricksennick/lcfetch/problem"

	// "github.com/ulricksennick/lcfetch/problem"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var difficulty string
var topics []string
var includePaid bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lcfetch",
	Short: "Program used to retrieve Leetcode problem URLs.",
	Long: `Get a random problem from Leetcode based on difficulty and/or topic.

Example:
	'lcfetch -d medium -t array,two-pointers'`,
	Run: func(cmd *cobra.Command, args []string) {
		// Fetch all of the problems, and print a random one after applying
		// the appropriate filters.
		database, err := db.CreateDB()
		must(err)

		problemSet, err := database.GetAllProblems()
		must(err)

		// Filter paid problems
		if !includePaid {
			problemSet = problem.FilterOutPaid(problemSet)
		}

		// Apply topic filters
		for _, topic := range topics {
			problemSet = problem.FilterByTopic(problemSet, topic)
		}
		if len(problemSet) == 0 {
			fmt.Println("No problems found with the provided topic...")
			fmt.Println("Run 'lcfetch list -t' to list all topics.")
			return
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
				fmt.Println("invalid difficulty rating... easy, medium, or hard")
				return
			}
			problemSet = problem.FilterByDifficulty(problemSet, difficultyRating)
		}
		if len(problemSet) == 0 {
			fmt.Println("No problems found with the provided topics/difficulty...")
			return
		}

		fmt.Println(len(problemSet))

		// Pick and print a random problem to the screen
		rand.Seed(time.Now().UnixNano())
		selected := problemSet[rand.Intn(len(problemSet))]
		fmt.Printf("#%d - %s\n", selected.DisplayId, selected.Name)
		fmt.Println(selected.Url)
	},
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
	rootCmd.Flags().StringVarP(&difficulty, "difficulty", "d", "all",
		"difficulty of problem to select")
	rootCmd.Flags().StringSliceVarP(&topics, "topics", "t", []string{},
		"topic(s) to select problem from (comma-separated, no spaces)")
	rootCmd.Flags().BoolVarP(&includePaid, "paid", "p", false,
		"include paid or premuim problems")

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
