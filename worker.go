/*
* @Author: souravray
* @Date:   2014-11-02 22:19:25
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-09 21:41:05
 */

package polybolos

import (
	"net/url"
	"time"
)

type Worker interface {
	Perform(payload url.Values) error
}

type WorkerConfig struct {
	// Number of tries/leases after which the task fails permanently and is deleted.
	// If AgeLimit is also set, both limits must be exceeded for the task to fail permanently.
	RetryLimit int32

	// Maximum time allowed since the task's first try before the task fails permanently and is deleted
	// If RetryLimit is also set, both limits must be exceeded for the task to fail permanently.
	AgeLimit time.Duration

	// Minimum time between successive tries
	MinBackoff time.Duration

	// Maximum time between successive tries
	MaxBackoff time.Duration

	// Maximum number of times to double the interval between successive tries before the intervals increase linearly
	MaxDoubling bool
}
