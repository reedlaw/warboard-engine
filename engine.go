package main

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
	"os"
	// "github.com/araddon/httpstream"
)

func main() {
	http.Handle("/js/", http.FileServer(http.Dir(".")))
	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/websocket/", websocket.Handler(wsHandler))
	err := http.ListenAndServe("localhost:8000", nil)
	checkError(err)
}

func wsHandler(ws *websocket.Conn) {
	fmt.Println("Processing websockets")
	go readKeyboard(ws)
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

func readKeyboard(ws *websocket.Conn) {
	for {
		in := bufio.NewReader(os.Stdin)
		input, err := in.ReadString('\n')
		if err != nil {
			// handle error
		}
		websocket.Message.Send(ws, input)
		fmt.Println("Keyboard message", input)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
