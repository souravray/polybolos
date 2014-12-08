/*
* @Author: souravray
* @Date:   2014-10-27 02:09:33
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-09 22:53:09
 */

package polybolos

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var i int = 0

type HTTPWorker struct {
	Config
	URI    url.URL
	Method string
}

func (w *HTTPWorker) Perform(payload url.Values) (err error) {
	responseError := make(chan error)
	go w.request(payload, responseError)
	select {
	case err = <-responseError:
		// return a response
	case <-time.After(time.Second * 20):
		fmt.Println("worker timeout ")
		err = errors.New("Worker time out")
	}
	close(responseError)
	return
}

func (w *HTTPWorker) request(payload url.Values, errC chan error) {
	var postParams url.Values
	var res *http.Response
	var req *http.Request
	var err error
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ResponseHeaderTimeout: 18 * time.Second,
	}
	client := &http.Client{Transport: tr}

	if w.Method == "GET" || w.Method == "DELETE" {
		w.URI.RawQuery = payload.Encode()
	} else {
		postParams = payload
	}

	urlStr := fmt.Sprintf("%v", &w.URI)
	req, err = http.NewRequest(w.Method, urlStr, bytes.NewBufferString(postParams.Encode()))
	if err != nil {
		errC <- err
	}
	res, err = client.Do(req)
	if err != nil {
		errC <- err
	}

	if res == nil {
		errC <- errors.New("Request doesn't return a response")
	}

	fmt.Println("worker status ", res.StatusCode)

	if res.StatusCode > 199 && res.StatusCode <= 299 {
		errC <- nil
	} else if res.StatusCode == 401 {
		errC <- errors.New("Authentication error")
	} else if res.StatusCode > 299 && res.StatusCode < 600 {
		errC <- errors.New("Request returns an error")
	}
}
