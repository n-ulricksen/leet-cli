package util

import (
	"log"
	"net/http"

	"github.com/ulricksennick/leet-cli/urls"
)

func GetCSRFToken() string {
	resp, err := http.Get(urls.LeetcodeBaseUrl)
	if err != nil {
		log.Fatal(err)
	}
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
