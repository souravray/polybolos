/*
* @Author: souravray
* @Date:   2014-10-26 17:40:14
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-17 09:00:13
 */
package queue

type Interface interface {
	PushTask(task *Task)
	PopTask() *Task
	DeleteTask(task *Task)
	Close()
}
