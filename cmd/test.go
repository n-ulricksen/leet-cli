/*
Copyright © 2021 Nicholas Ulricksen <n.ulricksen@gmail.com>

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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ulricksennick/lcfetch/db"
	"github.com/ulricksennick/lcfetch/parser"
	"github.com/ulricksennick/lcfetch/problem"
	"github.com/ulricksennick/lcfetch/urls"
	"github.com/ulricksennick/lcfetch/util"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run testcases",
	Long:  "Run testcases.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		codeFilename := args[0]
		id := getIdFromFilename(codeFilename)

		database, err := db.CreateDB()
		must(err)
		prob, err := database.GetProblemByDisplayId(id)
		must(err)

		probDetails := parser.GetProblemDetails(prob.Slug)
		typedCode, err := getUserTypedCode(codeFilename)
		must(err)

		language := getLangFromFilename(codeFilename)
		body := createTestRequestBody(prob.DisplayId, probDetails.SampleTestCase,
			language, typedCode)
		cookies := database.GetAllCookies()

		// Interpret solution
		httpRequest := newTestRequest(prob, body, cookies)
		client := &http.Client{}
		resp, err := client.Do(httpRequest)
		must(err)
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Invalid leetcode credentials. Please login using the `lcfetch login` command.\n")
			os.Exit(1)
		}
		defer resp.Body.Close()
		respBody, err := decodeTestResponseBody(resp.Body)
		must(err)

		cmdOutput := ""

		// Repeatedly check problem status
		checkStatusRequest := newCheckStatusRequest(prob, cookies, respBody.InterpretId)
		for {
			time.Sleep(500 * time.Millisecond)
			resp, err = client.Do(checkStatusRequest)
			defer resp.Body.Close()
			statusRespBody, err := decodeCheckStatusResponseBody(resp.Body)
			must(err)
			if statusRespBody.State == "PENDING" {
			}
			if statusRespBody.State == "STARTED" {
				// TODO: let the user know the server is checking thier code
			}
			if statusRespBody.State == "SUCCESS" {
				s := statusRespBody

				color := 0
				if s.CorrectAnswer {
					color = 82 // green
				} else {
					color = 196 // red
					if s.StatusMessage == "Accepted" {
						s.StatusMessage = "Wrong Answer"
					}
				}
				s.StatusMessage = fmt.Sprintf("\033[38;5;%dm%s\033[m",
					color, s.StatusMessage)

				cmdOutput += fmt.Sprintf("%s\nOutput: %v\nCode Answer: %v\nExpected Answer: %v\nRuntime: %s\nMemory Usage: %s\n",
					s.StatusMessage, s.CodeOutput, s.CodeAnswer, s.ExpectedCodeAnswer,
					s.StatusRuntime, s.StatusMemory)
				if s.RuntimeError != "" {
					cmdOutput += fmt.Sprintf("\nError:\n\t%s\n", s.RuntimeError)
				}
				if s.CompileError != "" {
					cmdOutput += fmt.Sprintf("\nError:\n\t%s\n", s.CompileError)
				}

				fmt.Print(cmdOutput)

				break
			}
		}

	},
}

func newTestRequest(prob *problem.Problem, body []byte, cookies map[string]string) *http.Request {
	var bodyBuf bytes.Buffer
	bodyBuf.Write(body)

	reqUrl := prob.Url + "/interpret_solution/"
	req, err := http.NewRequest("POST", reqUrl, &bodyBuf)
	must(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", urls.LeetcodeBaseUrl)
	req.Header.Set("Referer", prob.Url)

	req.Header.Set("x-csrftoken", cookies["csrftoken"])
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	return req
}

func newCheckStatusRequest(prob *problem.Problem, cookies map[string]string, interpretId string) *http.Request {
	reqUrl := "https://leetcode.com/submissions/detail/%s/check/"
	reqUrl = fmt.Sprintf(reqUrl, interpretId)
	req, err := http.NewRequest("GET", reqUrl, nil)
	must(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", urls.LeetcodeBaseUrl)
	req.Header.Set("Referer", prob.Url)

	req.Header.Set("x-csrftoken", cookies["csrftoken"])
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	return req
}

type testRequestBody struct {
	QuestionId int    `json:"question_id"`
	DataInput  string `json:"data_input"`
	Language   string `json:"lang"`
	TypedCode  string `json:"typed_code"`
}

func createTestRequestBody(id int, dataInput, language, typedCode string) []byte {
	reqBody := testRequestBody{
		QuestionId: id,
		DataInput:  dataInput,
		Language:   language,
		TypedCode:  typedCode,
	}
	jsn, err := json.Marshal(reqBody)
	must(err)

	return jsn
}

type testResponseBody struct {
	InterpretId string `json:"interpret_id"`
	TestCase    string `json:"test_case"`
}

func decodeTestResponseBody(body io.Reader) (*testResponseBody, error) {
	decoded := new(testResponseBody)
	err := json.NewDecoder(body).Decode(decoded)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

type checkStatusResponseBody struct {
	CodeAnswer         []string `json:"code_answer"`
	CodeOutput         []string `json:"code_output"`
	CompileError       string   `json:"compile_error"`
	CorrectAnswer      bool     `json:"correct_answer"`
	ExpectedCodeAnswer []string `json:"expected_code_answer"`
	RunSuccess         bool     `json:"run_success"`
	RuntimeError       string   `json:"runtime_error"`
	State              string   `json:"state"`
	StatusMemory       string   `json:"status_memory"`
	StatusMessage      string   `json:"status_msg"`
	StatusRuntime      string   `json:"status_runtime"`
}

func decodeCheckStatusResponseBody(body io.Reader) (*checkStatusResponseBody, error) {
	decoded := new(checkStatusResponseBody)
	err := json.NewDecoder(body).Decode(decoded)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func getUserTypedCode(codeFilename string) (string, error) {
	path, err := filepath.Abs(codeFilename)
	if err != nil {
		return "", err
	}
	userCode, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(userCode), nil
}

func getIdFromFilename(file string) int {
	toks := strings.Split(file, "-")
	id := int(toks[0][0] - '0')

	return id
}

func getLangFromFilename(file string) string {
	toks := strings.Split(file, ".")
	extension := toks[len(toks)-1]
	language := util.FileType[extension]

	return language
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
