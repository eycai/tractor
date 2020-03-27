import socketIOClient from "socket.io-client";
import { post } from "./api/fetch";

const endpoint = window.location.hostname + ":" + window.location.port;
export const socket = socketIOClient(endpoint);
socket.on("connect", data => {
  console.log("connected with socket");
  console.log(socket.id);
  post("/connect", {
    socketId: socket.id
  });
});
