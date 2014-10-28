/*
* @Author: souravray
* @Date:   2014-10-27 02:09:33
* @Last Modified by:   souravray
* @Last Modified time: 2014-10-29 01:55:37
 */

package polybolos

import (
	"net/url"
	"time"
)

type WorkerMethod int

const (
	GET WorkerMethod = iota
	POST
	PUT
	DELETE
)

type Worker struct {
	url.URL

	Method WorkerMethod
	// Number of tries/leases after which the task fails permanently and is deleted.
	// If AgeLimit is also set, both limits must be exceeded for the task to fail permanently.
	RetryLimit int32

	// Maximum time allowed since the task's first try before the task fails permanently and is deleted (only for push tasks).
	// If RetryLimit is also set, both limits must be exceeded for the task to fail permanently.
	AgeLimit time.Duration

	// Minimum time between successive tries (only for push tasks).
	MinBackoff time.Duration

	// Maximum time between successive tries (only for push tasks).
	MaxBackoff time.Duration

	// Maximum number of times to double the interval between successive tries before the intervals increase linearly (only for push tasks).
	MaxDoublings int32
}

// func NewWorker(host string, endpoint string) *Worker {

// }

// func (worker Worker) Name(name string) {

// }
