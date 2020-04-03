import React, { useEffect, useState } from "react";
import Lobby from "./components/pages/Lobby";
import Game from "./components/pages/Game";
import { socket } from "./client-socket";
import { testData } from "./utilities";
import { post, get } from "./api/fetch";
import "./utilities.css";

// Props: user, roomid, socket

let Room = props => {
  // RoomID
  let [roomInfo, setRoomInfo] = useState(null);

  useEffect(() => {
    get("/room_info", { roomId: props.roomid }).then(res => {
      if (res.status === 200) {
        setRoomInfo(res.payload);
      } else {
        console.error(
          `unexpected status code ${res.status} with message ${res.payload}`
        );
      }
    });
    // post("/test_set_room", testData);
    socket.on("update", data => {
      console.log("got an update on this room.");
      console.log(data);
      setRoomInfo(data.room);
      props.setUser(data.user);
    });
  }, []);

  return (
    <div
      style={{
        backgroundColor: "#F9F9F9",
        width: "100vw",
        height: "100vh",
        zIndex: "-1000"
      }}
    >
      {roomInfo && roomInfo.game ? (
        <Game roomInfo={roomInfo} {...props} />
      ) : (
        <Lobby roomInfo={roomInfo} {...props} />
      )}
    </div>
  );
};

export default Room;
