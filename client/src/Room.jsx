import React, { useEffect, useState } from "react";
import Lobby from "./components/pages/Lobby";
import { socket } from "./client-socket";
import { get } from "./api/fetch";
import "./utilities.css";

// Props: user, roomid, socket

let Room = props => {
  // RoomID
  let [roomInfo, setRoomInfo] = useState(null);

  useEffect(() => {
    get("/room_info", { roomId: props.roomid }).then(res => {
      console.log(res);
      setRoomInfo(res);
    });
    socket.on("update", data => {
      console.log("got an update on this room.");
      console.log(data);
      setRoomInfo(data.room);
      props.setUser(data.user);
    });
  }, []);

  return (
    <div>
      <Lobby roomInfo={roomInfo} {...props} />
    </div>
  );
};

export default Room;
