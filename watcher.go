package main

import (
	"fmt"
	"time"
)

type watcher struct {
	id         string
	hits       int
	reqChan    chan string
	maxRequest int
	autoExpire *time.Timer
}

func (w *watcher) spawn(done chan bool) {
	limiter := time.Tick(time.Second * time.Duration(requestTimeRange))
	fmt.Printf("Spawnning id:%s \n", w.id)
	done <- true
	for {
		select {
		case msg := <-w.reqChan:
			w.hits++
			fmt.Printf("received: %s on worker %s Hits %d \n", msg, w.id, w.hits)
			//update auto expire timer
			w.autoExpire = getTimer(sessionExpirancyLimit)
		case <-limiter:
			fmt.Printf("Expired on %s total hits %d \n", w.id, w.hits)
			if w.hits > w.maxRequest {
				fmt.Printf("Max Requests exceeded \n")
				AbuseControlRegisterChan <- w.id
			}
			w.hits = 0
		case <-w.autoExpire.C:
			fmt.Printf("Watcher from user %s die! \n", w.id)
			reg.unregister(w.id)

			return
		}
	}
}

func getTimer(expirancy int) *time.Timer {
	return time.NewTimer(time.Second * time.Duration(expirancy))
}
