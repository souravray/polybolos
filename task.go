/*
* @Author: souravray
* @Date:   2015-02-16 01:59:19
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-16 03:55:50
 */

package polybolos

import (
	Q "github.com/souravray/polybolos/queue"
	"net/url"
	"time"
)

type Task interface {
	SetDelay(delay string)
	SetETA(eta time.Time)
	IsEmpty() bool
}

func NewTask(path string, payload url.Values) (task Task) {
	task, _ = Q.NewTask(path, payload)
	return task
}
