/*
* @Author: souravray
* @Date:   2014-11-02 22:19:25
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-03 01:15:03
 */

package worker

import (
	"net/url"
	"time"
)

type Interface interface {
	Perform(payload url.Values) error
}

type Config struct {
	// Maximum time allocated to a worker to complete a job
	// If Timeout is not set expilicitly, then Timeout will same as DefaultWorkerTimeout.
	Timeout time.Duration

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

func (wc Config) SetTimeout(duration string) {
	timeduration, err := time.ParseDuration(duration)
	if err == nil {
		wc.AgeLimit = timeduration
	}
}

func (wc Config) SetRetryLimit(limit int32) {
	if limit > 0 {
		wc.RetryLimit = limit
	} else {
		wc.RetryLimit = 0
	}
}

func (wc Config) SetAgeLimit(duration string) {
	timeduration, err := time.ParseDuration(duration)
	if err == nil {
		wc.AgeLimit = timeduration
	}
}

func (wc Config) SetMinBackoff(duration string) {
	timeduration, err := time.ParseDuration(duration)
	if err == nil {
		wc.MinBackoff = timeduration
	}
}

func (wc Config) SetMaxBackoff(duration string) {
	timeduration, err := time.ParseDuration(duration)
	if err == nil {
		wc.MaxBackoff = timeduration
	}
}

func (wc Config) SetMaxDoubling(flag bool) {
	wc.MaxDoubling = flag
}
