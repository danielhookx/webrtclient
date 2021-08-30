package broadcast

import (
	"errors"
	"github.com/pion/webrtc/v3"
	"sync"
)

var once sync.Once
var roomManager *RoomManager

type Room struct {
	localTrackChan chan *webrtc.TrackLocalStaticRTP
	subscribe chan string
	membersConn []*webrtc.PeerConnection
}

func NewRoom() *Room{
	r := &Room{
		localTrackChan: make(chan *webrtc.TrackLocalStaticRTP, 100),
		subscribe:      make(chan string, 100),
		membersConn:    make([]*webrtc.PeerConnection, 0),
	}
	return r
}

func (r *Room) GetLocalTrackChan() chan *webrtc.TrackLocalStaticRTP{
	return r.localTrackChan
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
	m sync.RWMutex
	rooms map[string]*Room
}

func GetRoomManager() *RoomManager{
	once.Do(func() {
		roomManager = &RoomManager{
			rooms: make(map[string]*Room),
		}
	})
	return roomManager
}

func (rm *RoomManager)AddRoom(room string) *Room{
	rm.m.Lock()
	defer rm.m.Unlock()
	var r *Room
	var ok bool
	if r, ok = rm.rooms[room]; !ok || r==nil {
		r = NewRoom()
		rm.rooms[room] = r
		return r
	}
	return r
}

func (rm *RoomManager)GetRoom(room string) (*Room, error){
	rm.m.RLock()
	defer rm.m.RUnlock()
	if r, ok := rm.rooms[room]; !ok || r==nil {
		return nil, errors.New("room not find")
	}else {
		return r, nil
	}
}

func (rm *RoomManager)JoinIn(room, name string) {
	rm.m.RLock()
	defer rm.m.RUnlock()
	if r, ok := rm.rooms[room]; ok && r != nil{
		r.subscribe <- name
	}
}

func (rm *RoomManager)GetSubscribe(room string) (chan string, error){
	rm.m.RLock()
	defer rm.m.RUnlock()
	if r, ok := rm.rooms[room]; ok && r != nil{
		return r.subscribe, nil
	}
	return nil, errors.New("room not find")
}