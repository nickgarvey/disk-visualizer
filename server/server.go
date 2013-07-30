package main

import "encoding/json"
import "fmt"
import "os"
import "net/http"

import "code.google.com/p/go.net/websocket"

type clientListeners struct {
	clients map[*websocket.Conn]bool
	add     chan *websocket.Conn
	remove  chan *websocket.Conn
}

var cls = clientListeners{
	clients: make(map[*websocket.Conn]bool),
	add:     make(chan *websocket.Conn),
	remove:  make(chan *websocket.Conn),
}

func main() {
	// Intentionally unbuffered - if we are getting errors _that_ fast
	// then we probably want to slow things down
	errCh := make(chan error)
	go logErrors(errCh)

	traceCh := make(chan blkTrace)
	go sendToClients(traceCh, errCh)
	go traceBlocks(traceCh, errCh)

	go cls.run()

	http.HandleFunc("/", indexHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func wsHandler(ws *websocket.Conn) {
	cls.add <- ws
	defer func() { cls.remove <- ws }()
	for {
		var message [20]byte
		_, err := ws.Read(message[:])
		if err != nil {
			break
		}
	}
}

func (cls *clientListeners) run() {
	for {
		select {
		case ws := <-cls.add:
			cls.clients[ws] = true
		case ws := <-cls.remove:
			delete(cls.clients, ws)
			go ws.Close()
		}
	}
}

func sendToClients(traceCh chan blkTrace, errCh chan error) {
	for trace := range traceCh {
		go func(t blkTrace) {
			if t.Action == "C" && t.Blocks > 0 {
				json, err := json.Marshal(t)
				if err != nil {
					errCh <- err
				} else {
					for conn, _ := range cls.clients {
						conn.Write(json)
					}
				}
			}
		}(trace)
	}
}

func logErrors(ch chan error) {
	for err := range ch {
		fmt.Fprintf(os.Stderr, "ERROR - %s\n", err)
	}
}
