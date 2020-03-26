import React from "react";

let Room = props => {
  return (
    <div className="Room-container">
      <div className="u-big-text">Test Room</div>
      <div>{props.roomid}</div>
    </div>
  );
};

export default Room;
