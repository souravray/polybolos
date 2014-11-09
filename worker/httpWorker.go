/*
* @Author: souravray
* @Date:   2014-10-27 02:09:33
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-09 22:53:09
 */

package polybolos

import (
	"fmt"
	"net/url"
)

var i int = 0

type HTTPWorkerMethod int

const (
	GET HTTPWorkerMethod = iota
	POST
	PUT
	DELETE
)

type HTTPWorker struct {
	Config
	URI    url.URL
	Method HTTPWorkerMethod
}

func (w *HTTPWorker) Perform(payload url.Values) (err error) {
	i++
	fmt.Println("http worker called ", i, " times")
	return
}
