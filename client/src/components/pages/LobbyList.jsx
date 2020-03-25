import React, { useState } from "react";
import "../../utilities.css";
import "./LobbyList.css";
import { get } from "../../api/fetch";
import { useEffect } from "react";

import NewLobby from "../modules/NewLobby";

let LobbyList = props => {
  let [lobbies, setLobbies] = useState({});
  let [newLobbyActive, setNewLobbyActive] = useState(false);

  useEffect(() => {
    get("/room_list").then(res => {
      console.log(res);
      // setLobbies()
      let lobbiesObj = {};
      for (const l of res) {
        lobbiesObj[l.id] = l;
      }
      setLobbies(lobbiesObj);
      // lobbiesUI = res.map(l => <div> TEST </div>);
      console.log(lobbiesUI);
    });
  }, []);

  const enterLobby = lobby => {};

  let lobbiesUI = Object.keys(lobbies).map(l => {
    let lobby = lobbies[l];
    return (
      <div
        className="LobbyList-entry"
        key={lobby.id}
        onClick={() => enterLobby(lobby.id)}
      >
        <div>{lobby.id}</div>
        <div>{lobby.users.length}/9</div>
      </div>
    );
  });

  return (
    <div className="LobbyList-container">
      {newLobbyActive ? <NewLobby /> : null}
      <div className="u-large-text">Lobbies</div>
      <div className="LobbyList-body">{lobbiesUI}</div>
      <div className="LobbyList-footer">
        <div
          className="u-button"
          onClick={() => {
            setNewLobbyActive(true);
          }}
        >
          create new
        </div>
      </div>
    </div>
  );
};

export default LobbyList;
