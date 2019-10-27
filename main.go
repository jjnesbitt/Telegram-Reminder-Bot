package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	// body, err := r.GetBody()

	m := map[string]string{}
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &m)

	fmt.Printf("<%s>", (m))
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Printf("Listening!\n")
	http.ListenAndServe(":8080", nil)
}
