import React, { Component } from "react";
import "./App.css";
import Landing from "./components/pages/Landing";
import Lobby from "./components/pages/Lobby";
import Room from "./Room";
import { Router, Link } from "@reach/router";
import { get, post } from "./api/fetch";
import { pages } from "./utilities.js";
import { socket } from "./client-socket";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      user: null
    };
  }

  setUser = newUser => {
    this.setState({ user: newUser });
  };

  componentDidMount() {
    get("/whoami").then(res => {
      if (res.status === 200) {
        this.setState({ user: res.payload });
      }
    });
  }

  render() {
    return (
      <Router>
        <Landing path="/" user={this.state.user} />
        <Room path="/:roomid" user={this.state.user} setUser={this.setUser} />
      </Router>
    );
  }
}

export default App;
