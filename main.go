// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))
var requestLimit = 2 //max allowed requests by requestTimeRange
var requestTimeRange = 5
var sessionExpirancyLimit = 10 //seconds from last request when watcher dies
var AbuseControlRecvChan = make(chan string)
var AbuseControlSendChan = make(chan string)
var AbuseControlRegisterChan = make(chan string)
var AbuseControlUnRegisterChan = make(chan string)
var AbuseControlBlockTime = 1 * 60 // blocked seconds

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if r.FormValue("user_id") == "" {
		http.Error(w, "User_id parameter not found", 405)
		return
	}

	// register worker process by user
	user_id := r.FormValue("user_id")
	if !reg.exists(user_id) {
		worker := watcher{
			id:         user_id,
			hits:       0,
			reqChan:    make(chan string, 0),
			maxRequest: requestLimit,
			autoExpire: getTimer(sessionExpirancyLimit),
		}
		reg.register(user_id, worker.reqChan)
		done := make(chan bool)
		go worker.spawn(done)
		<-done
	}

	workerChan := reg.getWorkerChannel(user_id)
	workerChan <- "request"

	// Request abuse control about user_id
	AbuseControlRecvChan <- user_id
	fmt.Printf("Getting response from abuse control about user %s \n", user_id)
	result := <-AbuseControlSendChan
	fmt.Printf("Result is %s about user %s \n", result, user_id)

	if result == "true" {
		http.Error(w, "User_id is an abuser", 405)
		return
	}

	//if success manipulate header and forward to elbs
	handleResponse(w, r)
}

func handleResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

func main() {
	// initiate abuse register if required
	if !reg.exists("abuse_register") {
		fmt.Printf("Initiate abuse register \n")
		abuseRegister := abuserRegister{
			ids:            make(map[string]bool),
			recvChan:       AbuseControlRecvChan,
			sendChan:       AbuseControlSendChan,
			registerChan:   AbuseControlRegisterChan,
			unRegisterChan: AbuseControlUnRegisterChan,
		}

		done := make(chan bool)
		go abuseRegister.start(done)
		<-done
		fmt.Printf("Initiate abuse register DONE! \n")
	}

	//start http server
	flag.Parse()
	http.HandleFunc("/", serveHome)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
