/*
* @Author: souravray
* @Date:   2014-11-02 22:33:43
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-03 01:15:17
 */

package worker

import (
	"fmt"
	"net/url"
)

// type localInterface interface {
// 	Perform(errC chan err, payload url.Values) error
// }

type LocalWorker struct {
	Config
	Instance Interface
}

func (w *LocalWorker) Perform(payload url.Values) (err error) {
	//errC := make(chan error, 1)
	fmt.Println("local worker called")
	return
}
