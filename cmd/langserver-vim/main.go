package main

import (
	"context"
	"log"
	"os"

	"github.com/haya14busa/go-vimlparser/langserver"
	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	log.Println("langserver-vim: reading on stdin, writing on stdout")
	var connOpt []jsonrpc2.ConnOpt
	<-jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}), langserver.NewHandler(), connOpt...).DisconnectNotify()
	log.Println("langserver-vim: connections closed")
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
