/*
* Πολυβολος
* @Author: souravray
* @Date:   2014-10-11 19:52:00
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-17 09:06:36
 */

package polybolos

import (
	Q "github.com/souravray/polybolos/queue"
	W "github.com/souravray/polybolos/worker"
	"math"
	"net/url"
	"time"
)

type QueueType int

const (
	INMEMORY QueueType = iota
	INMEMORY_JOURNALING
)

type WorkerResource struct {
	W.Interface
}

type Queue struct {
	Q.Interface
	queType   QueueType
	bucket    *Bucket
	queueRate int32
	workers   map[string]*WorkerResource
	stop      chan bool
}

var queue *Queue = nil

func NewQueue(qtype QueueType) (queue Q.Interface) {
	switch qtype {
	case INMEMORY:
		queue = Q.NewInimemoryQueue()
	case INMEMORY_JOURNALING:
		queue = Q.NewJournalingInimemoryQueue()
	}
	return queue
}

func GetQueue(qtype QueueType, maxConcurrentWorker int32, maxDequeueRate int32) (*Queue, error) {
	if queue == nil {
		bucket, err := NewBucket(maxConcurrentWorker, maxDequeueRate)
		if err != nil {
			return nil, err
		}
		queueRate := int32(math.Ceil(float64(maxDequeueRate/3)) * 2)
		switch qtype {
		case INMEMORY:
			queue = &Queue{
				Q.NewInimemoryQueue(),
				qtype,
				bucket,
				queueRate,
				make(map[string]*WorkerResource),
				make(chan bool)}
		case INMEMORY_JOURNALING:
			queue = &Queue{
				Q.NewJournalingInimemoryQueue(),
				qtype,
				bucket,
				queueRate,
				make(map[string]*WorkerResource),
				make(chan bool)}
		}
	}
	return queue, nil
}

func (q *Queue) Start() {
	go q.bucket.Fill()
	for {
		select {
		case <-q.stop:
			q.bucket.Close()
			return
		default:
			n := <-q.bucket.Take(q.queueRate)
			for i := int32(0); i < n; i++ {
				item := q.PopTask()
				if item.Worker != "" {
					w, _ := q.workers[item.Worker]
					go func(q *Queue, w W.Interface, payload url.Values) {
						defer q.bucket.Spend()
						err := w.Perform(payload)
						if err != nil {

						} else {

						}
					}(q, w.Interface, item.Payload)
				}
			}
		}
	}
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
	worker = &W.HTTPWorker{W.Config{DefaultRetryLimit, DefaultAgeLimit, DefaultMinBackoff, DefaultMaxBackoff, DefaulrMaxDoubling},
		url,
		string(method)}
	return worker
}

func NewLocalWorker(instance W.Interface) (worker W.Interface) {
	worker = &W.LocalWorker{W.Config{DefaultRetryLimit, DefaultAgeLimit, DefaultMinBackoff, DefaultMaxBackoff, DefaulrMaxDoubling},
		instance}
	return worker
}

func (q *Queue) AddWorker(name string, worker W.Interface) {
	q.workers[name] = &WorkerResource{worker}
}

func NewTask(path string, payload url.Values, delay string, eta time.Time) (task *Q.Task) {
	task, _ = Q.NewTask(path, payload, delay, eta)
	return task
}
