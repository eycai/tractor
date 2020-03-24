import React, { useState } from "react";
import "../../utilities.css";
import "./LobbyList.css";
import { get } from "../../api/fetch";
import { useEffect } from "react";

let LobbyList = props => {
  let [lobbies, setLobbies] = useState({});
  // let lobbiesUI = [];
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

  // useEffect(() => {
  //   console.log("use effect lobbies");
  //   console.log(lobbiesUI);
  //   console.log(lobbiesUI);
  // });

  console.log(lobbies);
  let lobbiesUI = Object.keys(lobbies).map(l => {
    let lobby = lobbies[l];
    return (
      <div className="LobbyList-entry" key={lobby.id}>
        <div>{lobby.id}</div>
        <div>{lobby.users.length}/9</div>
      </div>
    );
  });

  return (
    <div className="LobbyList-container">
      <div className="u-large-text">Lobbies</div>
      <div className="LobbyList-body">{lobbiesUI}</div>
    </div>
  );
};

export default LobbyList;
