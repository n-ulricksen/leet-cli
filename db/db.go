package db

import (
	"encoding/json"
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

func (db *DB) SetProblemCompleted(displayId int) error {
	var updateId int
	coll := db.conn.Use(collectionName)

	// Lookup problem by displayId
	var err error
	coll.ForEachDoc(func(id int, doc []byte) bool {
		prob, e := documentToProblem(doc)
		if e != nil {
			err = e
			return false
		}

		if prob.DisplayId == displayId {
			updateId = id
			return false
		}
		return true
	})
	if err != nil {
		return err
	}

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

func (db *DB) SetQuestionBad(displayId int) error {
	var updateId int
	coll := db.conn.Use(collectionName)

	// Lookup problem by displayId
	var err error
	coll.ForEachDoc(func(id int, doc []byte) bool {
		prob, e := documentToProblem(doc)
		if e != nil {
			err = e
			return false
		}

		if prob.DisplayId == displayId {
			updateId = id
			return false
		}
		return true
	})
	if err != nil {
		return err
	}

	doc, err := coll.Read(updateId)
	if err != nil {
		return err
	}
	doc["badQuestion"] = true

	err = coll.Update(updateId, doc)
	if err != nil {
		return err
	}

	return nil
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
