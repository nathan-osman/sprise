package queue

import (
	"github.com/nathan-osman/sprise/db"
	"github.com/sirupsen/logrus"
)

// Queue manages uploads to S3 / Spaces.
type Queue struct {
	log         *logrus.Entry
	conn        *db.Conn
	triggerChan chan bool
	stopChan    chan bool
	stoppedChan chan bool
}

// run executes the main loop for the queue.
func (q *Queue) run() {
	defer close(q.stoppedChan)
	defer q.log.Info("queue has stopped")
	q.log.Info("starting queue")
	for {
		select {
		case <-q.triggerChan:
		case <-q.stopChan:
			return
		}
	}
}

// New creates a new upload queue.
func New(cfg *Config) *Queue {
	q := &Queue{
		log:         logrus.WithField("context", "queue"),
		conn:        cfg.Conn,
		triggerChan: make(chan bool, 1),
		stopChan:    make(chan bool),
		stoppedChan: make(chan bool),
	}
	go q.run()
	return q
}

// Trigger indicates that a new item has been queued.
func (q *Queue) Trigger() {
	q.triggerChan <- true
}

// Close shuts down the queue.
func (q *Queue) Close() {
	close(q.stopChan)
	<-q.stoppedChan
}
