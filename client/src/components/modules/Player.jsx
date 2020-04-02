import React from "react";
import "./Player.css";
import "../../utilities.css";

import Card from "./Card";

export let cardHeight = 150;
export let cardWidth = 0.714859 * cardHeight;
export let offset = cardHeight / 10;

let Player = props => {
  let cardsPlayed = props.playerInfo
    ? props.playerInfo.cardsPlayed.map((c, i) => (
        <div
          style={{ zIndex: i, transform: `translateX(${offset * i}px)` }}
          className="Card-container"
        >
          <Card card={c} />
        </div>
      ))
    : null;

  return (
    <div
      className="Player-container"
      style={{
        top: props.topOffset ? props.topOffset + "%" : null,
        left: props.leftOffset ? props.leftOffset + "%" : null,
        right: props.rightOffset ? props.rightOffset + "%" : null
      }}
    >
      <div className="Player-body">
        <div
          className="Card-played-container"
          style={{
            height: `${cardHeight}px`,
            transform: `translateY(calc(-100% - 20px)) translateX(-${(cardWidth +
              offset * props.playerInfo.cardsPlayed.length) /
              2}px)`
          }}
        >
          {cardsPlayed}
        </div>
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
