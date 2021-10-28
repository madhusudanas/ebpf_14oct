package main

import (
	"fmt"
	"net/http"
)

//go:noinline
func concat(n string, k string, l string, m string, o string) string {
  return n + k + l + m + o
}

func main() {
	addr := ":8000"
	http.HandleFunc("/concat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
                var n string = "default"
		t, ok := r.URL.Query()["n"]
		if ok && len(t[0]) >= 1 {
			n = t[0]
		}
                var a string = "a"
                var b string = "b"
                var c string = "c"
                var d string = "d"
		w.Write([]byte(fmt.Sprintf("Concat = %s\n", concat(n, a, b, c, d))))
	})

	fmt.Printf("Starting server on: %+v\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to run http server: %v\n", err)
	}
}
