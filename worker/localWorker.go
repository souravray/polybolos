/*
* @Author: souravray
* @Date:   2014-11-02 22:33:43
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-09 22:53:23
 */

package polybolos

import (
	"fmt"
	"net/url"
)

type LocalWorker struct {
	Config
	Instance Interface
}

func (w *LocalWorker) Perform(payload url.Values) (err error) {
	fmt.Println("local worker called")
	return
}
