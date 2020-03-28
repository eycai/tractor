import React from "react";
import "./Lobby.css";
import Header from "../modules/Header";
import Exit from "../modules/Exit";

//Props: user, roomid, roomInfo

const Lobby = props => {
  const usersList = props.roomInfo
    ? props.roomInfo.users.map(u => (
        <div className="u-medium-text" key={u}>
          {u}
        </div>
      ))
    : null;

  const startGame = () => {
    console.log("starting game");
  };

  const changeSettings = () => {
    console.log("changing settings");
  };

  const hostControls =
    props.roomInfo && props.roomInfo.host == props.user.username ? (
      <div className="Lobby-controls-container">
        <div onClick={startGame} className="u-button Lobby-button">
          start
        </div>
        <div onClick={changeSettings} className="u-button Lobby-button">
          settings
        </div>
      </div>
    ) : null;

  return (
    <>
      <Header />
      <div className="Lobby-container">
        <div className="Lobby-exit"></div>
        <div className="Lobby-body">
          <div className="u-title-text">
            {props.roomInfo ? `${props.roomInfo.host}'s Lobby` : "Loading..."}
          </div>
          <div className="u-large-text Lobby-room-code">
            Room Code: {props.roomid}
          </div>
          <div className="u-large-text">Connected Users:</div>
          <div className="Lobby-players-container"> {usersList} </div>
          {hostControls}
          {/* TODO @alex make this pretty */}
          <div
            onClick={() => props.navigate("/")}
            className="Lobby-exit-container"
          >
            <Exit />
          </div>
        </div>
      </div>
    </>
  );
};

export default Lobby;
