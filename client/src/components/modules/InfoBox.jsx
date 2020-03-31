import React from "react";
import "./InfoBox.css";

let InfoBox = props => {
  let points = Object.values(props.roomInfo.game.players).reduce(
    (t, a) => t + a.points,
    0
  );

  return (
    <div className="InfoBox-container">
      <div className="InfoBox-points-body">
        <span className="InfoBox-label">Points: </span>
        <span className="InfoBox-points">{points}</span>
        <span className="InfoBox-total-points">/80</span>
      </div>
      <div>
        <span className="InfoBox-label">Declared: </span>{" "}
        {props.roomInfo.game.trumpSuit} {props.roomInfo.game.trumpNumber}
      </div>
    </div>
  );
};

export default InfoBox;
