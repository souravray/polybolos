/*
* @Author: souravray
* @Date:   2014-10-26 17:40:14
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-02 22:34:37
 */
package queue

type Queue interface {
	Len() int
	PushTask(task *Task)
	PopTask() *Task
	DeleteTask(task *Task)
}
