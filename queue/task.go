/*
* @Author: souravray
* @Date:   2014-10-11 19:51:37
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-02 20:58:36
 */

package queue

import (
	"net/url"
	"time"
)

type Task struct {
	// Path is the worker URL for the task.
	Path string

	// Payload is the data for the task.
	// This will be delivered as the HTTP request body in case of POST or PUT
	// This will be converted to Querystring incase of GET or DELETE
	Payload url.Values

	// Minimum Delay specifies the duration the task queue service atleast wait
	// before executing the task.
	// Either Delay or ETA may be set, if both are set then MinDelay will be ignored.
	MinDelay time.Duration

	// ETA specifies the earliest time a task may be executed
	ETA time.Time

	// The number of times the task has been dispatched
	RetryCount int

	// private properties
	//Time when task is equed to queue for the first time
	enqueTS int64
	// index refers to the queue position
	index int
}

func (task *Task) priority() int32 {
	eta := task.ETA.Unix()
	now := time.Now().Unix()
	return int32(now - eta)
}

func NewTask(path string, payload url.Values, delay string, eta time.Time) (task *Task, err error) {
	task = new(Task)
	task.enqueTS = time.Now().Unix()
	task.Path = path
	task.Payload = payload
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
