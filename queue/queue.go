package queue

import (
	"time"

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
		err := func() error {
			u, err := q.selectUpload()
			if err != nil {
				return err
			}
			if u == nil {
				return nil
			}
			return q.processUpload(u)
		}()
		var retryChan <-chan time.Time
		if err != nil {
			q.log.Error(err.Error())
			retryChan = time.After(30 * time.Second)
		}
		select {
		case <-retryChan:
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
	select {
	case q.triggerChan <- true:
	default:
	}
}

// Close shuts down the queue.
func (q *Queue) Close() {
	close(q.stopChan)
	<-q.stoppedChan
}
