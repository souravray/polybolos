/*
* @Author: souravray
* @Date:   2014-10-26 22:28:33
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-09 21:53:00
 */

package queue

import (
	"time"
)

type DelayedQueue []*Task

func (dq DelayedQueue) Len() int { return len(dq) }

func (dq DelayedQueue) Less(i, j int) bool {
	return dq[i].priority() > dq[j].priority()
}

func (dq DelayedQueue) Swap(i, j int) {
	dq[i], dq[j] = dq[j], dq[i]
	dq[i].index = i
	dq[j].index = j
}

func (dq *DelayedQueue) Push(x interface{}) {
	n := len(*dq)
	task := x.(*Task)
	task.index = n
	*dq = append(*dq, task)
}

func (dq *DelayedQueue) Pop() interface{} {
	return dq.pop(true)
}

func (dq *DelayedQueue) PopWithoutWait() interface{} {
	return dq.pop(false)
}

func (dq *DelayedQueue) pop(shouldWait bool) interface{} {
	old := *dq
	n := len(old)
	task := old[n-1]
	if shouldWait && time.Now().Before(task.ETA) {
		return new(Task)
	}
	task.index = -1
	*dq = old[0 : n-1]
	return task
}
