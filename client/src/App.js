import React, { Component } from 'react';
import './App.css';
import Header from './components/Header';
import { Pane } from 'evergreen-ui';
import socketIOClient from 'socket.io-client';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      response: false,
      endpoint: 'http://localhost:8080/'
    };
  }

  componentDidMount() {
    const { endpoint } = this.state;
    const socket = socketIOClient(endpoint);
    socket.on('reply', data => this.setState({ response: data }));
  }

  render() {
    const { response } = this.state;
    return (
      <Pane>
        <Header />
        {response ? <p>Test: {response}</p> : <p>Loading...</p>}
      </Pane>
    );
  }
}

export default App;
