/*
* @Author: souravray
* @Date:   2014-10-26 17:40:14
* @Last Modified by:   souravray
* @Last Modified time: 2014-10-30 23:40:07
 */
package queue

/* There should be a abstuct factory for creating verious queue implemetation
 * For now the implementation cannot be done because of  this error https://groups.google.com/forum/#!topic/golang-nuts/-ZoCu5m0kJ4
 */
type Queue interface {
	Len() int
	PushTask(task *Task)
	PopTask() *Task
	DeleteTask(task *Task)
}
