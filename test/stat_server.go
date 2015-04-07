// Package test provides test package for testing gosiege
package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type stats struct {
	hitSec, hitsTotal uint32
}

var s stats
var t *time.Timer

func main() {
	log.SetFlags(0)
	http.HandleFunc("/", testHandler)

	log.Println("Listening on Port : ", 8888)

	t = time.AfterFunc(time.Second*1, printStats)

	log.Fatal(http.ListenAndServe(":8888", nil))

	t.Stop()
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint32(&s.hitSec, 1)
	s.hitsTotal++
}

func printStats() {
	log.Println("Hit/Sec : ", s.hitSec, " Total : ", s.hitsTotal)
	atomic.StoreUint32(&s.hitSec, 0)
	t.Reset(time.Second * 1)
}
