package main

import (
	"database/sql"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

func CreateDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "feedbacks.db")
	if err != nil {
		log.Fatal("Unable to open/create database", "err", err)
	}

	sql := `CREATE TABLE if not exists feedbacks (
          host TEXT PRIMARY KEY,  
          updates INTEGER 
       );`

	_, err = db.Exec(sql)
	if err != nil {
		log.Warn("Unable to create first table", "err", err)
	}

	sql = `CREATE TABLE if not exists answers (
          question TEXT,
		  answer TEXT,
		  feedbackID TEXT,
		  UNIQUE(question, feedbackID) ON CONFLICT REPLACE,
		  FOREIGN KEY(feedbackID) REFERENCES feedbacks(host)
       );`

	_, err = db.Exec(sql)
	if err != nil {
		log.Warn("Unable to create second table", "err", err)
	}

	return db
}

func InsertFeedback(db *sql.DB, f Feedback) {
	query, err := db.Prepare(`
		INSERT INTO feedbacks(host, updates) VALUES (?, 0) 
			ON CONFLICT(host) DO UPDATE SET updates=updates+1`)
	if err != nil {
		log.Fatal("Unable to prepare feedback upsert statement", "err", err)
	}
	_, err = query.Exec(f.Host)
	if err != nil {
		log.Fatal("Unable to upsert feedback from", "host", f.Host, "err", err)
	}

	InsertAnswers(db, f)

	log.Infof("%s has submitted a new feedback!")
}

func InsertAnswers(db *sql.DB, f Feedback) {
	for _, question := range f.Questions {
		a := AnswerEntity{
			Question:   question.Title,
			Answer:     question.Answer.Value(),
			FeedbackID: f.Host,
		}

		query, err := db.Prepare(`
			INSERT INTO answers(question, answer, feedbackID) VALUES (?, ?, ?) 
				ON CONFLICT (question, feedbackID) DO UPDATE 
				SET answer=excluded.answer`)
		if err != nil {
			log.Fatal("Unable to prepare answer upsert query", "err", err)
		}

		_, err = query.Exec(a.Question, a.Answer, a.FeedbackID)
		if err != nil {
			log.Fatal("Failed to upsert answer", "question", a.Question, "answer", a.Answer, "feedbackID", a.FeedbackID, "err", err)
		}

	}
}

type AnswerEntity struct {
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	FeedbackID string `json:"feedbackID"`
}

type FeedbackEntity struct {
	Host    string `json:"host"`
	Updates int    `json:"updates"`
}
