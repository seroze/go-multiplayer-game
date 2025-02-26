const canvas = document.getElementById("gameCanvas");
const ctx = canvas.getContext("2d");

const socket = new WebSocket("ws://localhost:8080/ws");

let players = {};

// Handle incoming messages from server
socket.onmessage = function (event) {
  console.log("Received from server:", event.data); // Debugging
  players = JSON.parse(event.data);
  drawPlayers();
};

// Handle key presses for movement
document.addEventListener("keydown", function (event) {
  let key = "";
  if (event.key === "ArrowUp") key = "up";
  if (event.key === "ArrowDown") key = "down";
  if (event.key === "ArrowLeft") key = "left";
  if (event.key === "ArrowRight") key = "right";

  if (key) {
    socket.send(key);
  }
});

// Draw all players on the canvas
function drawPlayers() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  for (let id in players) {
    let player = players[id];
    ctx.fillStyle = player.color;
    ctx.fillRect(player.X, player.Y, 20, 20); // Draw player as square
  }
}
