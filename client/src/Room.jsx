import React, { useEffect, useState } from "react";
import { get } from "./api/fetch";

let Room = props => {
  // RoomID
  useEffect(() => {
    get("/room_info", { RoomID: props.roomid }).then(res => {
      console.log(res);
    });
  }, []);

  return (
    <div className="Room-container">
      <div className="u-big-text">Test Room</div>
      <div>{props.roomid}</div>
    </div>
  );
};

export default Room;
