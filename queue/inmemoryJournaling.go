/*
* @Author: souravray
* @Date:   2014-10-26 20:52:28
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-08 09:21:29
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
	tq := JournalingInmemoryQueue{InmemoryQueue{DelayedQueue: make(DelayedQueue, 0),
		OneTheflyQueue: make(map[string]*Task, 0)}, model, make(chan bool)}
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

	var err error
	if task.RetryCount == 0 {
		err = tq.DB.Add(task.Id, wbuff.Bytes())
	} else {
		err = tq.DB.Update(task.Id, wbuff.Bytes())
	}

	if err == nil {
		tq.InmemoryQueue.PushTask(task)
	}
}

func (tq *JournalingInmemoryQueue) PopTask() *Task {
	task := tq.InmemoryQueue.PopTask()
	if task.Worker != "" {
		tq.InmemoryQueue.DeleteTask(task)
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
