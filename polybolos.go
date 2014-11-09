/*
* Πολυβολος
* @Author: souravray
* @Date:   2014-10-11 19:52:00
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-09 22:55:02
 */

package polybolos

import (
	"fmt"
	q "github.com/souravray/polybolos/queue"
	w "github.com/souravray/polybolos/worker"
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
	w.Interface
}

type Queue struct {
	q.Queue
	queType   QueueType
	bucket    *Bucket
	queueRate int32
	workers   map[string]*WorkerResource
	stop      chan bool
}

var queue *Queue = nil

func GetQueue(qtype QueueType, maxConcurrentWorker int32, maxDequeueRate int32) (*Queue, error) {
	if queue == nil {
		bucket, err := NewBucket(maxConcurrentWorker, maxDequeueRate)
		if err != nil {
			return nil, err
		}
		queueRate := int32(math.Ceil(float64(maxDequeueRate)))
		switch qtype {
		case INMEMORY:
			queue = &Queue{
				q.NewInimemoryQueue(),
				qtype,
				bucket,
				queueRate,
				make(map[string]*WorkerResource),
				make(chan bool)}
		case INMEMORY_JOURNALING:
			queue = &Queue{
				q.NewJournalingInimemoryQueue(),
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
				if q.Len() > 0 {
					item := q.PopTask()
					if item.Path != "" {
						fmt.Println(item)
						worker, _ := q.workers[item.Path]
						worker.Interface.Perform(item.Payload)
					}
				}
				q.bucket.Spend()
			}
		}
	}
}

func (q *Queue) Delete() bool {
	if q == queue {
		q.bucket.Close()
		queue = nil
		return true
	}
	return false
}

func (q *Queue) AddHTTPWorker(name string, url url.URL, method w.HTTPWorkerMethod, retryLimit int32, ageLimit, minBackoff, maxBackoff time.Duration, maxDoubling bool) {
	worker := &w.HTTPWorker{w.Config{retryLimit, ageLimit, minBackoff, maxBackoff, maxDoubling},
		url,
		method}
	q.workers[name] = &WorkerResource{worker}
}

func (q *Queue) AddLocalWorker(name string, instance w.Interface, retryLimit int32, ageLimit, minBackoff, maxBackoff time.Duration, maxDoubling bool) {
	worker := &w.LocalWorker{w.Config{retryLimit, ageLimit, minBackoff, maxBackoff, maxDoubling},
		instance}
	q.workers[name] = &WorkerResource{worker}
}

func NewTask(path string, payload url.Values, delay string, eta time.Time) (task *q.Task) {
	task, _ = q.NewTask(path, payload, delay, eta)
	return task
}
