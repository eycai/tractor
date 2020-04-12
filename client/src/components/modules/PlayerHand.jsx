import React from "react";

import "./PlayerHand.css";
import Card from "./Card";
import { cardHeight, cardWidth, offset } from "./PlayerCards";

let PlayerHand = (props) => {
  const initialOffset = (offset * (props.user.hand.length - 1) + cardWidth) / 2;

  return (
    <div className="PlayerHand-container">
      <div className="PlayerHand-body">
        {props.user.hand.map((c, i) => (
          <div
            style={{
              height: cardHeight,
              transform: `translateY(-100%) translateX(${i * offset -
                initialOffset}px)`,
            }}
            className="PlayerHand-card-container"
          >
            <div
              className={`PlayerHand-card-body ${
                props.selectedCards.includes(c) ? "PlayerHand-selected" : ""
              }`}
              onClick={() => {
                props.selectedCards.includes(c)
                  ? props.setSelectedCards((cards) =>
                      cards.filter((e) => e !== c)
                    )
                  : props.setSelectedCards((cards) => cards.concat(c));
              }}
            >
              <Card card={c} />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default PlayerHand;
