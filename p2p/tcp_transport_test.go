package p2p

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr:    ":4000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)
	go func() {
		for {
			msg := <-tr.rpcchan
			fmt.Printf("Msg: %+v", msg)
		}
	}()
	assert.Equal(t, tr.ListenAddr, ":4000")

	assert.Nil(t, tr.ListenAndAccept())

}
