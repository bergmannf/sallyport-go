package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var conn *sql.DB = nil

type Queue interface {
	Put(r Request) error
	Pop() (Request, error)
	NotificationChannel() (chan Request)
}

type SQLiteQueue struct {
	database *sql.DB
	Endpoint string
	RequestNotify chan Request
}

type MemoryQueue struct {
	RequestNotify chan Request
}

func persistRequest(s SQLiteQueue, r Request) error {
	stmnt, err := s.database.Prepare(
		fmt.Sprintf("INSERT INTO %s (requesturl, requestjson) VALUES(?,?)", s.Endpoint),
	)
	if err != nil {
		fmt.Println("Error creating PREPARED STATEMENT", err)
		return err
	}
	_, err = stmnt.Exec(r.Uri, r.json())
	if err != nil {
		fmt.Println("Error executing PREPARED STATEMENT", err)
		return err
	}
	return nil
}

func NewSQLiteQueue(databaseLocation string, endpoint string) SQLiteQueue {
	var err error
	if conn == nil {
		fmt.Println("Establishing new database connection")
		conn, err = sql.Open("sqlite3", databaseLocation)
		if err != nil {
		}
	}
	queue := SQLiteQueue{
		database: conn,
		RequestNotify: make(chan Request, 100),
	}
	queue.initialize(endpoint)
	return queue
}

func (s SQLiteQueue) Put(r Request) error {
	fmt.Println("Inserting new request.")
	err := persistRequest(s, r)
	s.RequestNotify <- r
	return err
}

func (s SQLiteQueue) Pop() (Request, error) {
	fmt.Println("Retrieving requests")
	r := <- s.RequestNotify
	return r, nil
}

func (s SQLiteQueue) initialize(endpoint string) error {
	fmt.Println("Creating required tables.")
	_, err := s.database.Query(`CREATE TABLE ?(
         requestnumber INTEGER PRIMARY KEY,
	 requesturl TEXT,
	 requestjson TEXT)`, endpoint)
	if err != nil {
		fmt.Println("Error creating required tables:", err)
		return err
	}
	return nil
}

func (s SQLiteQueue) NotificationChannel() (chan Request) {
	return s.RequestNotify
}

func NewMemoryQueue() MemoryQueue {
	return MemoryQueue {
		RequestNotify: make(chan Request, 100),
	}
}

func (q MemoryQueue) Put(r Request) error {
	fmt.Println("Inserting new request.")
	q.RequestNotify <- r
	return nil
}

func (q MemoryQueue) Pop() (Request, error) {
	fmt.Println("Retrieving requests")
	r := <- q.RequestNotify
	return r, nil
}

func (q MemoryQueue) NotificationChannel() (chan Request) {
	return q.RequestNotify
}
