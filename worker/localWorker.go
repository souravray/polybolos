/*
* @Author: souravray
* @Date:   2014-11-02 22:33:43
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-10 23:04:29
 */

package worker

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

type LocalInterface interface {
	Perform(payload url.Values, err chan error)
}

type LocalWorker struct {
	Config
	Instance LocalInterface
}

func (w *LocalWorker) Perform(payload url.Values) (err error) {
	c := make(chan error, 1)
	go w.Instance.Perform(payload, c)
	select {
	case err = <-c:
	case <-time.After(w.Timeout):
		err = errors.New("Worker timeout")
	}
	fmt.Println(err, w.Timeout)
	return
}
