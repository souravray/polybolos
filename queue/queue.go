/*
* @Author: souravray
* @Date:   2014-10-26 17:40:14
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-08 23:07:18
 */
package queue

type Interface interface {
	PushTask(task *Task)
	PopTask() *Task
	CleanTask(task *Task)
	DeleteTask(task *Task)
	Close()
}
