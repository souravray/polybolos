/*
* @Author: souravray
* @Date:   2014-10-26 20:52:28
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-17 09:08:11
 */

package queue

import (
	"bytes"
	"encoding/gob"
	"github.com/souravray/polybolos/queue/db"
	"github.com/souravray/polybolos/queue/heap"
	"time"
)

type JournalingInmemoryQueue struct {
	InmemoryQueue
	DB   *db.Model
	stop chan bool
}

func NewJournalingInimemoryQueue() Interface {
	model, _ := db.NewModel("./brue.sqlite", "queue")
	tq := JournalingInmemoryQueue{InmemoryQueue{make(DelayedQueue, 0),
		make(map[string]*Task, 0)}, model, make(chan bool)}
	heap.Init(&tq)
	go func(tq *JournalingInmemoryQueue) {
		tq.DB.BatchTransaction()
		ticker := time.NewTicker(1500 * time.Millisecond)
		for _ = range ticker.C {
			select {
			case <-tq.stop:
				ticker.Stop()
				tq.DB.TransactionEnd()
				return
			default:
				tq.DB.BatchTransaction()
			}
		}
	}(&tq)
	return &tq
}

func (tq *JournalingInmemoryQueue) PushTask(task *Task) {
	var wbuff bytes.Buffer
	enc := gob.NewEncoder(&wbuff)
	enc.Encode(task)
	err := tq.DB.Add(task.Id, wbuff.Bytes())
	if err == nil {
		tq.InmemoryQueue.PushTask(task)
	}
}

func (tq *JournalingInmemoryQueue) PopTask() *Task {
	task := tq.InmemoryQueue.PopTask()
	if task.Worker != "" {
		//fmt.Println("Journaling Pop - ", task.Id)
	}

	return task
}

func (tq *JournalingInmemoryQueue) DeleteTask(task *Task) {
	err := tq.DB.Delete(task.Id)
	if err == nil {
		tq.InmemoryQueue.DeleteTask(task)
	}
}

func (tq *JournalingInmemoryQueue) Close() {
	tq.stop <- true
	tq.InmemoryQueue.Close()
}
