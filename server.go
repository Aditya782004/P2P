package main

import (
	"fmt"
	"log"

	"github.com/Aditya-Vaghasiya/foreverstore/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransFormFunc PathTransFormFunc
	Transport         p2p.Transport
	Bootstrapnodes    []string
}

type FileServer struct {
	FileServerOpts
	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	StoreOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransFormFunc: opts.PathTransFormFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(StoreOpts),
		quitch:         make(chan struct{}),
	}
}

func (s *FileServer) loop() {
	defer func() {
		fmt.Println("file server stopped due to usser quit action")
	}()
	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Printf("msg: %s", msg.Payload)
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.Bootstrapnodes {
		go func(addr string) {
			fmt.Println("attempting to connect with remote hosts")
			err := s.Transport.Dial(addr)
			if err != nil {
				log.Println("dial error: ", err)
			}

		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstrapNetwork()
	s.loop()
	return nil
}
