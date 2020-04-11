package main

import (
	"fmt"
	"github.com/nwtgck/go-piping-duplex"
	"os"
)

func main() {
	// TODO: hard code
	server := "https://ppng.io"
	selfId := os.Args[1]
	peerId := os.Args[2]
	_, _ = fmt.Fprintf(os.Stderr, "[INFO] Server: %s\n", server)
	_, _ = fmt.Fprintf(os.Stderr, "[INFO] Establishing between '%s' and '%s'...\n", selfId, peerId)
	err := piping_duplex.Wait(server, selfId, peerId)
	if err != nil {
		panic(err)
	}
	_, _ = fmt.Fprintln(os.Stderr, "[INFO] Established!")
	input := os.Stdin
	output := os.Stdout
	err = piping_duplex.Duplex(server, selfId, peerId, input, output)
	if err != nil {
		panic(err)
	}
}
