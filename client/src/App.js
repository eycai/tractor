import React, { Component } from "react";
import "./App.css";
import socketIOClient from "socket.io-client";
import Landing from "./components/pages/Landing";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      response: false,
      endpoint: "http://localhost:8080/"
    };
  }

  componentDidMount() {
    const { endpoint } = this.state;
    const socket = socketIOClient(endpoint);
    socket.on("connect", data =>
      this.setState({ response: JSON.stringify(data) })
    );
  }

  render() {
    const { response } = this.state;
    return <Landing />;
  }
}

export default App;
