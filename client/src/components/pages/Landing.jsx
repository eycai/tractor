import React, { useState } from "react";
import "./Landing.css";
import "../../utilities.css";
import { post, get } from "../../api/fetch";

let Landing = props => {
  let [username, setUsername] = useState("");
  let [roomcode, setRoomcode] = useState(null);
  let [inputReadOnly, setInputReadOnly] = useState(false);

  let submitUsername = () => {
    post("/register", { Username: username }).then(res => {
      console.log(res);
      setInputReadOnly(true);
    });
  };

  let joinRoom = () => {
    if (roomcode) {
      post("/join_room", { RoomId: roomcode }).then(res => {
        console.log(res);
        // TODO @alex: make this not allow special chars, only letters
        props.navigate(roomcode);
      });
    }
  };

  let usernameInput = props.user ? (
    <input
      spellCheck="false"
      className="Landing-user-input"
      readOnly={true}
      value={props.user.username}
    ></input>
  ) : (
    <input
      spellCheck="false"
      className="Landing-user-input"
      readOnly={inputReadOnly}
      onChange={e => {
        setUsername(e.target.value);
      }}
    ></input>
  );

  let roomCodeWidget = props.user ? (
    <>
      <div className="u-small-text">room code</div>
      <input
        spellCheck="false"
        className="Landing-user-input"
        onChange={e => {
          setRoomcode(e.target.value);
        }}
      ></input>
    </>
  ) : null;

  return (
    <div className="Landing-container">
      <div className="Landing-body">
        <div className="Landing-title">tractor.io</div>
        <div className="u-small-text">enter username below</div>
        {usernameInput}
        {roomCodeWidget}
        <div
          className="u-button Landing-start-button"
          onClick={() => {
            if (!props.user) {
              submitUsername();
            } else {
              joinRoom();
            }
          }}
        >
          {props.user ? "connect to room" : "play!"}
        </div>
      </div>
    </div>
  );
};

export default Landing;
