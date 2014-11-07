package main

import "fmt"

type abuserRegister struct {
	ids            map[string]bool
	recvChan       chan string
	sendChan       chan string
	registerChan   chan string
	unRegisterChan chan string
}

func (a *abuserRegister) start(done chan bool) {
	done <- true
	for {
		select {
		case user_id := <-a.recvChan:
			if _, exists := a.ids[user_id]; exists {
				fmt.Printf("UserId %s Is abuser\n", user_id)
				a.sendChan <- "true"
			} else {
				fmt.Printf("UserId %s Is Not abuser \n", user_id)
				a.sendChan <- "false"
			}
		case user_id := <-a.registerChan:
			if _, exists := a.ids[user_id]; !exists {
				a.ids[user_id] = true
				abuser := abuser{
					id:         user_id,
					blockTimer: getTimer(AbuseControlBlockTime),
				}
				go abuser.spawn()
			}
		case user_id := <-a.unRegisterChan:
			fmt.Printf("Unblocking user %s \n", user_id)
			delete(a.ids, user_id)
		}
	}
}
