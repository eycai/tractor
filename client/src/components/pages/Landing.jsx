import React, { useState } from "react";
import "./Landing.css";
import "../../utilities.css";
import { post, get } from "../../api/fetch";
import { socket } from "../../client-socket";

let Landing = props => {
  let [username, setUsername] = useState("");
  let [roomcode, setRoomcode] = useState(null);
  let [inputReadOnly, setInputReadOnly] = useState(false);
  let [errorMessage, setErrorMessage] = useState(null);

  let createRoom = () => {
    post("/create_room", { name: "test name", capacity: 4 }).then(res => {
      if (res.status === 200) {
        if (res.payload) {
          props.navigate(res.payload.id);
        } else {
          console.log("ERROR, no payload in response");
        }
      } else {
        console.error(
          `unexpected status code ${res.status} with message ${res.payload}`
        );
      }
    });
  };

  let submitUsername = () => {
    post("/register", { Username: username }).then(res => {
      if (res.status === 200) {
        setInputReadOnly(true);
        post("/connect", {
          socketId: socket.id
        });
        get("/whoami").then(whoamires => {
          if (whoamires.status === 200) {
            props.setUser(whoamires.payload);
          }
        });
        setErrorMessage(null);
      } else if (res.status === 409) {
        setErrorMessage("Usename already taken! Try another.");
      } else {
        console.error(
          `unexpected status code ${res.status} with message ${res.payload}`
        );
      }
    });
  };

  let joinRoom = () => {
    if (roomcode) {
      post("/join_room", { roomId: roomcode.toUpperCase() }).then(res => {
        if (res.status === 200) {
          // if (props.updateRoom) {
          //   console.log("updating room");
          //   props.updateRoom();
          // } else {
          //   props.navigate(roomcode.toUpperCase());
          // }
          props.navigate(roomcode.toUpperCase());
          setErrorMessage(null);
        } else if (res.status === 400) {
          setErrorMessage(
            "Invalid room code. Perhaps you meant to create a new room?"
          );
        } else {
          console.error(
            `unexpected status code ${res.status} with message ${res.payload}`
          );
        }
      });
    }
  };

  let usernameInput =
    props.user && props.user.id ? (
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

  let roomCodeWidget =
    props.user && props.user.id ? (
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
          className="Landing-error u-small-text"
          style={
            errorMessage ? { visibility: "visible" } : { visibility: "hidden" }
          }
        >
          {errorMessage}
        </div>
        <div
          className="u-button Landing-start-button"
          onClick={() => {
            if (props.user && !props.user.id) {
              submitUsername();
            } else {
              joinRoom();
            }
          }}
        >
          {props.user && props.user.id ? "connect to room" : "play!"}
        </div>
        {props.user && props.user.id ? (
          <div className="u-small-text Landing-create-new" onClick={createRoom}>
            create new room...
          </div>
        ) : null}
      </div>
    </div>
  );
};

export default Landing;
