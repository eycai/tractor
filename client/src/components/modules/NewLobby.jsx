import React, { useState } from "react";
import "../../utilities.css";
import "./NewLobby.css";
import { get } from "../../api/fetch";
import { useEffect } from "react";

let NewLobby = props => {
  let [newLobbyInfo, setNewLobbyInfo] = useState({});

  return (
    <>
      <div className="NewLobby-overlay"></div>
      <div className="NewLobby-container">
        <div className="NewLobby-body">
          <div className="u-normal-text">Create New Lobby</div>
        </div>
      </div>
    </>
  );
};

export default NewLobby;
