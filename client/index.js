let ws;

const connectBtn = document.getElementById("connect-btn");
const disconnectBtn = document.getElementById("disconnect-btn");
const pingBtn = document.getElementById("ping-btn");

const pong = document.getElementById("pong-msg");

connectBtn.onclick = function () {
  if (!window["WebSocket"]) {
    alert("Not supporting websockets");
    return;
  }

  ws = new WebSocket(
    "ws://" + document.location.host + "/notification-service/ws"
  );
  connectBtn.textContent = "Connected";
  connectBtn.disabled = true;
  disconnectBtn.disabled = false;
  pingBtn.disabled = false;

  ws.onmessage = function (ev) {
    pong.textContent = ev.data;
  };

  ws.onclose = function (ev) {
    connectBtn.textContent = "Connect";
    connectBtn.disabled = false;
    disconnectBtn.disabled = true;
    pingBtn.disabled = true;
  };
};

disconnectBtn.onclick = function () {
  ws?.close();
};

class Event {
  constructor(type, payload) {
    this.type = type;
    this.payload = payload;
  }

  toJson() {
    return JSON.stringify(this);
  }
}

pingBtn.onclick = function () {
  const event = new Event("health", { msg: "ping", time: new Date() });
  console.log(event);
  ws.send(event.toJson());
};
