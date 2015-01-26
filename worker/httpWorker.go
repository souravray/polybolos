/*
* @Author: souravray
* @Date:   2014-10-27 02:09:33
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-23 14:54:36
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

func (w *HTTPWorker) Perform(payload url.Values) error {
	err := w.request(payload)
	/* Droping support to hard time out
	for httpWorkers, because it is causing
	race contions*/
	// select {
	// case err := <-errC:
	// 	return err
	// case <-time.After(time.Second * 20):
	// 	err := errors.New("Worker time out")
	// 	return err
	// }
	return err
}

func (w *HTTPWorker) request(payload url.Values) (err error) {
	var postParams url.Values
	var res *http.Response
	var req *http.Request
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ResponseHeaderTimeout: 18 * time.Second,
		// for now we are relying on ResponseHeaderTimeout
		// wich is erronus, because it is not hard timeout
		// function. It is a problem when network connection
		// is not available, in that case it will take
		// additional 30s for dialer failure.
		// solution: we need a custom implementation
		// of dialer and transport layer
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
		return err
	}
	res, err = client.Do(req)
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("Request doesn't return a response")
	}

	fmt.Println("worker status ", res.StatusCode)

	if res.StatusCode > 199 && res.StatusCode <= 299 {
		return nil
	} else if res.StatusCode == 401 {
		return errors.New("Authentication error")
	} else if res.StatusCode > 299 && res.StatusCode < 600 {
		return errors.New("Request returns an error")
	}

	return errors.New("Unhandeled reference code response")
}
