package db

import "testing"

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
