/*
* @Author: souravray
* @Date:   2014-10-11 19:50:44
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-08 23:09:53
 */

package queue

import (
	"github.com/souravray/polybolos/queue/heap"
	"sync"
)

type InmemoryQueue struct {
	DelayedQueue
	OneTheflyQueue map[string]*Task
	mutex          sync.Mutex
}

func NewInimemoryQueue() Interface {
	tq := InmemoryQueue{DelayedQueue: make(DelayedQueue, 0),
		OneTheflyQueue: make(map[string]*Task, 0)}
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
	//task clean-up code
}

func (tq *InmemoryQueue) Close() {
	//queue clean-up codes
}
