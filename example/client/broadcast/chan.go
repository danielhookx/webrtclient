package broadcast

import (
	"errors"
	"github.com/pion/webrtc/v3"
	"log"
	"sync"
)

var once sync.Once
var roomManager *RoomManager

type Room struct {
	m              sync.RWMutex
	localTrackChan chan *webrtc.TrackLocalStaticRTP
	subSignal      chan string
	membersConn    []*webrtc.PeerConnection
	tracks         []*webrtc.TrackLocalStaticRTP
}

func NewRoom() *Room {
	r := &Room{
		localTrackChan: make(chan *webrtc.TrackLocalStaticRTP, 100),
		subSignal:      make(chan string, 100),
		membersConn:    make([]*webrtc.PeerConnection, 0),
		tracks:         make([]*webrtc.TrackLocalStaticRTP, 0),
	}
	go func() {
		for track := range r.localTrackChan {
			r.m.Lock()
			log.Printf("append track: %v", track)
			r.tracks = append(r.tracks, track)
			r.m.Unlock()
		}
	}()
	return r
}

func (r *Room) Publish(src *webrtc.TrackLocalStaticRTP) {
	r.localTrackChan <- src
}

func (r *Room) GetTracks() []*webrtc.TrackLocalStaticRTP {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.tracks
}

//maybe need close
//defer func() {
//	if cErr := peerConnection.Close(); cErr != nil {
//		fmt.Printf("cannot close peerConnection: %v\n", cErr)
//	}
//}()
func (r *Room) SaveMember(conn *webrtc.PeerConnection) {
	r.membersConn = append(r.membersConn, conn)
}

type RoomManager struct {
	m     sync.RWMutex
	rooms map[string]*Room
}

func GetRoomManager() *RoomManager {
	once.Do(func() {
		roomManager = &RoomManager{
			rooms: make(map[string]*Room),
		}
	})
	return roomManager
}

func (rm *RoomManager) AddRoom(room string) *Room {
	rm.m.Lock()
	defer rm.m.Unlock()
	var r *Room
	var ok bool
	if r, ok = rm.rooms[room]; !ok || r == nil {
		r = NewRoom()
		rm.rooms[room] = r
		return r
	}
	return r
}

func (rm *RoomManager) GetRoom(room string) (*Room, error) {
	rm.m.RLock()
	defer rm.m.RUnlock()
	if r, ok := rm.rooms[room]; !ok || r == nil {
		return nil, errors.New("room not find")
	} else {
		return r, nil
	}
}

func (rm *RoomManager) JoinIn(room, name string) {
	rm.m.RLock()
	defer rm.m.RUnlock()
	if r, ok := rm.rooms[room]; ok && r != nil {
		r.subSignal <- name
	}
}

func (rm *RoomManager) GetSubscribe(room string) (chan string, error) {
	rm.m.RLock()
	defer rm.m.RUnlock()
	if r, ok := rm.rooms[room]; ok && r != nil {
		return r.subSignal, nil
	}
	return nil, errors.New("room not find")
}
