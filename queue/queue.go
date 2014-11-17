/*
* @Author: souravray
* @Date:   2014-10-26 17:40:14
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-10 00:22:26
 */
package queue

type Interface interface {
	Len() int
	PushTask(task *Task)
	PopTask() *Task
	DeleteTask(task *Task)
}
