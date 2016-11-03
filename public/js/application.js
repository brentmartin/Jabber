var socket = new WebSocket("ws://localhost:8080/socket");

socket.onmessage = function(e) {
        document.getElementById("server-message").innerHTML += e.data + "<br>";
};

function senddata() {
        var data = document.getElementById("sendtext").value;
        socket.send(data);
}
