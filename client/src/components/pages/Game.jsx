import React, { useState } from "react";
import Player from "../modules/Player";
import InfoBox from "../modules/InfoBox";
import PlayerHand from "../modules/PlayerHand";
import PlayerCards from "../modules/PlayerCards";
import Chat from "../modules/Chat";

import { post } from "../../api/fetch";

import "./Game.css";
import "../../utilities.css";

//Props: user, roomid, roomInfo

// left, top, right
const playerLocations = [
  [10, 45, null],
  [50, 28, null],
  [null, 45, 10],
];

const trumpNumber = 2;

const Game = (props) => {
  const [selectedCards, setSelectedCards] = useState([]);
  const [errorText, setErrorText] = useState("");

  const playCards = () => {
    console.log(`Playing ${JSON.stringify(selectedCards)}`);
    // Disallow leads for now TODO alex make leads possible
    post("/play_cards", { cards: [selectedCards] }).then((res) => {
      console.log(res);
      if (res.status === 200) {
        setSelectedCards([]);
        setErrorText("");
      } else {
        setErrorText("cards could not be played");
      }
    });
  };

  const startDrawing = () => {
    post("/begin_drawing", {});
  };

  const declarable = () => {
    for (const c of selectedCards) {
      if (
        (c.value !== trumpNumber && c.suit !== "JOKER") ||
        c.suit !== selectedCards[0].suit
      ) {
        return false;
      }
    }
    return true;
  };

  const declareTrump = () => {
    console.log("declaring");

    if (declarable()) {
      post("/flip_cards", {
        card: selectedCards[0],
        numCards: selectedCards.length,
      }).then((res) => {
        if (res.status === 200) {
          setSelectedCards([]);
          setErrorText("");
        } else {
          setErrorText("could not flip those cards");
        }
      });
    } else {
      setErrorText("could not flip those cards");
    }
  };

  const getKitty = () => {
    post("/get_kitty");
  };

  const setKitty = () => {
    post("/set_kitty", { kitty: selectedCards }).then((res) => {
      if (res.status === 200) {
        setSelectedCards([]);
        setErrorText("");
      } else {
        setErrorText(
          "could not set the kitty.. Do you have the right number of cards?"
        );
      }
    });
  };

  let playersBefore = [];
  let playersAfter = [];
  let reached = false;

  props.roomInfo.users.forEach((u) => {
    if (u.username === props.user.username) {
      reached = true;
    } else if (reached) {
      playersAfter.push(u);
    } else {
      playersBefore.push(u);
    }
  });

  console.log(playersBefore);
  console.log(playersAfter);

  let players = playersAfter
    .concat(playersBefore)
    .reverse()
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

  let error = (
    <div className="Game-error-text-container">
      <div className="Game-error-text u-center-anchored">{errorText}</div>
    </div>
  );

  let gameButton = null;

  switch (props.roomInfo.game.gamePhase) {
    case "START":
      if (props.roomInfo.host === props.user.username)
        gameButton = (
          <div className="Game-play-button" onClick={startDrawing}>
            START DRAWING
          </div>
        );
      break;
    case "PLAYING":
      if (props.roomInfo.game.turn === props.user.username) {
        gameButton = (
          <div className="Game-play-button" onClick={playCards}>
            PLAY
          </div>
        );
      }
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
            SET KITTY {selectedCards.length}/8
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
      {error}
      <div className="Game-play-button-container">{gameButton}</div>
      {props.user.username in props.roomInfo.game.players &&
      props.roomInfo.game.players[props.user.username].cardsPlayed &&
      props.roomInfo.game.gamePhase === "PLAYING" ? (
        <div className="Game-played-cards">
          <PlayerCards
            playerInfo={props.roomInfo.game.players[props.user.username]}
            {...props}
          />
        </div>
      ) : null}
    </>
  );
};

export default Game;
