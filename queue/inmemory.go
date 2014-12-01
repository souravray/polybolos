/*
* @Author: souravray
* @Date:   2014-10-11 19:50:44
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-17 09:01:22
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
	tq.mutex.Lock()
	task, _ := heap.Pop(tq).(*Task)
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

func (tq *InmemoryQueue) Close() {
	//cleanup codes need to be go here
}
