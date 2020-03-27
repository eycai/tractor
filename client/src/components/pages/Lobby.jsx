import React from "react";
import "./Lobby.css";
import Header from "../modules/Header";

//Props: user, roomid, roomInfo

let Lobby = props => {
  let usersList = props.roomInfo
    ? props.roomInfo.users.map(u => (
        <div className="u-medium-text" key={u}>
          {u}
        </div>
      ))
    : null;

  return (
    <>
      <Header />
      <div className="Lobby-container">
        <div className="Lobby-body">
          <div className="u-title-text">
            {props.user ? `${props.user.username}'s Lobby` : "Loading..."}
          </div>
          <div className="u-large-text Lobby-room-code">
            Room Code: {props.roomid}
          </div>
          <div className="u-large-text">Connected Users:</div>
          <div className="Lobby-players-container"> {usersList} </div>
        </div>
      </div>
    </>
  );
};

export default Lobby;
