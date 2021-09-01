package rtp_to_webrtc

import (
	"errors"
	"fmt"
	"github.com/pion/webrtc/v3"
	"io"
	"net"
)

var listener *net.UDPConn
var videoTrack *webrtc.TrackLocalStaticRTP

func RTPToWebInit(peerConnection *webrtc.PeerConnection, ip string, port int) {
	// Open a UDP Listener for RTP Packets on port 5004
	var err error
	listener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		panic(err)
	}

	// Create a video track
	videoTrack, err = webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	if err != nil {
		panic(err)
	}
	rtpSender, err := peerConnection.AddTrack(videoTrack)
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

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateFailed {
			if closeErr := peerConnection.Close(); closeErr != nil {
				panic(closeErr)
			}
		}
	})
}

func ServeVideoTrack() {
	// Read RTP packets forever and send them to the WebRTC Client
	inboundRTPPacket := make([]byte, 1600) // UDP MTU
	for {
		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			panic(fmt.Sprintf("error during read: %s", err))
		}

		if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				// The peerConnection has been closed.
				return
			}

			panic(err)
		}
	}
}

func Close() {
	if err := listener.Close(); err != nil {
		panic(err)
	}
}
