import React from "react";
import Player from "../modules/Player";
import InfoBox from "../modules/InfoBox";
import PlayerHand from "../modules/PlayerHand";

import "./Game.css";
import "../../utilities.css";

//Props: user, roomid, roomInfo

const playerLocations = [
  [5, 45, null],
  [50, 5, null],
  [null, 45, 5]
];

const Game = props => {
  console.log(playerLocations);
  let game = props.roomInfo.game;

  let players = props.roomInfo.users
    .filter(p => p.username !== props.user.username)
    .map((p, i) => (
      <Player
        key={p.username}
        playerInfo={props.roomInfo.game.players[p.username]}
        leftOffset={playerLocations[i][0]}
        topOffset={playerLocations[i][1]}
        rightOffset={playerLocations[i][2]}
        centered={
          Object.keys(game.players).length % 2 === 0 &&
          i === Object.keys(game.players).length / 2 - 1
        }
      />
    ));

  let centerText = (
    <div className="Game-center-text-container">
      <div className="Game-center-text u-y-centered">
        <span className="Game-current-turn">{`${props.roomInfo.game.turn}'s`}</span>
        <span> Turn</span>
      </div>
    </div>
  );
  return (
    <>
      {players}
      <InfoBox {...props} />
      <PlayerHand {...props} />
      {centerText}
      {true ? (
        <div className="Game-play-button-container">
          <div className="Game-play-button">PLAY</div>
        </div>
      ) : null}
    </>
  );
};

export default Game;
