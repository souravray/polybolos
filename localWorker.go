/*
* @Author: souravray
* @Date:   2014-11-02 22:33:43
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-09 21:41:42
 */

package polybolos

import (
	"fmt"
	"net/url"
)

type LocalWorker struct {
	WorkerConfig
	Instance Worker
}

func (w *LocalWorker) Perform(payload url.Values) (err error) {
	fmt.Println("local worker called")
	return
}
