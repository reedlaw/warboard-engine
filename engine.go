package main

import (
	"fmt"
	"os"
	"net/http"
	"code.google.com/p/go.net/websocket"
	// "github.com/araddon/httpstream"
)

func main() {
	http.Handle("/js/", http.FileServer(http.Dir(".")))
	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/websocket/", websocket.Handler(socketHandler))
	err := http.ListenAndServe("localhost:8000", nil)
	checkError(err)
}

func socketHandler(ws *websocket.Conn) {
	fmt.Println("Processing websockets")
	// go readWebsocket(ws)
	var msg string

	for {
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			fmt.Println("Websocket error", err)
			break
		}

		fmt.Println("Websocket message", msg)
	}
	fmt.Println("Exit")
}

func readWebsocket(ws *websocket.Conn) {
	for {
		// do something
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
