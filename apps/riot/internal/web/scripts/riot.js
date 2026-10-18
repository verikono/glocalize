// riot browser client: opens a websocket to "/ws" and reloads the page when
// the server dispatches a refresh message.
(function () {
	"use strict";

	var RECONNECT_MS = 1000;

	function connect() {
		var scheme = window.location.protocol === "https:" ? "wss:" : "ws:";
		var url = scheme + "//" + window.location.host + "/ws";
		var socket = new WebSocket(url);

		socket.addEventListener("message", function (event) {
			var message = String(event.data || "");
			if (message.indexOf("riot::refresh") === 0) {
				window.location.reload();
			}
		});

		socket.addEventListener("close", function () {
			// The server restarts on every hot rebuild, so reconnect and wait
			// for the next refresh message.
			window.setTimeout(connect, RECONNECT_MS);
		});

		socket.addEventListener("error", function () {
			socket.close();
		});
	}

	connect();
})();
