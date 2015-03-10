/*
* @Author: souravray
* @Date:   2014-10-15 02:23:23
* @Last Modified by:   souravray
* @Last Modified time: 2015-03-10 07:31:01
 */

package polybolos

import (
	"math"
	"sync/atomic"
	"time"
)

const (
	MIN_PERIOD int32 = 10
)

type bucket struct {
	tokens     int32
	usedTokens int32
	capacity   int32
	period     time.Duration
	fillSize   int32
	stop       chan bool
}

func newBucket(capacity int32, rate int32) (b *bucket, err error) {

	b = &bucket{capacity: capacity}

	if rate > 100 {
		b.fillSize = int32(math.Ceil(float64(rate / MIN_PERIOD)))
		rate = MIN_PERIOD
	} else {
		b.fillSize = 1
	}

	period := time.Duration(1e9 / int64(rate))
	if period > 0 {
		b.period = period
	} else {
		b.period = 1
	}
	return
}

func (b *bucket) setupUsedTokens(delta int32) {
	for {
		usedTokens := atomic.LoadInt32(&b.usedTokens)
		if !atomic.CompareAndSwapInt32(&b.usedTokens, usedTokens, usedTokens+delta) {
			continue
		} else {
			break
		}
	}
	return
}

func (b *bucket) setdownUsedTokens(delta int32) {
	for {
		if usedTokens := atomic.LoadInt32(&b.usedTokens); usedTokens < delta {
			if !atomic.CompareAndSwapInt32(&b.usedTokens, usedTokens, 0) {
				continue
			} else {
				break
			}
		} else {
			if !atomic.CompareAndSwapInt32(&b.usedTokens, usedTokens, usedTokens-delta) {
				continue
			} else {
				break
			}
		}
	}
	return
}

func (b *bucket) getWaitTime(n int32) time.Duration {
	l := math.Ceil(float64(n / b.fillSize))
	return time.Duration(l * float64(b.period.Nanoseconds()))
}

func (b *bucket) Put(n int32) (success bool) {
	for {
		tokens := atomic.LoadInt32(&b.tokens)
		usedTokens := atomic.LoadInt32(&b.usedTokens)
		if tokens+usedTokens < b.capacity {
			if !atomic.CompareAndSwapInt32(&b.tokens, tokens, tokens+n) {
				continue
			} else {
				break
			}
		} else {
			break
		}
	}
	return
}

func (b *bucket) Take(n int32) <-chan int32 {
	c := make(chan int32)
	go func(c chan int32, b *bucket, n int32) {
		var tokens int32
		defer close(c)
	TryToTake:
		for {
			if tokens = atomic.LoadInt32(&b.tokens); tokens == 0 {
				break
			} else if n <= tokens {
				if !atomic.CompareAndSwapInt32(&b.tokens, tokens, tokens-n) {
					continue
				} else {
					break
				}
			} else {
				if !atomic.CompareAndSwapInt32(&b.tokens, tokens, 0) {
					continue
				} else {
					break
				}
			}
		}

		if tokens > 0 {
			b.setupUsedTokens(tokens)
			c <- tokens
		} else {
			time.Sleep(b.getWaitTime(n))
			goto TryToTake
		}
	}(c, b, n)
	return c
}

func (b *bucket) Spend() {
	b.setdownUsedTokens(1)
	return
}

func (b *bucket) Fill() {
	b.stop = make(chan bool, 0)
	go func(b *bucket) {
		defer close(b.stop)
		ticker := time.NewTicker(b.period)
		for _ = range ticker.C {
			select {
			case <-b.stop:
				ticker.Stop()
				return
			default:
				go b.Put(b.fillSize)
			}
		}
	}(b)
}

func (b *bucket) Close() {
	b.stop <- true
	return
}
