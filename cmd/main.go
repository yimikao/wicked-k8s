package main

import (
	"fmt"
	"net/http"
)

func main() {

	m := http.NewServeMux()
	m.HandleFunc("/f", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("first"))
	})

	m2 := http.NewServeMux()
	m2.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("second"))
	})

	fmt.Println("Hakuna Matata!")
	go func() {
		http.ListenAndServe(":8080", m)
	}()
	http.ListenAndServe(":8081", m2)
}
