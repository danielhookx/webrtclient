/* eslint-env browser */

let pc = new RTCPeerConnection({
    iceServers: [
        {
            urls: 'stun:stun.l.google.com:19302'
        }
    ]
})
var log = msg => {
    document.getElementById('logs').innerHTML += msg + '<br>'
}

window.start = () => {
    navigator.mediaDevices.getUserMedia({ video: true, audio: true })
        .then(stream => {
            stream.getTracks().forEach(track => pc.addTrack(track, stream))
            document.getElementById('video1').srcObject = stream
            pc.createOffer().then(d => pc.setLocalDescription(d)).catch(log)
        }).catch(log)

    pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)
    pc.onicecandidate = event => {
        if (event.candidate === null) {
            console.log(vm.name);
            const v = btoa(JSON.stringify(pc.localDescription))
            document.getElementById('localSessionDescription').value = v
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
}

window.startSession = () => {
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
