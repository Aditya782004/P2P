package main

import (
	"fmt"
	"log"

	"github.com/Aditya-Vaghasiya/foreverstore/p2p"
)

func Onpeer(p2p.Peer) error {
	fmt.Println("Doing some Logic With peer outside the TCPTransport")
	return nil
}

func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddr:    ":4000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        func(p p2p.Peer) error { return fmt.Errorf("error establishing the connection") },
	}

	tr := p2p.NewTCPTransport(opts)
	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Msg: %+v\n", msg)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
