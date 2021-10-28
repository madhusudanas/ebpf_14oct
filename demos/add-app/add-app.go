package main

import (
	"fmt"
	"net/http"
	"strconv"
)

//go:noinline
func add(n uint64, k uint64, l uint64, m uint64, o uint64) uint64 {
  return n + k + l + m + o
}

func main() {
	addr := ":8000"
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		nth := uint64(1)
		n, ok := r.URL.Query()["n"]
		if ok && len(n[0]) >= 1 {
			val, err := strconv.ParseUint(n[0], 10, 64)
			if err != nil || val <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			nth = val
		}

		w.Write([]byte(fmt.Sprintf("Sum = %d\n", add(nth, nth*2, nth*3, nth*4, nth*5))))
	})

	fmt.Printf("Starting server on: %+v\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to run http server: %v\n", err)
	}
}
