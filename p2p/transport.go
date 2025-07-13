package p2p

type Peer interface {
	//ListenAndAccept() error
	Close() error
}

type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}

// type Closech interface {
// 	Close()
// }
