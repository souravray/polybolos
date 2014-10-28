/*
* @Author: souravray
* @Date:   2014-10-11 19:52:00
* @Last Modified by:   souravray
* @Last Modified time: 2014-10-29 02:02:27
 */

package polybolos

import (
	"fmt"
	q "github.com/souravray/polybolos/queue"
	"math"
	"net/url"
)

type QueueType int

const (
	INMEMORY QueueType = iota
	INMEMORY_JOURNALING
)

type Queue struct {
	q.Queue
	queType   QueueType
	bucket    *Bucket
	queueRate int32
	workers   map[string]*Worker
	stop      chan bool
}

var queue *Queue = nil

func GetQueue(qtype QueueType, maxConcurrentWorker int32, maxDequeueRate int32) (*Queue, error) {
	if queue == nil {
		bucket, err := NewBucket(maxConcurrentWorker, maxDequeueRate)
		if err != nil {
			return nil, err
		}
		queueRate := int32(math.Ceil(float64(maxDequeueRate / 2)))
		switch qtype {
		case INMEMORY:
			queue = &Queue{
				q.NewInimemoryQueue(),
				qtype,
				bucket,
				queueRate,
				make(map[string]*Worker),
				make(chan bool)}
		case INMEMORY_JOURNALING:
			queue = &Queue{
				q.NewJournalingInimemoryQueue(),
				qtype,
				bucket,
				queueRate,
				make(map[string]*Worker),
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
					}
					q.bucket.Spend()
				}
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

func NewTask(path string, payload url.Values, delay string) (task *q.Task) {
	task, _ = q.NewTask(path, payload, delay)
	return task
}
