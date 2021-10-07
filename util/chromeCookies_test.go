package util

import (
	"fmt"
	"testing"
)

func TestGetCredentialCookies(t *testing.T) {
	chromeCookies, err := GetCredentialCookies("chromium")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", chromeCookies)
}
