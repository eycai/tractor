import React from "react";
import "./Player.css";
import "../../utilities.css";

let Player = props => {
  return (
    <div
      className="Player-container"
      style={{
        top: props.topOffset ? props.topOffset + "%" : null,
        left: props.leftOffset ? props.leftOffset + "%" : null,
        right: props.rightOffset ? props.rightOffset + "%" : null
      }}
    >
      <div
        className={`u-y-centered ${props.centered ? "u-center-anchored" : ""}`}
      >
        <div
          className={`Player-username ${
            props.playerInfo.team === "BOSSES"
              ? "Player-team-a"
              : "Player-team-b"
          }`}
        >
          {props.playerInfo.username}
        </div>
      </div>
    </div>
  );
};

export default Player;
