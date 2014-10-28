/*
* @Author: souravray
* @Date:   2014-10-15 02:23:23
* @Last Modified by:   souravray
* @Last Modified time: 2014-10-29 01:27:22
 */

package polybolos

import (
	"sync"
	"time"
)

type Bucket struct {
	tokens     int32
	usedTokens int32
	capacity   int32
	freq       time.Duration
	mutex      sync.Mutex
	stop       chan bool
}

func NewBucket(capacity int32, rate int32) (b *Bucket, err error) {

	b = &Bucket{capacity: capacity}
	minFreq := time.Duration(1e9 / int64(capacity))
	freq := time.Duration(1e9 / int64(rate))
	if minFreq > freq {
		b.freq = minFreq
	} else {
		b.freq = freq
	}
	return
}

func (b *Bucket) Put() (success bool) {

	if b.tokens+b.usedTokens < b.capacity {
		b.mutex.Lock()
		b.tokens++
		b.mutex.Unlock()
	}
	return
}

func (b *Bucket) Take(n int32) <-chan int32 {

	c := make(chan int32)
	go func(c chan int32, b *Bucket, n int32) {
		// waiting loop
		// @todo: a better solution than waiting for loop
		for b.tokens == 0 {
			time.Sleep(time.Duration(n * int32(b.freq.Nanoseconds())))
		}
		b.mutex.Lock()
		if n > b.tokens {
			n = b.tokens
		}
		b.tokens -= n
		b.usedTokens += n
		b.mutex.Unlock()
		c <- n
		defer close(c)
		return
	}(c, b, n)
	return c
}

func (b *Bucket) Spend() (success bool) {

	b.mutex.Lock()
	if b.usedTokens > 0 {
		b.usedTokens--
	}
	b.mutex.Unlock()
	return
}

func (b *Bucket) Fill() {

	b.stop = make(chan bool)
	ticker := time.NewTicker(b.freq)
	for _ = range ticker.C {
		select {
		case <-b.stop:
			ticker.Stop()
			return
		default:
			b.Put()
		}
	}
}

func (b *Bucket) Close() {

	b.stop <- true
	defer close(b.stop)
	return
}
