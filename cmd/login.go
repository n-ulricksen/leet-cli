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
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/ulricksennick/leet-cli/db"
	"github.com/ulricksennick/leet-cli/util"
)

var supportedBrowsers = map[string]bool{
	"chromium": true,
	"chrome":   false,
	"edge":     false,
	"safari":   false,
}

var requiredLoginCookies = []string{"csrftoken", "LEETCODE_SESSION"}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Leetcode",
	Long:  "Login to Leetcode using your browser's cookies.\nFirst, login to leetcode.com using your web browser, then close the web browser, and run `lcfetch login <browser>` again.",
	Args:  validateLoginArgs, // Web browser
	Run: func(cmd *cobra.Command, args []string) {
		browser := args[0]
		cookies, err := util.GetCredentialCookies(browser)
		must(err)

		// verify csrf and session cookies are present
		for _, cookie := range requiredLoginCookies {
			if _, ok := cookies[cookie]; !ok {
				log.Fatalf("Could not retrieve `%s` from %s cookies.\nPlease login to leetcode.com using your web browser, close the web browser, and run `lcfetch login <browser>` again.", cookie, browser)
			}
		}

		database, err := db.CreateDB()
		must(err)
		err = database.InsertCookies(cookies)
		must(err)
		fmt.Println("Successfully logged in.")
	},
}

// Should receive 1 argument, the name of the user's web browser.
func validateLoginArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires a browser argument\n" + getSupportedBrowsers())
	}

	browser := args[0]
	if isSupported := supportedBrowsers[browser]; !isSupported {
		return errors.New("unsupported browser...\n" + getSupportedBrowsers())
	}

	return nil
}

func getSupportedBrowsers() string {
	var supported []string
	for b, s := range supportedBrowsers {
		if s {
			supported = append(supported, b)
		}
	}

	return fmt.Sprintf("supported browsers: %s\n", supported)
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
