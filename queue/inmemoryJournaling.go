/*
* @Author: souravray
* @Date:   2014-10-26 20:52:28
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-09 21:53:38
 */

package queue

import (
	"fmt"
	"github.com/souravray/polybolos/queue/heap"
)

type JournalingInmemoryQueue struct {
	InmemoryQueue
}

func NewJournalingInimemoryQueue() Queue {
	tq := JournalingInmemoryQueue{InmemoryQueue{make(DelayedQueue, 0),
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

func (tq *JournalingInmemoryQueue) DeleteTask(task *Task) {
	tq.InmemoryQueue.DeleteTask(task)
	fmt.Println("Journaling Delete - ", task.Path)
}
