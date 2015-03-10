/*
* @Author: souravray
* @Date:   2014-10-26 20:52:28
* @Last Modified by:   souravray
* @Last Modified time: 2015-03-10 08:26:45
 */

package queue

import (
	"bytes"
	"encoding/gob"
	"github.com/souravray/polybolos/queue/db"
	"github.com/souravray/polybolos/queue/heap"
	"path"
	"time"
)

const (
	JOURNAL_COLLECTION     string        = "queue"
	JOURNAL_WRITE_INTERVAL time.Duration = 1500 * time.Millisecond
)

type JournalingInmemoryQueue struct {
	InmemoryQueue
	DB   *db.Model
	stop chan bool
}

func NewJournalingInimemoryQueue(journalPath, queueName string) Interface {
	dbPath := path.Join(journalPath, queueName)
	model, _ := db.NewModel(dbPath, "queue")

	tq := JournalingInmemoryQueue{
		InmemoryQueue{DelayedQueue: make(DelayedQueue, 0)},
		model,
		make(chan bool),
	}
	heap.Init(&tq)

	tq.DB.BatchTransaction()
	tq.recover()
	go func(tq *JournalingInmemoryQueue) {
		ticker := time.NewTicker(JOURNAL_WRITE_INTERVAL)
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

func (tq *JournalingInmemoryQueue) recover() {
	pendingTask := tq.DB.Count()
	limit := 1000
	p := 0
	for offset := 0; offset < pendingTask; offset = offset + limit {
		results := tq.DB.Read(offset, limit)
		for result := range results {
			p++
			rbuff := bytes.NewBuffer(result)
			enc := gob.NewDecoder(rbuff)
			task := Task{}
			enc.Decode(&task)
			if task.IsEmpty() == false {
				tq.InmemoryQueue.PushTask(&task)
			}
		}
	}
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
	return task
}

func (tq *JournalingInmemoryQueue) DeleteTask(task *Task) {
	err := tq.DB.Delete(task.Id)
	if err == nil {
		tq.InmemoryQueue.DeleteTask(task)
	}
}

func (tq *JournalingInmemoryQueue) CleanTask(task *Task) {
	tq.InmemoryQueue.CleanTask(task)
	tq.DB.Delete(task.Id)
}

func (tq *JournalingInmemoryQueue) Close() {
	tq.stop <- true
	tq.InmemoryQueue.Close()
}
