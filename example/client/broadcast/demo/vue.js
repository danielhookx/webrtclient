const comp = {
    data() {
        return {
            page: 1,
            mode: null, // 'publish' | 'subscribe'
            logMsg: '',

            name: 'broadcast',
            server: '172.16.101.131:19801',
            center: '172.16.101.131:12345',
            room: 'r1',
        }
    },

    methods: {
        publish: function () {
            this.page = 2
            this.mode = 'publish'
        },

        subscribe: function () {
            this.page = 3
            this.mode = 'subscribe'
        },

        start() {
            wsStart()
            if (this.mode === 'publish') {
                navigator.mediaDevices
                    .getUserMedia({ video: true, audio: true })
                    .then((stream) => {
                        stream.getTracks().forEach((track) => pc.addTrack(track, stream))
                        document.getElementById('remoteVideos').srcObject = stream
                        pc.createOffer()
                            .then((d) => pc.setLocalDescription(d))
                            .catch((err) => (this.logMsg = err))
                    })
                    .catch((err) => (this.logMsg = err))

                axios({
                    url: 'http://' + this.center + '/publish',
                    method: 'post',
                    data: { name: this.name, room: this.room },
                }).then(function (myJson) {
                    console.log(myJson)
                })

                pc.oniceconnectionstatechange = e => this.logMsg= pc.iceConnectionState
                pc.onicecandidate = (event) => {
                    if (event.candidate === null) {
                        console.log(this.name)
                        const v = btoa(JSON.stringify(pc.localDescription))
                        document.getElementById('localSessionDescription').value = btoa(
                            JSON.stringify(pc.localDescription)
                        )
                        axios({
                            url: 'http://' + this.server + '/pub',
                            method: 'post',
                            data: { name: this.name, sdp: v },
                        }).then(function (myJson) {
                            console.log(myJson)
                        })
                    }
                }
            } else {
                axios({
                    url: 'http://' + this.center + '/subscribe',
                    method: 'post',
                    data: { name: this.name, room: this.room },
                }).then(function (myJson) {
                    console.log(myJson)
                })

                pc.ontrack = function (event) {
                    var el = document.createElement(event.track.kind)
                    el.srcObject = event.streams[0]
                    el.autoplay = true
                    el.controls = true

                    document.getElementById('remoteVideos').appendChild(el)
                }

                pc.oniceconnectionstatechange = (e) => (this.logMsg = pc.iceConnectionState)
                pc.onicecandidate = (event) => {
                    if (event.candidate === null) {
                        console.log(this.name)
                        const v = btoa(JSON.stringify(pc.localDescription))
                        document.getElementById('localSessionDescription').value = btoa(
                            JSON.stringify(pc.localDescription)
                        )
                        axios({
                            url: 'http://' + this.server + '/pub',
                            method: 'post',
                            data: { name: this.name, sdp: v },
                        }).then(function (myJson) {
                            console.log(myJson)
                        })
                    }
                }

                // pc.addTransceiver('remoteVideos')
                pc.addTransceiver('audio', {'direction': 'recvonly'})
                pc.addTransceiver('video', {'direction': 'recvonly'})
                pc.createOffer()
                    .then((d) => pc.setLocalDescription(d))
                    .catch((err) => (this.logMsg = err))

                // pc.ontrack = function (event) {
                //     var el = document.getElementById('remoteVideos')
                //     el.srcObject = event.streams[0]
                //     el.autoplay = true
                //     el.controls = true
                // }
            }
        },

        startSession() {
            let sd = document.getElementById('remoteSessionDescription').value
            if (sd === '') {
                return alert('Session Description must not be empty')
            }

            try {
                pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(sd))))
            } catch (e) {
                alert(e)
            }
        },
    },
}

const vm = Vue.createApp(comp).mount('#app')
