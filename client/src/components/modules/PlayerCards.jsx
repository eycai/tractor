import React from "react";

import "./Player.css";
import Card from "./Card";

export let cardHeight = 150;
export let cardWidth = 0.714859 * cardHeight;
export let offset = cardHeight / 10;

let PlayerCards = (props) => {
  let cardsPlayed =
    props.playerInfo && props.playerInfo.cardsPlayed
      ? props.playerInfo.cardsPlayed
          .reduce((a, e) => a.concat(e))
          .map((c, i) => (
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
      className="Card-played-container"
      style={{
        height: `${cardHeight}px`,
        transform: `translateY(calc(-100% - 20px)) translateX(-${(cardWidth +
          offset * props.playerInfo.cardsPlayed.length) /
          2}px)`,
      }}
    >
      {cardsPlayed}
    </div>
  );
};

export default PlayerCards;
