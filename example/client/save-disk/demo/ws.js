// Create WebSocket connection.
const socket = new WebSocket('ws://'+vm.server+'/ws');

// Connection opened
socket.addEventListener('open', function (event) {
    socket.send(JSON.stringify({
        type: 1,
        data: JSON.stringify({
            name: 'demo-'+vm.name
        })
    }));
});

// Listen for messages
socket.addEventListener('message', function (event) {
    document.getElementById('remoteSessionDescription').value = event.data
});