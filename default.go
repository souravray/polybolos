/*
* @Author: souravray
* @Date:   2014-10-12 01:31:40
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-03 00:34:25
 */

package polybolos

import (
	"time"
)

// default configuration for Worker
var DefaultRetryLimit int32 = 3

var DefaultWorkerTimeout = 30 * time.Second
var DefaultAgeLimit time.Duration = 0 * time.Second
var DefaultMinBackoff time.Duration = 0 * time.Second
var DefaultMaxBackoff time.Duration = 0 * time.Second

var DefaulrMaxDoubling bool = false

// default configuration for Bucket
var DefaulrMaxConcurrentWorker int32 = 10
var DefaulrMaxDequeueRate int32 = 5

// constants
type HTTPWorkerMethod string

const (
	GET    HTTPWorkerMethod = "GET"
	POST   HTTPWorkerMethod = "POST"
	PUT    HTTPWorkerMethod = "PUT"
	DELETE HTTPWorkerMethod = "DELETE"
)
