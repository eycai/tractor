import React from "react";

import { ReactComponent as ONECLUB } from "../../img/1CLUB.svg";
import { ReactComponent as ONESPADE } from "../../img/1SPADE.svg";
import { ReactComponent as ONEDIAMOND } from "../../img/1DIAMOND.svg";
import { ReactComponent as ONEHEART } from "../../img/1HEART.svg";
import { ReactComponent as TWOCLUB } from "../../img/2CLUB.svg";
import { ReactComponent as TWOSPADE } from "../../img/2SPADE.svg";
import { ReactComponent as TWODIAMOND } from "../../img/2DIAMOND.svg";
import { ReactComponent as TWOHEART } from "../../img/2HEART.svg";
import "./Player.css";

let convertCard = card => {
  return `${card.value.toString()}${card.suit}`;
};

let Card = props => {
  let cardSVG = null;
  switch (convertCard(props.card)) {
    case "1CLUB":
      cardSVG = <ONECLUB className="Card-body" />;
      break;
    case "1SPADE":
      cardSVG = <ONESPADE className="Card-body" />;
      break;
    case "1HEART":
      cardSVG = <ONEHEART className="Card-body" />;
      break;
    case "1DIAMOND":
      cardSVG = <ONEDIAMOND className="Card-body" />;
      break;
    case "2CLUB":
      cardSVG = <TWOCLUB className="Card-body" />;
      break;
    case "2SPADE":
      cardSVG = <TWOSPADE className="Card-body" />;
      break;
    case "2HEART":
      cardSVG = <TWOHEART className="Card-body" />;
      break;
    case "2DIAMOND":
      cardSVG = <TWODIAMOND className="Card-body" />;
      break;
  }

  return <>{cardSVG}</>;
};

export default Card;
