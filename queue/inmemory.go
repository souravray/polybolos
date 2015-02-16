/*
* @Author: souravray
* @Date:   2014-10-11 19:50:44
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-17 00:51:44
 */

package queue

import (
	"github.com/souravray/polybolos/queue/heap"
	"sync"
)

type InmemoryQueue struct {
	DelayedQueue
	onTheFlyJobs sync.WaitGroup
	mutex        sync.Mutex
}

func NewInimemoryQueue() Interface {
	tq := InmemoryQueue{DelayedQueue: make(DelayedQueue, 0)}
	heap.Init(&tq)
	return &tq
}

func (tq *InmemoryQueue) PushTask(task *Task) {
	tq.mutex.Lock()
	heap.Push(tq, task)
	tq.mutex.Unlock()
}

func (tq *InmemoryQueue) PopTask() *Task {
	var task *Task
	tq.mutex.Lock()
	if tq.Len() > 0 {
		task, _ = heap.Pop(tq).(*Task)
		tq.onTheFlyJobs.Add(1)
	} else {
		task = new(Task)
	}
	tq.mutex.Unlock()
	return task
}

func (tq *InmemoryQueue) DeleteTask(task *Task) {
	if task.index >= 0 {
		tq.mutex.Lock()
		task, _ = heap.Remove(tq, task.index).(*Task)
		tq.mutex.Unlock()
	}
}

func (tq *InmemoryQueue) CleanTask(task *Task) {
	tq.onTheFlyJobs.Done()
}

func (tq *InmemoryQueue) Close() {
	tq.onTheFlyJobs.Wait()
}
