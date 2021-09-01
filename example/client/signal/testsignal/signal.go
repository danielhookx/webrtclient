package testsignal

import (
	"errors"
	"sync"
)

type Signal struct {
	sync.RWMutex
	members map[string]*Channel
}

func NewSignal() *Signal {
	sig := &Signal{
		members: make(map[string]*Channel, 100),
	}
	return sig
}

//信令接收者注册
func (s *Signal) Reg(name string, ch *Channel) {
	s.Lock()
	s.members[name] = ch
	s.Unlock()
}

func (s *Signal) Push(name string, sdp string) error {
	s.Lock()
	defer s.Unlock()
	if ch := s.members[name]; ch == nil {
		return errors.New("member not find")
	} else {
		ch.Push([]byte(sdp))
		return nil
	}
}
