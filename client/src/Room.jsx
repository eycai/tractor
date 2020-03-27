import React, { useEffect, useState } from "react";
import { get } from "./api/fetch";
import "./utilities.css";

let Room = props => {
  // RoomID
  let [roomInfo, setRoomInfo] = useState(null);

  useEffect(() => {
    get("/room_info", { RoomID: props.roomid }).then(res => {
      console.log(res);
      setRoomInfo(res);
    });
  }, []);

  let usersList = roomInfo
    ? roomInfo.users.map(u => <div>{u.username}</div>)
    : null;

  return (
    <div className="Room-container">
      <div className="u-large-text">
        {props.user ? `${props.user.username}'s Room` : "Loading..."}
      </div>
      <div>{props.roomid}</div>
      <div className="u-normal-text">Connected Users:</div>
      <div> {usersList} </div>
    </div>
  );
};

export default Room;
