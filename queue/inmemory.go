/*
* @Author: souravray
* @Date:   2014-10-11 19:50:44
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-02 20:59:03
 */

package queue

import (
	"github.com/souravray/polybolos/queue/heap"
)

type InmemoryQueue struct {
	PriorityWaitQueue
	OneTheflyQueue map[string]*Task
	DoneQueue      map[string]*Task
}

func NewInimemoryQueue() Queue {
	tq := InmemoryQueue{make(PriorityWaitQueue, 0),
		make(map[string]*Task, 0),
		make(map[string]*Task, 0)}
	heap.Init(&tq)
	return &tq
}

func (tq *InmemoryQueue) PushTask(task *Task) {
	heap.Push(tq, task)
}

func (tq *InmemoryQueue) PopTask() *Task {
	task, _ := heap.Pop(tq).(*Task)
	return task
}

func (tq *InmemoryQueue) DeleteTask(task *Task) {
	if task.index >= 0 {
		task, _ = heap.Remove(tq, task.index).(*Task)
	}
}
