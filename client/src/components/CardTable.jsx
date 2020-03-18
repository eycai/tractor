import React, { Component } from 'react';

class CardTable extends Component {
  render() {
    const messages = this.props.cardTable.map((msg, index) => (
      <p key={index}>{msg.data}</p>
    ));

    return (
      <div className="CardTable">
        <h2>Card Table</h2>
        {messages}
      </div>
    );
  }
}

export default CardTable;
