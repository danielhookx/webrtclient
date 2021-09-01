package signal

type Signal interface {
	Sub(target string, sdp []byte) error
	Pub() ([]byte, error)
}
