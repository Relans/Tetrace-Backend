package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"weboscket_dispatcher"
)

func main() {
	http.HandleFunc("/ws", weboscket_dispatcher.WsHandler)
	http.HandleFunc("/", rootHandler)

	panic(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}