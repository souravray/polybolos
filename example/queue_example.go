/*
* @Author: souravray
* @Date:   2014-10-26 20:04:00
* @Last Modified by:   souravray
* @Last Modified time: 2014-11-17 02:16:04
 */

package main

import (
	"fmt"
	p "github.com/souravray/polybolos"
	// q "github.com/souravray/polybolos/queue"
	w "github.com/souravray/polybolos/worker"
	"net/url"
	"sync"
	"time"
)

type LazyWorker struct{}

func (l LazyWorker) Perform(payload url.Values) (err error) {
	time.Sleep(time.Second)
	return
}

const longForm = "Jan 2, 2006 at 3:04pm (MST)"

var lock sync.Mutex

func main() {

	pq, err := p.GetQueue(p.INMEMORY_JOURNALING, 500, 500)
	pq.AddHTTPWorker("http-w", url.URL{}, w.GET, 5, 5*time.Minute, time.Second, 4*time.Second, true)
	pq.AddLocalWorker("local-w", LazyWorker{}, 5, 5*time.Minute, time.Second, 4*time.Second, true)
	if err != nil {
		fmt.Println(err)
	}

	go pq.Start()

	for o := 0; o < 10; o++ {
		for m := 0; m < 50; m++ {
			go func(pq *p.Queue) {
				// var task *q.Task
				for i := 0; i < 50; i++ {
					var delay string
					delay = ""
					item := p.NewTask("http-w", url.Values{}, delay, time.Time{})
					pq.PushTask(item)
					// if i%3 == 0 {
					// 	time.Sleep(6 * time.Second)
					// } else if i%5 == 0 {
					// 	task = item
					// } else if i%7 == 0 {
					// 	fmt.Println("call delete on", task)
					// 	pq.DeleteTask(task)
					// }
				}
			}(pq)
		}
	}
	time.Sleep(60 * time.Second)
	pq.Delete()
}
