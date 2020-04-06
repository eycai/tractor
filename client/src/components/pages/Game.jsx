import React, { useState } from "react";
import Player from "../modules/Player";
import InfoBox from "../modules/InfoBox";
import PlayerHand from "../modules/PlayerHand";
import Chat from "../modules/Chat";

import { post } from "../../api/fetch";

import "./Game.css";
import "../../utilities.css";

//Props: user, roomid, roomInfo

const playerLocations = [
  [10, 45, null],
  [50, 28, null],
  [null, 45, 10]
];

const Game = props => {
  console.log(playerLocations);
  let game = props.roomInfo.game;
  const [selectedCards, setSelectedCards] = useState([]);

  const playCards = () => {
    const cards = selectedCards.map(i => props.user.hand[i]);
    console.log(`Playing ${JSON.stringify(cards)}`);
  };

  const startDrawing = () => {
    post("/begin_drawing", {});
  };

  let players = props.roomInfo.users
    .filter(p => p.username !== props.user.username)
    .map((p, i) => (
      <Player
        key={p.username}
        playerInfo={props.roomInfo.game.players[p.username]}
        leftOffset={playerLocations[i][0]}
        topOffset={playerLocations[i][1]}
        rightOffset={playerLocations[i][2]}
      />
    ));

  let centerText = (
    <div className="Game-center-text-container">
      <div className="Game-center-text u-center-anchored">
        <span className="Game-current-turn">{`${props.roomInfo.game.turn}'s`}</span>
        <span>Turn</span>
      </div>
    </div>
  );

  let gameButton = null;

  switch (props.roomInfo.game.gamePhase) {
    case "START":
      gameButton = (
        <div className="Game-play-button" onClick={startDrawing}>
          START DRAWING
        </div>
      );
      break;
    case "PLAYING":
      gameButton = (
        <div className="Game-play-button" onClick={playCards}>
          PLAY
        </div>
      );
      break;
  }
  return (
    <>
      {players}
      <div className="Game-player-dash">
        <InfoBox {...props} />
        {props.user.hand ? (
          <PlayerHand
            selectedCards={selectedCards}
            setSelectedCards={setSelectedCards}
            {...props}
          />
        ) : null}
        <Chat />
      </div>
      {centerText}
      {true ? (
        <div className="Game-play-button-container">{gameButton}</div>
      ) : null}
    </>
  );
};

export default Game;
