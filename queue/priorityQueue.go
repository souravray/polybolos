/*
* @Author: souravray
* @Date:   2014-10-26 22:28:33
* @Last Modified by:   souravray
* @Last Modified time: 2014-10-26 22:33:34
 */

package queue

import (
	"time"
)

type PriorityQueue []*Task

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority() > pq[j].priority()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	task := x.(*Task)
	task.index = n
	*pq = append(*pq, task)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	task := old[n-1]
	if time.Now().Before(task.ETA) {
		return new(Task)
	}
	task.index = -1
	*pq = old[0 : n-1]
	return task
}
