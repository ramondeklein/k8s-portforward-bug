package main

import (
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

func main() {
	i := int32(0)
	err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqNr := atomic.AddInt32(&i, 1)
		log.Printf("%d: got request from %s", reqNr, req.RemoteAddr)
		size := 32 * 1024 * 1024
		w.Header().Add("Content-Type", "application/data")
		w.Header().Add("Content-Length", strconv.Itoa(size))

		buf := make([]byte, 64*1024)
		for size > 0 {
			maxBytes := size
			if maxBytes > len(buf) {
				maxBytes = len(buf)
			}
			n, err := w.Write(buf[:maxBytes])
			if err != nil {
				log.Printf("%d: error writing data (%d bytes left to write): %s", reqNr, size, err)
				return
			}
			size = size - n
		}
		log.Printf("%d: completed request", reqNr)
	}))
	if err != nil {
		log.Fatalln("Error listening: " + err.Error())
	}
}
