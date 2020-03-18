import React, { Component } from 'react';
import './App.css';
import API from './api/api.js';

class App extends Component {
  constructor() {
    super();
    this.state = {
      pong: 'pending'
    };
  }
  componentWillMount() {
    API.getHelloWorld()
      .then(response => {
        console.log('HELP');
        this.setState(() => {
          return { pong: response.help };
        });
      })
      .catch(function(error) {
        console.log(error);
      });
  }

  render() {
    return <h1>Ping {this.state.pong}</h1>;
  }
}

export default App;
