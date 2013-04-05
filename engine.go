package main

import (
  "bufio"
  "code.google.com/p/go.net/websocket"
	// "encoding/json"
  "flag"
  "fmt"
  "github.com/araddon/httpstream"
  "log"
  "net/http"
  "os"
	"time"
)

var (
  pwd      *string = flag.String("pwd", "password", "Password")
  user     *string = flag.String("user", "username", "username")
  track    *string = flag.String("track", "", "Twitter terms to track")
  logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
)

func main() {

  flag.Parse()

  if *user == "username" || *pwd == "Password" {
    fmt.Println("To pull Twitter data do: go run engine.go -user=twitter_username -pwd=password");
  } else {
		go fetchTwitter()
	}
		
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
	go ticker(ws)
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

func fetchTwitter() {
  httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), *logLevel)
  stream := make(chan []byte, 1000)
  done := make(chan bool)

  client := httpstream.NewBasicAuthClient(*user, *pwd, httpstream.OnlyTweetsFilter(func(line []byte) {
    stream <- line
  }))
	
	err := client.Sample(done)
	if err != nil {
		httpstream.Log(httpstream.ERROR, err.Error())
	} else {
		
		go func() {
			ct := 0
			for tw := range stream {
				println(string(tw))
				// heavy lifting
				ct++
				if ct > 1 {
					done <- true
				}
			}
		}()
		_ = <-done
	}
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

func ticker(ws *websocket.Conn) {
	for {
    time.Sleep(1 * 1e9 )
		t := time.Now()
		var jsonBlob = "{\"event\":\"time\",\"data\":{\"name\":\"warboard\",\"message\":\"" +	t.Format("20060102150405") + "\"}}"
		websocket.Message.Send(ws, jsonBlob)
	}
}

func checkError(err error) {
  if err != nil {
    fmt.Println("Fatal error ", err.Error())
    os.Exit(1)
  }
}
