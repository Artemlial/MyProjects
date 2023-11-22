package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	msg string `json:"msg"`
}

var port *string

func init() {
	port = flag.String("port", "9000", "specify port")
}

func connect() (*websocket.Conn, error) {
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", randIP())
}

func randIP() string {
	var arr [4]int64
	for i := 0; i < 4; i++ {
		source := rand.NewSource(time.Now().UnixNano())
		arr[i] = source.Int63()
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}

func main() {
	flag.Parse()
	ws, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	var m Message
	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				log.Println("could not read message ", err.Error())
				break
			}
			fmt.Println(m)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		m = Message{msg: text}
		err := websocket.JSON.Send(ws, m)
		if err != nil {
			log.Println("could not send message ", err.Error())
			break
		}
	}
}
