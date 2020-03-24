import React, { Component } from "react";
import "./App.css";
import socketIOClient from "socket.io-client";
import Landing from "./components/pages/Landing";
import Lobby from "./components/pages/Lobby";
import LobbyList from "./components/pages/LobbyList";
import { get, post } from "./api/fetch";
import { pages } from "./utilities.js";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      user: null,
      page: pages.LANDING
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

  setPage = page => {
    this.setState({ page: page });
  };

  render() {
    const { response } = this.state;
    switch (this.state.page) {
      case pages.LANDING:
        return <Landing setPage={this.setPage} user={this.state.user} />;
      case pages.LOBBYLIST:
        return <LobbyList />;
      default:
        return <div>Oops</div>;
    }
  }
}

export default App;
