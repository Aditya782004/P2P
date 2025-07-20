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

func OnPeer(p p2p.Peer) error {
	return nil
}

func makeServer(ListenAddr string, nodes ...string) *FileServer {
	TcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ListenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		//TODO OnPeer:        OnPeer,
	}
	tcpTransport := p2p.NewTCPTransport(TcptransportOpts)

	fileServer := FileServerOpts{
		StorageRoot:       "4000_practice",
		PathTransFormFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		Bootstrapnodes:    nodes,
	}
	return NewFileServer(fileServer)
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", "")

	go func() {
		err := s1.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()
	s2.Start()

}
