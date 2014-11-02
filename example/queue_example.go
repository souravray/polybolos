/*
* @Author: souravray
* @Date:   2014-10-26 20:04:00
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-02 21:00:54
 */

package main

import (
	"fmt"
	p "github.com/souravray/polybolos"
	q "github.com/souravray/polybolos/queue"
	"net/url"
	"sync"
	"time"
)

const longForm = "Jan 2, 2006 at 3:04pm (MST)"

var lock sync.Mutex

func main() {

	pq, err := p.GetQueue(p.INMEMORY_JOURNALING, 15, 11)
	if err != nil {
		fmt.Println(err)
	}

	go func(pq *p.Queue) {
		for i := 0; i < 33; i++ {
			item := p.NewTask("dummyA", url.Values{}, "7s", time.Time{})
			pq.PushTask(item)
			if i%5 == 0 {
				time.Sleep(3 * time.Second)
			}
		}
	}(pq)
	go func(pq *p.Queue) {
		for i := 0; i < 88; i++ {
			item := p.NewTask("dummyD", url.Values{}, "22s", time.Time{})
			pq.PushTask(item)
		}
	}(pq)

	go func(pq *p.Queue) {
		for i := 0; i < 8; i++ {
			eta, err := time.Parse(longForm, "Nov 3, 2014 at 7:54pm (IST)")
			fmt.Println(eta, "   ", time.RFC3339, " ", err)
			item := p.NewTask("dummyE", url.Values{}, "1h", eta)
			pq.PushTask(item)
		}
	}(pq)

	go func(pq *p.Queue) {
		var task *q.Task
		for i := 0; i < 18; i++ {
			var delay string
			path := fmt.Sprintf("Path%d", i)
			if i%2 == 0 {
				delay = "2s"
			} else {
				delay = "21s"
			}
			item := p.NewTask(path, url.Values{}, delay, time.Time{})
			pq.PushTask(item)
			if i%3 == 0 {
				time.Sleep(6 * time.Second)
			} else if i%5 == 0 {
				task = item
			} else if i%7 == 0 {
				fmt.Println("call delete on", task)
				pq.DeleteTask(task)
			}
		}
	}(pq)

	go func(pq *p.Queue) {
		for i := 0; i < 44; i++ {
			item := p.NewTask("dummyB", url.Values{}, "6s", time.Time{})
			pq.PushTask(item)
			if i%7 == 0 {
				time.Sleep(3 * time.Second)
			}
		}
	}(pq)

	go func(pq *p.Queue) {
		for i := 0; i < 60; i++ {
			item := p.NewTask("dummyC", url.Values{}, "5s", time.Time{})
			pq.PushTask(item)
			if i%5 == 0 {
				time.Sleep(2 * time.Second)
			}
		}
	}(pq)

	//Take the items out
	// they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := pq.PopTask()
		if item.Path != "" {
			fmt.Println(item)
		}
	}
	go pq.Start()
	time.Sleep(1 * time.Minute)
	pq.Delete()
}
