/*
* @Author: souravray
* @Date:   2014-11-02 22:19:25
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-16 03:30:17
 */

package worker

import (
	"net/url"
	"time"
)

type Interface interface {
	Perform(payload url.Values) error
	SetTimeout(duration string)
	SetRetryLimit(limit int32)
	GetRetryLimit() int32
	SetAgeLimit(duration string)
	GetAgeLimit() time.Duration
	SetMinBackoff(duration string)
	SetMaxBackoff(duration string)
	SetMaxDoubling(attempts int32)
	GetInterval(retryAttempts int32) time.Duration
}

type Config struct {
	// Maximum time allocated to a worker to complete a job
	// If Timeout is not set expilicitly, then Timeout will same as DefaultWorkerTimeout.
	Timeout time.Duration

	// Number of tries after which the task fails permanently and is deleted.
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
	MaxDoubling int32
}

func (wc *Config) SetTimeout(duration string) {
	timeduration, err := time.ParseDuration(duration)
	if err == nil {
		wc.Timeout = timeduration
	}
}

func (wc *Config) SetRetryLimit(limit int32) {
	if limit > 0 {
		wc.RetryLimit = limit
	} else {
		wc.RetryLimit = 0
	}
}

func (wc *Config) GetRetryLimit() int32 {
	return wc.RetryLimit
}

func (wc *Config) SetAgeLimit(duration string) {
	timeduration, err := time.ParseDuration(duration)
	if err == nil {
		wc.AgeLimit = timeduration
	}
}

func (wc *Config) GetAgeLimit() time.Duration {
	return wc.AgeLimit
}

func (wc *Config) SetMinBackoff(duration string) {
	timeduration, err := time.ParseDuration(duration)
	if err == nil {
		wc.MinBackoff = timeduration
	}
}

func (wc *Config) SetMaxBackoff(duration string) {
	timeduration, err := time.ParseDuration(duration)
	if err == nil {
		wc.MaxBackoff = timeduration
	}
}

func (wc *Config) SetMaxDoubling(attempts int32) {
	wc.MaxDoubling = attempts
}

func (wc *Config) GetInterval(retryAttempts int32) time.Duration {
	var interval time.Duration
	if wc.MaxDoubling <= 0 {
		interval = wc.MinBackoff
	}

	if wc.MaxDoubling+1 > retryAttempts {
		interval = time.Duration(retryAttempts) * wc.MinBackoff
	} else {
		interval = time.Duration(wc.MaxDoubling+1) * wc.MinBackoff
	}
	return wc.intervalPassFilter(interval)
}

func (wc *Config) intervalPassFilter(interval time.Duration) time.Duration {
	if wc.MaxBackoff == time.Duration(0) ||
		wc.MinBackoff > wc.MaxBackoff {
		return interval
	} else if wc.MaxBackoff > interval {
		return interval
	} else {
		return wc.MaxBackoff
	}

	return time.Duration(0)
}
