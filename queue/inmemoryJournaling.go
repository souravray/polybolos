/*
* @Author: souravray
* @Date:   2014-10-26 20:52:28
* @Last Modified by:   souravray
* @Last Modified time: 2014-10-26 23:42:12
 */

package queue

import (
	"container/heap"
	"fmt"
)

type JournalingInmemoryQueue struct {
	InmemoryQueue
}

func NewJournalingInimemoryQueue() Queue {
	tq := JournalingInmemoryQueue{InmemoryQueue{make(PriorityQueue, 0),
		make(map[string]*Task, 0),
		make(map[string]*Task, 0)}}
	heap.Init(&tq)
	return &tq
}

func (tq *JournalingInmemoryQueue) PushTask(task *Task) {
	fmt.Println("Journaling Push - ", task.Path)
	tq.InmemoryQueue.PushTask(task)
}

func (tq *JournalingInmemoryQueue) PopTask() *Task {
	task := tq.InmemoryQueue.PopTask()
	if task.Path != "" {
		fmt.Println("Journaling Pop - ", task.Path)
	}
	return task
}
