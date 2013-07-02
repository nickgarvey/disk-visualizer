package main

import "encoding/json"
import "fmt"
import "os"
import "time"

func main() {
	// Intentionally unbuffered - if we are getting errors _that_ fast
	// then we probably want to slow things down
	errCh := make(chan error)
	go logErrors(errCh)

	traceCh := make(chan blkTrace)
	go sendToClients(traceCh, errCh)

	go traceBlocks(traceCh, errCh)
	for {
		time.Sleep(1 * time.Second)
	}

}

func sendToClients(traceCh chan blkTrace, errCh chan error) {
	for trace := range traceCh {
		go func(t blkTrace) {
			json, err := json.Marshal(t)
			if err != nil {
				errCh <- err
			} else {
				fmt.Println(string(json))
			}
		}(trace)
	}
}

func logErrors(ch chan error) {
	for err := range ch {
		fmt.Fprintf(os.Stderr, "ERROR - %s\n", err)
	}
}
