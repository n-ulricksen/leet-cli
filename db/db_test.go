package db

import (
	"fmt"
	"testing"
)

func TestInsertCookies(t *testing.T) {
	cookies := map[string]string{
		"hello":   "world",
		"goodbye": "space",
	}

	database, err := CreateDB()
	if err != nil {
		t.Error(err)
	}

	err = database.InsertCookies(cookies)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAllCookies(t *testing.T) {
	database, err := CreateDB()
	if err != nil {
		t.Error(err)
	}

	_ = database.GetAllCookies()

}

func TestSetProblemTestCase(t *testing.T) {
	testcase := "testcase #1"
	id := 235

	database, err := CreateDB()
	if err != nil {
		t.Fatal(err)
	}
	database.SetProblemTestCase(id, testcase)

	problem, err := database.GetProblemByDisplayId(id)
	if err != nil {
		t.Errorf("Unable to find problem with id %d\n", id)
	}

	if problem.SampleTestCase != testcase {
		t.Errorf("did not properly set sampleTestCase field")
	}

	fmt.Printf("found problem: %#v\n", problem)
}
