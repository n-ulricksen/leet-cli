package db

import (
	"encoding/json"
	"errors"

	tiedot "github.com/HouzuoGuo/tiedot/db"
	// "github.com/HouzuoGuo/tiedot/dberr"
	"github.com/ulricksennick/lcfetch/problem"
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

// Insert a problem into the database, returns the problem's database ID
func (db *DB) InsertProblem(problem *problem.Problem) (int, error) {
	coll := db.conn.Use(collectionName)

	problemJson := problemToJsonMap(problem)

	docID, err := coll.Insert(problemJson)
	if err != nil {
		return -1, err
	}

	return docID, err
}

// Get a slice of pointers to all problems in the database
func (db *DB) GetAllProblems() ([]*problem.Problem, error) {
	var ret []*problem.Problem
	var err error

	coll := db.conn.Use(collectionName)

	coll.ForEachDoc(func(id int, doc []byte) bool {
		prob, e := documentToProblem(doc)
		if e != nil {
			err = e
			return false
		}
		ret = append(ret, prob)
		return true
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Set a problem's "completed" field as true
func (db *DB) SetProblemCompleted(displayId int) error {
	updateId, err := db.getProblemId(displayId)
	if err != nil {
		return err
	}

	coll := db.conn.Use(collectionName)

	doc, err := coll.Read(updateId)
	if err != nil {
		return err
	}
	doc["completed"] = true

	err = coll.Update(updateId, doc)
	if err != nil {
		return err
	}

	return nil
}

// Set a problem's "completed" field as false
func (db *DB) SetProblemIncomplete(displayId int) error {
	updateId, err := db.getProblemId(displayId)
	if err != nil {
		return err
	}

	coll := db.conn.Use(collectionName)

	doc, err := coll.Read(updateId)
	if err != nil {
		return err
	}
	doc["completed"] = false

	err = coll.Update(updateId, doc)
	if err != nil {
		return err
	}

	return nil
}

// Set a problem's "isBad" field as true
func (db *DB) SetProblemBad(displayId int) error {
	updateId, err := db.getProblemId(displayId)
	if err != nil {
		return err
	}

	coll := db.conn.Use(collectionName)
	doc, err := coll.Read(updateId)
	if err != nil {
		return err
	}
	doc["isBad"] = true

	err = coll.Update(updateId, doc)
	if err != nil {
		return err
	}

	return nil
}

// Drop the "problems" collection, create a new empty one
func (db *DB) DropAllProblems() {
	db.conn.Drop(collectionName)
	db.conn.Create(collectionName)
}

// Lookup a problem ID in the database by its leetcode displayId
func (db *DB) getProblemId(displayId int) (int, error) {
	var problemId int = -1
	var err error

	coll := db.conn.Use(collectionName)

	coll.ForEachDoc(func(id int, doc []byte) bool {
		prob, e := documentToProblem(doc)
		if e != nil {
			err = e
			return false
		}

		if prob.DisplayId == displayId {
			problemId = id
			return false
		}
		return true
	})
	if err != nil {
		return -1, err
	}

	if problemId == -1 {
		return -1, errors.New("No problem found with given id.")
	}

	return problemId, nil
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
