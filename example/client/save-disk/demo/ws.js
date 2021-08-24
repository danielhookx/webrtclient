// Create WebSocket connection.
const socket = new WebSocket('ws://172.16.101.131:19801/ws');

// Connection opened
socket.addEventListener('open', function (event) {
    socket.send(JSON.stringify({
        type: 1,
        data: JSON.stringify({
            name: 'demo'
        })
    }));
});

// Listen for messages
socket.addEventListener('message', function (event) {
    document.getElementById('remoteSessionDescription').value = event.data
});