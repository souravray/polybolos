/*
* @Author: souravray
* @Date:   2015-02-03 01:27:30
* @Last Modified by:   souravray
* @Last Modified time: 2015-02-04 03:05:16
 */

package sys

import "syscall"

func SetFDLimits(newLimit uint64) (uint64, error) {
	var limits syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limits); err != nil {
		return 0, err
	}

	if limits.Cur >= newLimit {
		return limits.Cur, nil
	}

	if newLimit > limits.Max {
		newLimit = limits.Max
	}
	oldLimit := limits.Cur
	limits.Cur = newLimit
	limits.Max = newLimit
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limits); err != nil {
		return oldLimit, err
	}
	return newLimit, nil
}
