package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func proxy(tcpsock net.Conn, wsock *websocket.Conn) {
	for {
		tcpsock.SetReadDeadline(time.Now().Add(3 * time.Second))
		message, _ := bufio.NewReader(tcpsock).ReadString('\n')
		fmt.Print("Message from TCP server: " + message)
		fmt.Printf("%s sent: %s\n", wsock.RemoteAddr(), message)
		err := wsock.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("error")
		}
	}
}

func main() {
	tcpsock, _ := net.Dial("tcp", "127.0.0.1:8081")

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		wsock, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		// Close connection when this function ends
		defer func() {
			fmt.Println("Closing connection...")
			wsock.Close()
		}()

		go proxy(tcpsock, wsock)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":8080", nil)
}
