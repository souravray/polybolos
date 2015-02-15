/*
* @Author: souravray
* @Date:   2014-10-11 19:51:37
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-15 20:12:20
 */

package queue

import (
	"github.com/nu7hatch/gouuid"
	"net/url"
	"time"
)

type Task struct {
	// Id is rfc4122 UUID for the Task
	Id string
	// Worker Name for the task.
	Worker string

	// Payload is the data for the task.
	// This will be delivered as the HTTP request body in case of POST or PUT
	// This will be converted to Querystring incase of GET or DELETE
	Payload url.Values

	// Minimum Delay specifies the duration the task queue service atleast wait
	// before executing the task for the first time.
	// Either Delay or ETA may be set, if both are set then MinDelay will be ignored.
	MinDelay time.Duration

	// ETA specifies the earliest time a task may be executed
	ETA time.Time

	// The number of times the task has been dispatched
	RetryCount int32

	//Time when task is equed to queue for the first time
	EnqueTime time.Time
	// index refers to the queue position
	index int
}

func (task *Task) priority() int32 {
	eta := task.ETA.Unix()
	now := time.Now().Unix()
	return int32(now - eta)
}

func NewTask(worker string, payload url.Values, delay string, eta time.Time) (task *Task, err error) {
	var uid *uuid.UUID
	uid, err = uuid.NewV4()
	if err != nil {
		return nil, err
	}

	task = &Task{
		Id:         uid.String(),
		EnqueTime:  time.Now(),
		Worker:     worker,
		Payload:    payload,
		RetryCount: 0,
	}

	if !eta.IsZero() {
		task.ETA = eta
	} else if delay != "" {
		task.MinDelay, err = time.ParseDuration(delay)
		if err != nil {
			return
		}
		task.ETA = time.Now().Add(task.MinDelay)
	} else {
		task.ETA = time.Now()
	}
	return
}
