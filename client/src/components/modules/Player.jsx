import React from "react";
import "./Player.css";
import "../../utilities.css";

import Card from "./Card";
import PlayerCards from "./PlayerCards";

export let cardHeight = 150;
export let cardWidth = 0.714859 * cardHeight;
export let offset = cardHeight / 10;

let Player = (props) => {
  return (
    <div
      className="Player-container"
      style={{
        top: props.topOffset ? props.topOffset + "%" : null,
        left: props.leftOffset ? props.leftOffset + "%" : null,
        right: props.rightOffset ? props.rightOffset + "%" : null,
      }}
    >
      <div className="Player-body">
        {props.playerInfo.cardsPlayed ? <PlayerCards {...props} /> : null}
        <div className="Player-username-container u-center-anchored">
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
    </div>
  );
};

export default Player;
