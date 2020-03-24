import React, { useState } from "react";
import "./Landing.css";
import "../../utilities.css";
import { post, get } from "../../api/fetch";
import { pages } from "../../utilities.js";

let Landing = props => {
  let [username, setUsername] = useState("");
  let [inputReadOnly, setInputReadOnly] = useState(false);

  let submitUsername = () => {
    post("/register", { Username: username }).then(res => {
      console.log(res);
      setInputReadOnly(true);
    });
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

  return (
    <div className="Landing-container">
      <div className="Landing-body">
        <div className="Landing-title">tractor.io</div>
        <div className="u-small-text">enter username below</div>
        {usernameInput}
        <div
          className="u-button"
          onClick={() => {
            if (!props.user) {
              submitUsername();
            }
            props.setPage(pages.LOBBYLIST);
          }}
        >
          play!
        </div>
      </div>
    </div>
  );
};

export default Landing;
