import React from "react";

import { ReactComponent as ONECLUB } from "../../img/1CLUB.svg";
import { ReactComponent as ONESPADE } from "../../img/1SPADE.svg";
import { ReactComponent as ONEDIAMOND } from "../../img/1DIAMOND.svg";
import { ReactComponent as ONEHEART } from "../../img/1HEART.svg";
import { ReactComponent as TWOCLUB } from "../../img/2CLUB.svg";
import { ReactComponent as TWOSPADE } from "../../img/2SPADE.svg";
import { ReactComponent as TWODIAMOND } from "../../img/2DIAMOND.svg";
import { ReactComponent as TWOHEART } from "../../img/2HEART.svg";
import { ReactComponent as THREECLUB } from "../../img/3CLUB.svg";
import { ReactComponent as THREESPADE } from "../../img/3SPADE.svg";
import { ReactComponent as THREEDIAMOND } from "../../img/3DIAMOND.svg";
import { ReactComponent as THREEHEART } from "../../img/3HEART.svg";
import { ReactComponent as FOURCLUB } from "../../img/4CLUB.svg";
import { ReactComponent as FOURSPADE } from "../../img/4SPADE.svg";
import { ReactComponent as FOURDIAMOND } from "../../img/4DIAMOND.svg";
import { ReactComponent as FOURHEART } from "../../img/4HEART.svg";
import { ReactComponent as FIVECLUB } from "../../img/5CLUB.svg";
import { ReactComponent as FIVESPADE } from "../../img/5SPADE.svg";
import { ReactComponent as FIVEDIAMOND } from "../../img/5DIAMOND.svg";
import { ReactComponent as FIVEHEART } from "../../img/5HEART.svg";
import { ReactComponent as SIXCLUB } from "../../img/6CLUB.svg";
import { ReactComponent as SIXSPADE } from "../../img/6SPADE.svg";
import { ReactComponent as SIXDIAMOND } from "../../img/6DIAMOND.svg";
import { ReactComponent as SIXHEART } from "../../img/6HEART.svg";
import { ReactComponent as SEVENCLUB } from "../../img/7CLUB.svg";
import { ReactComponent as SEVENSPADE } from "../../img/7SPADE.svg";
import { ReactComponent as SEVENDIAMOND } from "../../img/7DIAMOND.svg";
import { ReactComponent as SEVENHEART } from "../../img/7HEART.svg";
import { ReactComponent as EIGHTCLUB } from "../../img/8CLUB.svg";
import { ReactComponent as EIGHTSPADE } from "../../img/8SPADE.svg";
import { ReactComponent as EIGHTDIAMOND } from "../../img/8DIAMOND.svg";
import { ReactComponent as EIGHTHEART } from "../../img/8HEART.svg";
import { ReactComponent as NINECLUB } from "../../img/9CLUB.svg";
import { ReactComponent as NINESPADE } from "../../img/9SPADE.svg";
import { ReactComponent as NINEDIAMOND } from "../../img/9DIAMOND.svg";
import { ReactComponent as NINEHEART } from "../../img/9HEART.svg";
import { ReactComponent as TENCLUB } from "../../img/10CLUB.svg";
import { ReactComponent as TENSPADE } from "../../img/10SPADE.svg";
import { ReactComponent as TENDIAMOND } from "../../img/10DIAMOND.svg";
import { ReactComponent as TENHEART } from "../../img/10HEART.svg";

import JACKCLUB from "../../img/11CLUB.png";
import JACKSPADE from "../../img/11SPADE.png";
import JACKDIAMOND from "../../img/11DIAMOND.png";
import JACKHEART from "../../img/11HEART.png";
import QUEENCLUB from "../../img/12CLUB.png";
import QUEENSPADE from "../../img/12SPADE.png";
import QUEENDIAMOND from "../../img/12DIAMOND.png";
import QUEENHEART from "../../img/12HEART.png";
import KINGCLUB from "../../img/13CLUB.png";
import KINGSPADE from "../../img/13SPADE.png";
import KINGDIAMOND from "../../img/13DIAMOND.png";
import KINGHEART from "../../img/13HEART.png";

import BIGJOKER from "../../img/1JOKER.png";
import LITTLEJOKER from "../../img/2JOKER.png";

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
    case "3CLUB":
      cardSVG = <THREECLUB className="Card-body" />;
      break;
    case "3SPADE":
      cardSVG = <THREESPADE className="Card-body" />;
      break;
    case "3HEART":
      cardSVG = <THREEHEART className="Card-body" />;
      break;
    case "3DIAMOND":
      cardSVG = <THREEDIAMOND className="Card-body" />;
      break;
    case "4CLUB":
      cardSVG = <FOURCLUB className="Card-body" />;
      break;
    case "4SPADE":
      cardSVG = <FOURSPADE className="Card-body" />;
      break;
    case "4HEART":
      cardSVG = <FOURHEART className="Card-body" />;
      break;
    case "4DIAMOND":
      cardSVG = <FOURDIAMOND className="Card-body" />;
      break;
    case "5CLUB":
      cardSVG = <FIVECLUB className="Card-body" />;
      break;
    case "5SPADE":
      cardSVG = <FIVESPADE className="Card-body" />;
      break;
    case "5HEART":
      cardSVG = <FIVEHEART className="Card-body" />;
      break;
    case "5DIAMOND":
      cardSVG = <FIVEDIAMOND className="Card-body" />;
      break;
    case "6CLUB":
      cardSVG = <SIXCLUB className="Card-body" />;
      break;
    case "6SPADE":
      cardSVG = <SIXSPADE className="Card-body" />;
      break;
    case "6HEART":
      cardSVG = <SIXHEART className="Card-body" />;
      break;
    case "6DIAMOND":
      cardSVG = <SIXDIAMOND className="Card-body" />;
      break;
    case "7CLUB":
      cardSVG = <SEVENCLUB className="Card-body" />;
      break;
    case "7SPADE":
      cardSVG = <SEVENSPADE className="Card-body" />;
      break;
    case "7HEART":
      cardSVG = <SEVENHEART className="Card-body" />;
      break;
    case "7DIAMOND":
      cardSVG = <SEVENDIAMOND className="Card-body" />;
      break;
    case "8CLUB":
      cardSVG = <EIGHTCLUB className="Card-body" />;
      break;
    case "8SPADE":
      cardSVG = <EIGHTSPADE className="Card-body" />;
      break;
    case "8HEART":
      cardSVG = <EIGHTHEART className="Card-body" />;
      break;
    case "8DIAMOND":
      cardSVG = <EIGHTDIAMOND className="Card-body" />;
      break;
    case "9CLUB":
      cardSVG = <NINECLUB className="Card-body" />;
      break;
    case "9SPADE":
      cardSVG = <NINESPADE className="Card-body" />;
      break;
    case "9HEART":
      cardSVG = <NINEHEART className="Card-body" />;
      break;
    case "9DIAMOND":
      cardSVG = <NINEDIAMOND className="Card-body" />;
      break;
    case "10CLUB":
      cardSVG = <TENCLUB className="Card-body" />;
      break;
    case "10SPADE":
      cardSVG = <TENSPADE className="Card-body" />;
      break;
    case "10HEART":
      cardSVG = <TENHEART className="Card-body" />;
      break;
    case "10DIAMOND":
      cardSVG = <TENDIAMOND className="Card-body" />;
      break;
    case "11CLUB":
      cardSVG = <img style={{ height: "inherit" }} src={JACKCLUB} alt="card" />;
      break;
    case "11SPADE":
      cardSVG = (
        <img style={{ height: "inherit" }} src={JACKSPADE} alt="card" />
      );
      break;
    case "11HEART":
      cardSVG = (
        <img style={{ height: "inherit" }} src={JACKHEART} alt="card" />
      );
      break;
    case "11DIAMOND":
      cardSVG = (
        <img style={{ height: "inherit" }} src={JACKDIAMOND} alt="card" />
      );
      break;
    case "12CLUB":
      cardSVG = (
        <img style={{ height: "inherit" }} src={QUEENCLUB} alt="Logo" />
      );
      break;
    case "12SPADE":
      cardSVG = (
        <img style={{ height: "inherit" }} src={QUEENSPADE} alt="card" />
      );
      break;
    case "12HEART":
      cardSVG = (
        <img style={{ height: "inherit" }} src={QUEENHEART} alt="card" />
      );
      break;
    case "12DIAMOND":
      cardSVG = (
        <img style={{ height: "inherit" }} src={QUEENDIAMOND} alt="card" />
      );
      break;
    case "13CLUB":
      cardSVG = <img style={{ height: "inherit" }} src={KINGCLUB} alt="card" />;
      break;
    case "13SPADE":
      cardSVG = (
        <img style={{ height: "inherit" }} src={KINGSPADE} alt="card" />
      );
      break;
    case "13HEART":
      cardSVG = (
        <img style={{ height: "inherit" }} src={KINGHEART} alt="card" />
      );
      break;
    case "13DIAMOND":
      cardSVG = (
        <img style={{ height: "inherit" }} src={KINGDIAMOND} alt="card" />
      );
      break;
    case "1JOKER":
      cardSVG = <img style={{ height: "inherit" }} src={BIGJOKER} alt="card" />;
      break;
    case "2JOKER":
      cardSVG = (
        <img style={{ height: "inherit" }} src={LITTLEJOKER} alt="card" />
      );
      break;
    default:
      console.log(`invalid card: ${JSON.stringify(props.card)}`);
  }

  return <>{cardSVG}</>;
};

export default Card;
