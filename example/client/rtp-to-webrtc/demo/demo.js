/* eslint-env browser */

let pc = new RTCPeerConnection({
    iceServers: [
        {
            urls: 'stun:stun.l.google.com:19302'
        }
    ]
})
let log = msg => {
    document.getElementById('logs').innerHTML += msg + '<br>'
}

window.start = () => {

    pc.ontrack = function (event) {
        var el = document.createElement(event.track.kind)
        el.srcObject = event.streams[0]
        el.autoplay = true
        el.controls = true

        document.getElementById('remoteVideos').appendChild(el)
    }

    pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)
    pc.onicecandidate = event => {
        if (event.candidate === null) {
            console.log(vm.name);
            const v = btoa(JSON.stringify(pc.localDescription))
            document.getElementById('localSessionDescription').value = btoa(JSON.stringify(pc.localDescription))
            axios({
                url: 'http://'+vm.server+'/pub',
                method: 'post',
                data: { name: vm.name, sdp: v }
            })
                .then(function (myJson) {
                    console.log(myJson);
                });
        }
    }

// Offer to receive 1 audio, and 2 video tracks
    pc.addTransceiver('audio', {'direction': 'recvonly'})
    pc.addTransceiver('video', {'direction': 'recvonly'})
    pc.addTransceiver('video', {'direction': 'recvonly'})
    pc.createOffer().then(d => pc.setLocalDescription(d)).catch(log)
}

window.startSession = () => {
    console.log("click start");
    let sd = document.getElementById('remoteSessionDescription').value
    if (sd === '') {
        return alert('Session Description must not be empty')
    }

    try {
        pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(sd))))
    } catch (e) {
        alert(e)
    }
}