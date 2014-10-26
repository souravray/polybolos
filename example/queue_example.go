/*
* @Author: souravray
* @Date:   2014-10-26 20:04:00
* @Last Modified by:   souravray
* @Last Modified time: 2014-10-26 23:33:27
 */

package main

import (
	"fmt"
	q "github.com/souravray/polybolos/queue"
	"net/url"
)

func main() {

	pq := q.NewJournalingInimemoryQueue()

	for i := 0; i < 6; i++ {
		var delay string
		path := fmt.Sprintf("Path%d", i)
		if i%2 == 0 {
			delay = "35s"
		} else {
			delay = fmt.Sprintf("%ds", (25 - 5*i))
		}
		item, _ := q.NewTask(path, url.Values{}, delay)
		pq.PushTask(item)
	}
	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := pq.PopTask()
		if item.Path != "" {
			fmt.Println(item)
		}
	}
}
