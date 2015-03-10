/*
* @Author: souravray
* @Date:   2014-10-12 01:31:40
* @Last Modified by:   souravray
* @Last Modified time: 2015-03-02 07:30:44
 */

package polybolos

import (
	"time"
)

// default configuration for Worker
var DefaultRetryLimit int32 = 3

var DefaultWorkerTimeout = 90 * time.Second
var DefaultAgeLimit time.Duration = 0 * time.Second
var DefaultMinBackoff time.Duration = 1 * time.Second
var DefaultMaxBackoff time.Duration = 6 * time.Second

var DefaulrMaxDoubling int32 = 0

// default configuration for Bucket
var DefaulrMaxConcurrentWorker int32 = 10
var DefaulrMaxDequeueRate int32 = 2

// constants
