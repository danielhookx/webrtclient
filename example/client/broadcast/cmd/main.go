package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/oofpgDLD/webrtclient/example/client/broadcast"
	"github.com/oofpgDLD/webrtclient/example/client/signal/testsignal/client"
	"github.com/oofpgDLD/webrtclient/internal/signal"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"io"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	rtcpPLIInterval = time.Second * 3
)

var addr = flag.String("addr", os.Getenv("SIGNALADDR"), "http service address")
var port = flag.Int("port", 12345, "http service address")

var peerConnectionConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

//map[room]trackChan
var roomManager *broadcast.RoomManager

func main() {
	flag.Parse()
	var u = url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	pubChan, subChan := signal.HTTPPubSubServer(*port)
	log.Printf("serve http hub %v", *port)
	roomManager = broadcast.GetRoomManager()
	go func() {
		for {
			hub := <-pubChan
			log.Printf("publish name=%v,room=%v", hub.Name, hub.Room)
			go func() {
				Publisher(hub.Name, hub.Room, u)
				log.Printf("ServeSubscribe room=%v", hub.Room)
				err := ServeSubscribe(hub.Room, u)
				if err != nil {
					log.Printf("ServeSubscribe failed:%v", err)
				}
			}()
		}
	}()

	go func() {
		for {
			hub := <-subChan
			log.Printf("subscribe name=%v,room=%v", hub.Name, hub.Room)
			roomManager.JoinIn(hub.Room, hub.Name)
		}
	}()
	// Block forever
	select {}
}


func Publisher(goName, room string, u url.URL) {
	jsName := "demo-"+goName
	//初始化 信令服务
	c, err := client.NewClient(goName, u)
	if err != nil {
		log.Printf("Publisher failed:%v", err.Error())
		return
	}
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		panic(err)
	}

	// Allow us to receive 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}

	roomEng := roomManager.AddRoom(room)
	localTrackChan := roomEng.GetLocalTrackChan()
	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(remoteTrack.SSRC())}}); rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		// Create a local track, all our SFU clients will be fed via this track
		localTrack, newTrackErr := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, "video", "pion")
		if newTrackErr != nil {
			panic(newTrackErr)
		}
		localTrackChan <- localTrack
		log.Print("add local track")
		rtpBuf := make([]byte, 1400)
		for {
			i, _, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				panic(readErr)
			}

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err = localTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				panic(err)
			}
		}
	})

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	//signal.Decode(signal.MustReadStdin(), &offer)
	signal.Decode(string(c.Read()), &offer)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	fmt.Println("ICE Gathering start")
	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete
	fmt.Println("ICE Gathering complete")

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))
	err = client.Put("http://"+ *addr + "/pub", jsName, signal.Encode(*peerConnection.LocalDescription()))
	if err != nil {
		fmt.Println("put sdp err:", err.Error())
	}
	roomEng.SaveMember(peerConnection)
}

func ServeSubscribe(room string, u url.URL) error{
	roomEng, err := roomManager.GetRoom(room)
	if err != nil {
		return err
	}
	localTrackChan := roomEng.GetLocalTrackChan()
	localTrack := <-localTrackChan
	log.Printf("subscribe start listen:room=%v", room)
	for {
		subChan, err := roomManager.GetSubscribe(room)
		if err != nil {
			return err
		}
		subName := <-subChan
		log.Printf("subscribe start name=%v,room=%v", subName, room)
		go Subscribe(localTrack, subName, u)
	}
}

func Subscribe(localTrack *webrtc.TrackLocalStaticRTP, goName string, u url.URL) {
	fmt.Println("")
	fmt.Println("Curl an base64 SDP to start sendonly peer connection")

	jsName := "demo-"+goName
	//初始化 信令服务
	c, err := client.NewClient(goName, u)
	if err != nil {
		log.Printf("Subscribe failed: %v\n", err.Error())
	}
	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	//signal.Decode(signal.MustReadStdin(), &offer)
	signal.Decode(string(c.Read()), &offer)

	// Create a new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		panic(err)
	}

	rtpSender, err := peerConnection.AddTrack(localTrack)
	if err != nil {
		panic(err)
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete
	fmt.Println("ICE Gathering complete")

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))
	err = client.Put("http://"+ *addr + "/pub", jsName, signal.Encode(*peerConnection.LocalDescription()))
	if err != nil {
		fmt.Println("put sdp err:", err.Error())
	}
}