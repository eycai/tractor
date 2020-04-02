import React, { useState } from "react";
import Player from "../modules/Player";
import InfoBox from "../modules/InfoBox";
import PlayerHand from "../modules/PlayerHand";
import Chat from "../modules/Chat";

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
  return (
    <>
      {players}
      <div className="Game-player-dash">
        <InfoBox {...props} />
        <PlayerHand
          selectedCards={selectedCards}
          setSelectedCards={setSelectedCards}
          {...props}
        />
        <Chat />
      </div>
      {centerText}
      {true ? (
        <div className="Game-play-button-container">
          <div className="Game-play-button" onClick={playCards}>
            PLAY
          </div>
        </div>
      ) : null}
    </>
  );
};

export default Game;
