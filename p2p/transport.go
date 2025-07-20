package p2p

type Peer interface {
	//ListenAndAccept() error
	Close() error
}

type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}

// type Closech interface {
// 	Close()
// }
