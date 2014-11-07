package main

import (
	"fmt"
	"time"
)

type abuser struct {
	id         string
	blockTimer *time.Timer
}

func (a *abuser) spawn() {
	fmt.Printf("Spawnning Abuser on user id:%s \n", a.id)
	for {
		select {
		case <-a.blockTimer.C:
			fmt.Printf("Abuse control Block time from user %s expired, unblock! \n", a.id)
			AbuseControlUnRegisterChan <- a.id

			return
		}
	}
}
