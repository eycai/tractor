import React, { Component } from 'react';
import './App.css';
// import API from './api/api';
import { connect, sendMsg } from './api/websocket';
import Header from './components/Header';
import CardTable from './components/CardTable';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      pong: 'pending',
      cardTable: []
    };
    connect();
  }

  componentDidMount() {
    connect(msg => {
      console.log('New Message');
      this.setState(prevState => ({
        cardTable: [...this.state.cardTable, msg]
      }));
      console.log(this.state);
    });
  }

  send() {
    console.log('hello');
    sendMsg('hello');
  }

  render() {
    return (
      <div className="App">
        <Header />
        <CardTable cardTable={this.state.cardTable} />
        <button onClick={this.send}>Hit</button>
      </div>
    );
  }
}

export default App;
