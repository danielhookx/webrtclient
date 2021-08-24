package testsignal

const defaultSrvBuf = 1000

type Channel struct {
	signal   chan []byte
}

func NewChannel() *Channel{
	ch := &Channel{
		signal: make(chan []byte, defaultSrvBuf),
	}
	return ch
}

func (ch *Channel) Push(data []byte) {
	ch.signal<- data
}

func (ch *Channel) Ready() []byte{
	return <-ch.signal
}