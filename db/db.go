package db

import (
	"encoding/json"
	"strings"
	// "fmt"

	tiedot "github.com/HouzuoGuo/tiedot/db"
	// "github.com/HouzuoGuo/tiedot/dberr"
	"github.com/ulricksennick/leetcode-fetcher/problem"
)

const (
	defaultDbFilePath = "/tmp/leetcode"
	collectionName    = "problems"
)

type DB struct {
	conn *tiedot.DB
}

// Create a database if it does exist
func CreateDB() (*DB, error) {
	// TODO: get optional user provided path from config file
	conn, err := tiedot.OpenDB(defaultDbFilePath)
	if err != nil {
		return nil, err
	}

	// Create collection of problems
	conn.Create(collectionName)
	return &DB{conn}, nil
}

func (db *DB) InsertProblem(problem *problem.Problem) (int, error) {
	coll := db.conn.Use(collectionName)

	problemJson := problemToJsonMap(problem)

	docID, err := coll.Insert(problemJson)
	if err != nil {
		return -1, err
	}

	return docID, err
}

func (db *DB) GetAllProblems() ([]*problem.Problem, error) {
	var ret []*problem.Problem
	var err error

	coll := db.conn.Use(collectionName)

	coll.ForEachDoc(func(id int, doc []byte) bool {
		prob, e := documentToProblem(doc)
		if e != nil {
			err = e
		}
		ret = append(ret, prob)
		return true
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// TODO: move these functionalities to filter functions, accepting the full
// list of problems, and returning a subset of filtered problems
func (db *DB) GetProblemsByDifficulty(difficulty int) ([]*problem.Problem, error) {
	var ret []*problem.Problem
	var err error

	coll := db.conn.Use(collectionName)

	coll.ForEachDoc(func(id int, doc []byte) bool {
		prob, e := documentToProblem(doc)
		if e != nil {
			err = e
		}

		if prob.Difficulty == difficulty {
			ret = append(ret, prob)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *DB) GetProblemsByTopic(topic string) ([]*problem.Problem, error) {
	var ret []*problem.Problem
	var err error

	// Lowercase search term
	topic = strings.ToLower(topic)

	coll := db.conn.Use(collectionName)

	coll.ForEachDoc(func(id int, doc []byte) bool {
		prob, e := documentToProblem(doc)
		if e != nil {
			err = e
		}

		for _, t := range prob.Topics {
			if t == topic {
				ret = append(ret, prob)
			}
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *DB) DropAllProblems() {
	db.conn.Drop(collectionName)
	db.conn.Create(collectionName)
}

// This function is so hacky... ＞︿＜
func problemToJsonMap(problem *problem.Problem) map[string]interface{} {
	var jsonMap map[string]interface{}

	jsn, _ := json.Marshal(problem)
	json.Unmarshal(jsn, &jsonMap)

	return jsonMap
}

func documentToProblem(doc []byte) (*problem.Problem, error) {
	newProb := &problem.Problem{}
	err := json.Unmarshal(doc, newProb)
	if err != nil {
		return nil, err
	}
	return newProb, nil
}
