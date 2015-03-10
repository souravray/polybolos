/*
* @Author: souravray
* @Date:   2015-02-16 00:54:54
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-21 07:17:05
 */

package polybolos

import (
	"errors"
	"fmt"
	Q "github.com/souravray/polybolos/queue"
	"github.com/souravray/polybolos/sys"
	"math"
	"os"
	"path"
	"time"
)

type standardQueueType int

const (
	INMEMORY standardQueueType = iota
	INMEMORY_JOURNALING
)

const (
	DEFAULT_QUEUE string = "default"
)

type workerResource struct {
	Worker
}

// Interface for the push Queue, the standardQueue structure implements the
// interface. You can start dispatching task to worker by calling Start(),
// and can pause by calling Pause(). Calling Close() will close the queue
// and return a boolean true on success. RegisterWorker register a new worker
// to the worker pool. New task can be added or removed from the Queue using
// AddTask and RemoveTask.
type Queue interface {
	Start()
	Pause()
	Close() bool
	RegisterWorker(name string, worker Worker)
	AddTask(task Task)
	RemoveTask(task Task)
}

type standardQueue struct {
	Q.Interface                            // delayed queue (priority queue)
	queType     standardQueueType          // queutye INMEMORY or INMEORY_JOURNALING
	bucket      *bucket                    // token bucket
	queueRate   int32                      // rate, at wich the queue will fetch token
	workers     map[string]*workerResource // worker pool
	stop        chan bool                  // quit channel
}

var queue *standardQueue = nil

func newQueue(qtype standardQueueType, journalPath, queueName string) (tq Q.Interface) {
	switch qtype {
	case INMEMORY:
		tq = Q.NewInimemoryQueue()
	case INMEMORY_JOURNALING:
		tq = Q.NewJournalingInimemoryQueue(journalPath, DEFAULT_QUEUE)
	}
	return tq
}

func validatedJournalPath(journalPath string) (string, error) {
	cleanPath := path.Clean(journalPath)
	pathInfo, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("Journal Path does not exist")
		}
		return "", err
	}

	if !pathInfo.IsDir() {
		return "", errors.New("Journal Path is not a  directory")
	}

	return cleanPath, nil
}

func GetQueue(qtype standardQueueType, maxConcurrentWorker, maxDequeueRate int32, journalPath string) (Queue, error) {
	var err error
	if qtype == INMEMORY_JOURNALING {
		journalPath, err = validatedJournalPath(journalPath)
		if err != nil {
			return nil, err
		}
	}

	if queue == nil {
		maxConcurrentWorker = getMaxConcurrentWorker(maxConcurrentWorker)
		b, err := newBucket(maxConcurrentWorker, maxDequeueRate)
		if err != nil {
			return nil, err
		}
		queueRate := int32(math.Ceil(float64(maxDequeueRate/3)) * 2)
		queue = &standardQueue{
			Interface: newQueue(qtype, journalPath, DEFAULT_QUEUE),
			queType:   qtype,
			bucket:    b,
			queueRate: queueRate,
			workers:   make(map[string]*workerResource),
		}
	}
	return queue, err
}

func (q *standardQueue) Start() {
	q.stop = make(chan bool, 1)
	q.bucket.Fill()
	go func(q *standardQueue) {
		defer close(q.stop)
		for {
			select {
			case <-q.stop:
				q.bucket.Close()
				return
			default:
				n := <-q.bucket.Take(q.queueRate)
				for i := int32(0); i < n; i++ {
					task := q.PopTask()
					if task.Worker != "" {
						w, ok := q.workers[task.Worker]
						if ok {
							go q.dispatch(w.Worker, task)
						} else {
							fmt.Println("Not found worker ", task.Worker)
							q.done(task)
							q.bucket.Spend()
						}

					} else {
						q.bucket.Spend()
					}
				}
			}
		}
	}(q)
}

func (q *standardQueue) dispatch(w Worker, task *Q.Task) {
	defer q.bucket.Spend()
	err := w.Perform(task.Payload)
	if err != nil {
		q.reenqueue(w, task)
	} else {
		q.done(task)
	}
}

func (q *standardQueue) reenqueue(w Worker, task *Q.Task) {
	retryLimit := w.GetRetryLimit()
	ageLimit := w.GetAgeLimit()
	retryAttempt := task.RetryCount + 1
	taskAge := time.Since(task.EnqueTime)
	taskEnque := func(q *standardQueue, w Worker, task *Q.Task, retryAttempt int32) {
		delay := w.GetInterval(retryAttempt)
		task.ETA = time.Now().Add(delay)
		task.RetryCount = retryAttempt
		q.PushTask(task)
	}

	if retryLimit == int32(0) || retryLimit > retryAttempt {
		if ageLimit == time.Duration(0) || ageLimit > taskAge {
			taskEnque(q, w, task, retryAttempt)
		} else {
			q.done(task)
		}
	} else {
		if ageLimit == time.Duration(0) || taskAge >= ageLimit {
			q.done(task)
		} else {
			taskEnque(q, w, task, retryAttempt)
		}
	}
}

func (q *standardQueue) done(task *Q.Task) {
	q.CleanTask(task)
}

func (q *standardQueue) Pause() {
	q.stop <- true
}

func (q *standardQueue) Close() bool {
	if q == queue {
		q.stop <- true
		q.Interface.Close()
		queue = nil
		return true
	}
	return false
}

func (q *standardQueue) RegisterWorker(name string, worker Worker) {
	q.workers[name] = &workerResource{worker}
}

func (q *standardQueue) AddTask(task Task) {
	if t, ok := task.(*Q.Task); ok {
		q.PushTask(t)
	}
}

func (q *standardQueue) RemoveTask(task Task) {
	if t, ok := task.(*Q.Task); ok {
		q.DeleteTask(t)
	}
}

func getMaxConcurrentWorker(maxConcurrentWorker int32) int32 {
	maxFDLimit := uint64(maxConcurrentWorker * 2)
	newLimit, _ := sys.SetFDLimits(maxFDLimit)
	if newLimit > 0 && newLimit < maxFDLimit {
		maxConcurrentWorker = int32(math.Floor(float64(newLimit / 2)))
	}
	return maxConcurrentWorker
}
