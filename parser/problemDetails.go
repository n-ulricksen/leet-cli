package parser

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	problemDetailsUrl string = "https://leetcode.com/graphql"
	gqlOperation      string = "getQuestionDetail"
	gqlQuery          string = `query getQuestionDetail($titleSlug: String!) {
               question(titleSlug: $titleSlug) {
                 content
                 stats
                 likes
                 dislikes
                 codeDefinition
                 sampleTestCase
                 enableRunCode
                 metaData
				 translatedContent
               }
            }`
)

type RequestPayload struct {
	OperationName string           `json:"operationName"`
	Query         string           `json:"query"`
	Variables     PayloadVariables `json:"variables"`
}

type PayloadVariables struct {
	TitleSlug string `json:"titleSlug"`
}

type ProblemDetailsResponse struct {
	Data struct {
		Question struct {
			Content        string `json:"content"`
			Stats          string `json:"stats"`
			Likes          int    `json:"likes"`
			Dislikes       int    `json:"dislikes"`
			CodeDefinition string `json:"codeDefinition"`
			SampleTestCase string `json:"sampleTestCase"`
			EnableRunCode  bool   `json:"enableRunCode"`
			MetaData       string `json:"metaData"`
		} `json:"question"`
	} `json:"data"`
}

type CodeDefinition []struct {
	Value       string `json:"value"`
	Test        string `json:"test"`
	DefaultCode string `json:"defaultCode"`
}

type ProblemDetails struct {
	CodeDefinitions map[string]string
	Likes           int
	Dislikes        int
	Content         string // TODO: figure out structure and missing data
	Stats           string // TODO: figure this out too..
	SampleTestCase  string // TODO: find other test cases, setup tests
}

func GetProblemDetails(titleSlug string) *ProblemDetails {
	requestPayload := createRequestPayload(titleSlug)

	var buf bytes.Buffer
	buf.Write(requestPayload)

	req, err := http.NewRequest("POST", problemDetailsUrl, &buf)
	must(err)

	// Generate a CSRF token by sending a request to Leetcode
	csrfToken := getCSRFToken()
	if len(csrfToken) == 0 {
		panic("Unable to generate CSRF token from Leetcode...")
	}
	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: csrfToken})

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", leetcodeBaseUrl)
	req.Header.Set("Referer", leetcodeProblemUrl+titleSlug)

	client := &http.Client{}
	resp, err := client.Do(req)
	must(err)
	defer resp.Body.Close()

	problemDetailsResponse := new(ProblemDetailsResponse)
	err = json.NewDecoder(resp.Body).Decode(problemDetailsResponse)
	must(err)

	codeDefinition := new(CodeDefinition)
	json.Unmarshal([]byte(problemDetailsResponse.Data.Question.CodeDefinition),
		codeDefinition)

	// Map languages to code definitions
	codeDefinitions := make(map[string]string)
	for _, lang := range *codeDefinition {
		codeDefinitions[lang.Value] = lang.DefaultCode
	}

	return &ProblemDetails{
		CodeDefinitions: codeDefinitions,
		Likes:           problemDetailsResponse.Data.Question.Likes,
		Dislikes:        problemDetailsResponse.Data.Question.Dislikes,
		Content:         problemDetailsResponse.Data.Question.Content,
		Stats:           problemDetailsResponse.Data.Question.Stats,
		SampleTestCase:  problemDetailsResponse.Data.Question.SampleTestCase,
	}
}

func createRequestPayload(titleSlug string) []byte {
	jsn, err := json.Marshal(&RequestPayload{
		OperationName: gqlOperation,
		Query:         gqlQuery,
		Variables:     PayloadVariables{titleSlug},
	})
	if err != nil {
		log.Println(err)
	}

	return jsn
}

func getCSRFToken() string {
	resp, err := http.Get(leetcodeBaseUrl)
	must(err)
	defer resp.Body.Close()

	var csrfToken string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "csrftoken" {
			csrfToken = cookie.Value
			break
		}
	}
	return csrfToken
}
