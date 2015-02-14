/*
* @Author: souravray
* @Date:   2015-02-14 23:40:14
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-14 23:50:57
 */

package sys

import (
	"os"
	"runtime"
)

func UseMaxCPUs() {
	if envGOMAXPROCS := os.Getenv("GOMAXPROCS"); envGOMAXPROCS != "" {
		return
	}
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
}
