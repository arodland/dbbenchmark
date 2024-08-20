package main

import "time"

type Database struct {
	requests chan Request
}

func NewDatabase() Database {
	return Database{
		requests: make(chan Request),
	}
}

type Request struct {
	response chan Response
}
type Response struct{}

const maxConcurrentRequests = 3

// Run spawns a long-running Goroutine that processes requests.
func (d *Database) Run() chan struct{} {
	close := make(chan struct{})

	// Our database is modern and concurrent, so we can handle a number of
	// concurrent requests at once, but not too many. If too many come in, some
	// will have to wait.
	semaphore := make(chan struct{}, maxConcurrentRequests)

	go func() {
		for {
			select {
			case r := <-d.requests:
				semaphore <- struct{}{}
				go func() {
					d.processRequest(r)
					<-semaphore
				}()
			case <-close:
				return
			}
		}
	}()

	return close
}

// Request is called by clients to simulate making a request to the database.
func (d *Database) Request() {
	response := make(chan Response)
	d.requests <- Request{response: response}
	<-response
}

func (d *Database) processRequest(r Request) {
	time.Sleep(10 * time.Millisecond)
	r.response <- Response{}
}
