/*
* @Author: souravray
* @Date:   2014-10-27 02:09:33
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-10 00:19:44
 */

package worker

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var i int = 0

type HTTPWorker struct {
	Config
	URI    url.URL
	Method string
}

func (w *HTTPWorker) Perform(payload url.Values) error {
	var postParams url.Values
	var res *http.Response
	var req *http.Request
	var err error

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		// Timeout is supported from Go version 1.3 onward
		Timeout: w.Timeout,
	}

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

	if res.StatusCode > 199 && res.StatusCode <= 299 {
		return nil
	} else if res.StatusCode == 401 {
		return errors.New("Authentication error")
	} else if res.StatusCode > 299 && res.StatusCode < 600 {
		return errors.New("Request returns an error")
	}

	return errors.New("Unhandeled reference code response")
}
