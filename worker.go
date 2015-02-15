/*
* @Author: souravray
* @Date:   2015-02-16 01:44:19
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-16 03:25:56
 */

package polybolos

import (
	W "github.com/souravray/polybolos/worker"
	"net/url"
	"time"
)

type HTTPWorkerMethod string

const (
	GET    HTTPWorkerMethod = "GET"
	POST   HTTPWorkerMethod = "POST"
	PUT    HTTPWorkerMethod = "PUT"
	DELETE HTTPWorkerMethod = "DELETE"
)

type Worker interface {
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

func NewHTTPWorker(url url.URL, method HTTPWorkerMethod) (worker Worker) {
	worker = &W.HTTPWorker{W.Config{DefaultWorkerTimeout, DefaultRetryLimit, DefaultAgeLimit, DefaultMinBackoff, DefaultMaxBackoff, DefaulrMaxDoubling},
		url,
		string(method)}
	return worker
}

func NewLocalWorker(instance W.LocalInterface) (worker Worker) {
	worker = &W.LocalWorker{W.Config{DefaultWorkerTimeout, DefaultRetryLimit, DefaultAgeLimit, DefaultMinBackoff, DefaultMaxBackoff, DefaulrMaxDoubling},
		instance}
	return worker
}
