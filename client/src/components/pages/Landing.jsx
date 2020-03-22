import React from "react";
import "./Landing.css";
import "../../utilities.css";

let Landing = props => {
  return (
    <div className="Landing-container">
      <div className="Landing-body">
        <div className="Landing-title">tractor.io</div>
        <div className="u-small-text">enter username below</div>
        <input spellCheck="false" className="Landing-user-input"></input>
        <div className="u-button">play!</div>
      </div>
    </div>
  );
};

export default Landing;
