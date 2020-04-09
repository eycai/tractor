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

const trumpNumber = 2;

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

  const declarable = () => {
    const cards = selectedCards.map(i => props.user.hand[i]);
    for (const c of cards) {
      if (c.value !== trumpNumber && c.suit === cards[0].value) {
        return false;
      }
    }
    return true;
  };

  const declareTrump = () => {
    console.log("declaring");
    const cards = selectedCards.map(i => props.user.hand[i]);
    if (declarable()) {
      post("/flip_cards", { card: cards[0], numCards: cards.length });
    } else {
      console.log("invalid");
    }
  };

  const getKitty = () => {
    console.log("getting the kitty");
    post("/get_kitty");
  };

  const setKitty = () => {
    post("/set_kitty");
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
    case "DRAWING":
      gameButton = (
        <div className="Game-play-button" onClick={declareTrump}>
          DECLARE
        </div>
      );
      break;
    case "DRAWING_COMPLETE":
      if (props.user.username === props.roomInfo.game.trumpFlipUser) {
        gameButton = (
          <div className="Game-play-button" onClick={getKitty}>
            GET KITTY
          </div>
        );
      } else {
        gameButton = (
          <div className="Game-play-button" onClick={declareTrump}>
            DECLARE
          </div>
        );
      }
      break;
    case "SET_KITTY":
      if (props.user.username === props.roomInfo.game.trumpFlipUser) {
        gameButton = (
          <div className="Game-play-button" onClick={setKitty}>
            SET KITTY
          </div>
        );
      }
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
