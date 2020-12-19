package parser

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ulricksennick/lcfetch/urls"
	"github.com/ulricksennick/lcfetch/util"
)

const (
	gqlOperation string = "getQuestionDetail"
	gqlQuery            = `query getQuestionDetail($titleSlug: String!) {
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

type ProblemDetails struct {
	CodeDefinitions map[string]string
	Likes           int
	Dislikes        int
	Content         string // TODO: figure out structure and missing data
	Stats           string // TODO: figure this out too..
	SampleTestCase  string // TODO: find other test cases, setup tests
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

type CodeDefinitionList []struct {
	Value       string `json:"value"`
	Test        string `json:"test"`
	DefaultCode string `json:"defaultCode"`
}

type RequestPayload struct {
	OperationName string           `json:"operationName"`
	Query         string           `json:"query"`
	Variables     PayloadVariables `json:"variables"`
}

type PayloadVariables struct {
	TitleSlug string `json:"titleSlug"`
}

func GetProblemDetails(titleSlug string) *ProblemDetails {
	req := buildProblemDetailsRequest(titleSlug)

	client := &http.Client{}
	resp, err := client.Do(req)
	must(err)
	defer resp.Body.Close()

	problemDetailsResponse := new(ProblemDetailsResponse)
	err = json.NewDecoder(resp.Body).Decode(problemDetailsResponse)
	must(err)

	codeDefinitionList := new(CodeDefinitionList)
	json.Unmarshal([]byte(problemDetailsResponse.Data.Question.CodeDefinition),
		codeDefinitionList)

	// Map languages to code definitions
	codeDefinitions := make(map[string]string)
	for _, lang := range *codeDefinitionList {
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

func buildProblemDetailsRequest(titleSlug string) *http.Request {
	requestPayload := createRequestPayload(titleSlug)

	var buf bytes.Buffer
	buf.Write(requestPayload)

	req, err := http.NewRequest("POST", urls.ProblemDetailsUrl, &buf)
	must(err)

	// Generate a CSRF token by sending a request to Leetcode
	csrfToken := util.GetCSRFToken()
	if len(csrfToken) == 0 {
		panic("Unable to generate CSRF token from Leetcode...")
	}
	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: csrfToken})

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", urls.LeetcodeBaseUrl)
	req.Header.Set("Referer", urls.LeetcodeProblemUrl+titleSlug)

	return req
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
