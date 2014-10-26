/*
* @Author: souravray
* @Date:   2014-10-11 19:50:44
* @Last Modified by:   souravray
* @Last Modified time: 2014-10-27 00:32:12
 */

package queue

import (
	"container/heap"
)

type InmemoryQueue struct {
	PriorityQueue
	OneTheflyQueue map[string]*Task
	DoneQueue      map[string]*Task
}

func NewInimemoryQueue() Queue {
	tq := InmemoryQueue{make(PriorityQueue, 0),
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
