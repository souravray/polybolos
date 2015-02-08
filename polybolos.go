/*
* Πολυβολος
* @Author: souravray
* @Date:   2014-10-11 19:52:00
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-08 09:29:39
 */

package polybolos

import (
	Q "github.com/souravray/polybolos/queue"
	"github.com/souravray/polybolos/sys"
	W "github.com/souravray/polybolos/worker"
	"math"
	"net/url"
	"time"
)

type queueType int

const (
	INMEMORY queueType = iota
	INMEMORY_JOURNALING
)

type WorkerResource struct {
	W.Interface
}

type Queue struct {
	Q.Interface
	queType   queueType
	bucket    *bucket
	queueRate int32
	workers   map[string]*WorkerResource
	stop      chan bool
}

var queue *Queue = nil

func NewQueue(qtype queueType) (queue Q.Interface) {
	switch qtype {
	case INMEMORY:
		queue = Q.NewInimemoryQueue()
	case INMEMORY_JOURNALING:
		queue = Q.NewJournalingInimemoryQueue()
	}
	return queue
}

func getWorkerCapacity(maxConcurrentWorker int32) int32 {
	maxFDLimit := uint64(maxConcurrentWorker * 2)
	newLimit, _ := sys.SetFDLimits(maxFDLimit)
	if newLimit > 0 && newLimit < maxFDLimit {
		maxConcurrentWorker = int32(math.Floor(float64(newLimit / 2)))
	}
	return maxConcurrentWorker
}

func GetQueue(qtype queueType, maxConcurrentWorker int32, maxDequeueRate int32) (*Queue, error) {
	if queue == nil {
		maxConcurrentWorker = getWorkerCapacity(maxConcurrentWorker)
		b, err := newBucket(maxConcurrentWorker, maxDequeueRate)
		if err != nil {
			return nil, err
		}
		queueRate := int32(math.Ceil(float64(maxDequeueRate/3)) * 2)
		switch qtype {
		case INMEMORY:
			queue = &Queue{
				Q.NewInimemoryQueue(),
				qtype,
				b,
				queueRate,
				make(map[string]*WorkerResource),
				make(chan bool)}
		case INMEMORY_JOURNALING:
			queue = &Queue{
				Q.NewJournalingInimemoryQueue(),
				qtype,
				b,
				queueRate,
				make(map[string]*WorkerResource),
				make(chan bool)}
		}
	}
	return queue, nil
}

func (q *Queue) Start() {
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
					go q.dispatch(w.Interface, task)
				} else {
					q.bucket.Spend()
				}
			}
		}
	}
}

func (q *Queue) dispatch(w W.Interface, task *Q.Task) {
	defer q.bucket.Spend()
	err := w.Perform(task.Payload)
	if err != nil {
		q.reenqueue(w, task)
	} else {
		q.done(task)
	}
	q.bucket.Spend()
}

func (q *Queue) reenqueue(w W.Interface, task *Q.Task) {
	retryLimit := w.GetRetryLimit()
	ageLimit := w.GetAgeLimit()
	retryAttempt := task.RetryCount + 1
	taskAge := time.Since(task.EnqueTime)

	if retryLimit == int32(0) || retryLimit > retryAttempt {
		if ageLimit == time.Duration(0) || ageLimit > taskAge {
			task.MinDelay = w.GetInterval(retryAttempt)
			task.RetryCount = retryAttempt
			q.PushTask(task)
		} else {
			q.done(task)
		}
	} else {
		if ageLimit == time.Duration(0) || taskAge >= ageLimit {
			q.done(task)
		} else {
			task.MinDelay = w.GetInterval(retryAttempt)
			task.RetryCount = retryAttempt
			q.PushTask(task)
		}
	}
}

func (q *Queue) done(task *Q.Task) {

}

func (q *Queue) Delete() bool {
	if q == queue {
		q.Close()
		q.bucket.Close()
		queue = nil
		return true
	}
	return false
}

func NewHTTPWorker(url url.URL, method HTTPWorkerMethod) (worker W.Interface) {
	worker = &W.HTTPWorker{W.Config{DefaultWorkerTimeout, DefaultRetryLimit, DefaultAgeLimit, DefaultMinBackoff, DefaultMaxBackoff, DefaulrMaxDoubling},
		url,
		string(method)}
	return worker
}

func NewLocalWorker(instance W.Interface) (worker W.Interface) {
	worker = &W.LocalWorker{W.Config{DefaultWorkerTimeout, DefaultRetryLimit, DefaultAgeLimit, DefaultMinBackoff, DefaultMaxBackoff, DefaulrMaxDoubling},
		instance}
	return worker
}

func (q *Queue) RegisterWorker(name string, worker W.Interface) {
	q.workers[name] = &WorkerResource{worker}
}

func NewTask(path string, payload url.Values, delay string, eta time.Time) (task *Q.Task) {
	task, _ = Q.NewTask(path, payload, delay, eta)
	return task
}
