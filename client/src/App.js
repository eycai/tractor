import React, { Component } from "react";
import "./App.css";
import socketIOClient from "socket.io-client";
import Landing from "./components/pages/Landing";
import Lobby from "./components/pages/Lobby";
import Room from "./Room";
import { Router, Link } from "@reach/router";
import { get, post } from "./api/fetch";
import { pages } from "./utilities.js";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      user: null
    };
  }

  componentDidMount() {
    const { endpoint } = this.state;
    const socket = socketIOClient(endpoint);
    get("/whoami").then(res => {
      console.log(res);
      if (res.id) {
        this.setState({ user: res });
      }
    });
    socket.on("connect", data => {
      // this.setState({ response: JSON.stringify(data) })
      console.log("connected with socket");
      console.log(data);
      console.log(socket.id);
      if (this.state.user) {
        post("/connect", { socketId: socket.id, userId: this.state.user.id });
      }
    });
  }

  render() {
    return (
      <Router>
        <Landing path="/" user={this.state.user} />
        <Room path="/:roomid" user={this.state.user} />
      </Router>
    );
  }
}

export default App;
