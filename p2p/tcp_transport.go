package p2p

import (
	"fmt"
	"log"
	"net"
	"time"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcchan  chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		rpcchan:          make(chan RPC),
		TCPTransportOpts: opts,
	}
}

// Consume Implements the Transport Interface which will return the read-only channel
// for Reading the readonly incoming message from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcchan
}

func (t *TCPTransport) NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close Implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// Close Implements Transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial Implements the Transport Interfac.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAccpetLoop()
	log.Printf("tcp transport listening on port 4000")
	return nil

}

func (t *TCPTransport) startAccpetLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Acccept Error %s\n", err)
			fmt.Println(conn)
		}
		go t.handleConn(conn, false)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	defer func() {
		fmt.Printf("connections closed %+v\n", err)
		conn.Close()
	}()
	peer := t.NewTCPPeer(conn, outbound)
	// fmt.Printf("new incoming connections %+v\n", peer)
	if err = t.HandshakeFunc(peer); err != nil {
		return
	}
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	//buf := make([]byte, 20000)
	rpc := RPC{}
	for {
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Printf("error while reading from the connection: %+v\n", err)
		// }

		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			logConnectionError(err)
			// fmt.Printf("error while reading from the connection: %+v\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcchan <- rpc
		// fmt.Printf("Msg: %+v\n", *rpc)
	}

}

var lastLogTime time.Time
var logThrottle = 3 * time.Second

func logConnectionError(err error) {
	if time.Since(lastLogTime) > logThrottle {
		log.Printf("connection error: %v", err)
		lastLogTime = time.Now()
	}
}
