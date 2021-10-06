package db

import (
	"encoding/json"
	"errors"
	"sort"

	tiedot "github.com/HouzuoGuo/tiedot/db"
	"github.com/ulricksennick/lcfetch/problem"
)

const (
	defaultDbFilePath  = "/tmp/leetcode"
	problemsCollection = "problems"
	topicsCollection   = "topics"
	cookiesCollection  = "cookies"
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
	conn.Create(problemsCollection)
	conn.Create(topicsCollection)
	conn.Create(cookiesCollection)
	return &DB{conn}, nil
}

// Insert multiple problems into the database
func (db *DB) InsertProblems(problems map[int]*problem.Problem) error {
	var err error

	for _, problem := range problems {
		_, err = db.InsertProblem(problem)
		if err != nil {
			return err
		}
	}

	return nil
}

// Insert a problem into the database, returns the problem's database ID
func (db *DB) InsertProblem(problem *problem.Problem) (int, error) {
	coll := db.conn.Use(problemsCollection)

	problemJson := problemToJsonMap(problem)

	docID, err := coll.Insert(problemJson)
	if err != nil {
		return -1, err
	}

	return docID, err
}

// Insert topics into the database
func (db *DB) InsertTopics(topics []*problem.Topic) error {
	var err error

	for _, topic := range topics {
		_, err = db.InsertTopic(topic)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) InsertTopic(topic *problem.Topic) (int, error) {
	coll := db.conn.Use(topicsCollection)

	topicJson := topicToJsonMap(topic)

	docID, err := coll.Insert(topicJson)
	if err != nil {
		return -1, err
	}

	return docID, err
}

// Get a slice of pointers to all topics in the database
func (db *DB) GetAllTopics() ([]*problem.Topic, error) {
	var ret []*problem.Topic
	var err error

	coll := db.conn.Use(topicsCollection)

	coll.ForEachDoc(func(id int, doc []byte) bool {
		topic, e := documentToTopic(doc)
		if e != nil {
			err = e
			return false
		}
		ret = append(ret, topic)
		return true
	})

	return ret, err
}

func GetSortedTopicStrings() ([]string, error) {
	db, err := CreateDB()
	if err != nil {
		return nil, err
	}

	topics, err := db.GetAllTopics()
	if err != nil {
		return nil, err
	}

	sorted := make([]string, len(topics))
	for i, t := range topics {
		sorted[i] = t.Slug
	}
	sort.Strings(sorted)

	return sorted, nil
}

// Get a slice of pointers to all problems in the database
func (db *DB) GetAllProblems() ([]*problem.Problem, error) {
	var ret []*problem.Problem
	var err error

	coll := db.conn.Use(problemsCollection)

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

// Return a slice of pointers to one or more problems found by displayId
func (db *DB) GetProblemsByDisplayId(ids []int) ([]*problem.Problem, error) {
	var ret []*problem.Problem
	var err error

	coll := db.conn.Use(problemsCollection)

	// Used to check for target IDs when iterating over all problems
	displayIds := make(map[int]bool)
	for _, id := range ids {
		displayIds[id] = true
	}

	// Iterate over all problems checking displayId
	coll.ForEachDoc(func(id int, doc []byte) bool {
		prob, e := documentToProblem(doc)
		if e != nil {
			err = e
			return false
		}
		if displayIds[prob.DisplayId] {
			ret = append(ret, prob)
		}
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

	coll := db.conn.Use(problemsCollection)

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

	coll := db.conn.Use(problemsCollection)

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

	coll := db.conn.Use(problemsCollection)
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

// Insert items into cookies table as a single document
func (db *DB) InsertCookies(cookies map[string]string) error {
	db.DropAllCookies()

	var err error
	var doc map[string]interface{}

	jsn, err := json.Marshal(cookies)
	if err != nil {
		return err
	}
	json.Unmarshal(jsn, &doc)

	coll := db.conn.Use(cookiesCollection)
	_, err = coll.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAllCookies() map[string]string {
	cookies := make(map[string]string)

	coll := db.conn.Use(cookiesCollection)
	coll.ForEachDoc(func(id int, doc []byte) bool {
		json.Unmarshal(doc, &cookies)
		return true
	})

	return cookies
}

// Drop the "problems" collection, create a new empty one
func (db *DB) DropAllProblems() {
	db.conn.Drop(problemsCollection)
	db.conn.Create(problemsCollection)
}

// Drop the "topics" collection, create a new empty one
func (db *DB) DropAllTopics() {
	db.conn.Drop(topicsCollection)
	db.conn.Create(topicsCollection)
}

// Drop the "cookies" collection, create a new empty one
func (db *DB) DropAllCookies() {
	db.conn.Drop(cookiesCollection)
	db.conn.Create(cookiesCollection)
}

// Lookup a problem ID in the database by its leetcode displayId
func (db *DB) getProblemId(displayId int) (int, error) {
	var problemId int = -1
	var err error

	coll := db.conn.Use(problemsCollection)

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

func topicToJsonMap(topic *problem.Topic) map[string]interface{} {
	var jsonMap map[string]interface{}

	jsn, _ := json.Marshal(topic)
	json.Unmarshal(jsn, &jsonMap)

	return jsonMap
}

func documentToTopic(doc []byte) (*problem.Topic, error) {
	newTopic := &problem.Topic{}
	err := json.Unmarshal(doc, newTopic)
	if err != nil {
		return nil, err
	}
	return newTopic, nil
}
