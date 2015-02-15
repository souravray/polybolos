/*
* @Author: souravray
* @Date:   2015-02-16 00:54:54
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-16 04:15:32
 */

package polybolos

import (
	Q "github.com/souravray/polybolos/queue"
	"github.com/souravray/polybolos/sys"
	"math"
	"time"
)

type standardQueueType int

const (
	INMEMORY standardQueueType = iota
	INMEMORY_JOURNALING
)

type workerResource struct {
	Worker
}

type Queue interface {
	Start()
	Delete() bool
	RegisterWorker(name string, worker Worker)
	AddTask(task Task)
	RemoveTask(task Task)
}

type standardQueue struct {
	Q.Interface
	queType   standardQueueType
	bucket    *bucket
	queueRate int32
	workers   map[string]*workerResource
	stop      chan bool
}

var queue *standardQueue = nil

func newQueue(qtype standardQueueType) (tq Q.Interface) {
	switch qtype {
	case INMEMORY:
		tq = Q.NewInimemoryQueue()
	case INMEMORY_JOURNALING:
		tq = Q.NewJournalingInimemoryQueue()
	}
	return tq
}

func GetQueue(qtype standardQueueType, maxConcurrentWorker int32, maxDequeueRate int32) (Queue, error) {
	if queue == nil {
		maxConcurrentWorker = getMaxConcurrentWorker(maxConcurrentWorker)
		b, err := newBucket(maxConcurrentWorker, maxDequeueRate)
		if err != nil {
			return nil, err
		}
		queueRate := int32(math.Ceil(float64(maxDequeueRate/3)) * 2)
		queue = &standardQueue{
			newQueue(qtype),
			qtype,
			b,
			queueRate,
			make(map[string]*workerResource),
			make(chan bool)}
	}
	return queue, nil
}

func (q *standardQueue) Start() {
	q.bucket.Fill()
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
					w, _ := q.workers[task.Worker]
					go q.dispatch(w.Worker, task)
				} else {
					q.bucket.Spend()
				}
			}
		}
	}
}

func (q *standardQueue) dispatch(w Worker, task *Q.Task) {
	defer q.bucket.Spend()
	err := w.Perform(task.Payload)
	if err != nil {
		q.reenqueue(w, task)
	} else {
		q.done(task)
	}
	q.bucket.Spend()
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

func (q *standardQueue) Delete() bool {
	if q == queue {
		q.Close()
		q.bucket.Close()
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
