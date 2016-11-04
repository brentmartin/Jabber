var socket = new ReconnectingWebSocket(location.protocol.replace("http","ws") + "//" + location.host + "/socket");

socket.onmessage = function(e) {
        document.getElementById("server-message").innerHTML += e.data + "<br>";
};

function senddata() {
        var data = document.getElementById("sendtext").value;
        socket.send(data);
}
