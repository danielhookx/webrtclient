package main

import (
	"flag"
	"fmt"
	rtp_to_webrtc "github.com/oofpgDLD/webrtclient/example/client/rtp-to-webrtc"
	"github.com/oofpgDLD/webrtclient/example/client/signal/testsignal/client"
	"github.com/oofpgDLD/webrtclient/internal/signal"
	"github.com/pion/webrtc/v3"
	"net/url"
	"os"
	"strconv"
)

var addr = flag.String("addr", os.Getenv("SIGNALADDR"), "http service address")
var inputAddr = flag.String("input", "127.0.0.1:4002", "http service address")
var goName = "rtp-to-webrtc"
var jsName = "demo-" + goName

func main() {
	flag.Parse()
	var u = url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var inputUrl = url.URL{Host: *inputAddr}
	engine := &webrtc.MediaEngine{}
	if err := engine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeVP8,
			ClockRate:    90000,
			Channels:     0,
			SDPFmtpLine:  "",
			RTCPFeedback: nil,
		},
		PayloadType: 0,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}

	if err := engine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeOpus,
			ClockRate:    48000,
			Channels:     0,
			SDPFmtpLine:  "",
			RTCPFeedback: nil,
		},
		PayloadType: 0,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}

	//初始化 信令服务
	c, err := client.NewClient(goName, u)
	if err != nil {
		panic(err)
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(engine))

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	//
	port, err := strconv.Atoi(inputUrl.Port())
	if err != nil {
		panic(err)
	}
	rtp_to_webrtc.RTPToWebInit(peerConnection, inputUrl.Host, port)

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
	err = client.Put("http://"+*addr+"/pub", jsName, signal.Encode(*peerConnection.LocalDescription()))
	if err != nil {
		fmt.Println("put sdp err:", err.Error())
	}

	go rtp_to_webrtc.ServeVideoTrack()
	defer rtp_to_webrtc.Close()
	// Block forever
	select {}
}
