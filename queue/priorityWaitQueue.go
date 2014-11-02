/*
* @Author: souravray
* @Date:   2014-10-26 22:28:33
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-02 20:57:53
 */

package queue

import (
	"time"
)

type PriorityWaitQueue []*Task

func (pwq PriorityWaitQueue) Len() int { return len(pwq) }

func (pwq PriorityWaitQueue) Less(i, j int) bool {
	return pwq[i].priority() > pwq[j].priority()
}

func (pwq PriorityWaitQueue) Swap(i, j int) {
	pwq[i], pwq[j] = pwq[j], pwq[i]
	pwq[i].index = i
	pwq[j].index = j
}

func (pwq *PriorityWaitQueue) Push(x interface{}) {
	n := len(*pwq)
	task := x.(*Task)
	task.index = n
	*pwq = append(*pwq, task)
}

func (pwq *PriorityWaitQueue) Pop() interface{} {
	return pwq.pop(true)
}

func (pwq *PriorityWaitQueue) PopWithoutWait() interface{} {
	return pwq.pop(false)
}

func (pwq *PriorityWaitQueue) pop(shouldWait bool) interface{} {
	old := *pwq
	n := len(old)
	task := old[n-1]
	if shouldWait && time.Now().Before(task.ETA) {
		return new(Task)
	}
	task.index = -1
	*pwq = old[0 : n-1]
	return task
}
