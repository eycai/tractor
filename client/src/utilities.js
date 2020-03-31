export const pages = {
  LANDING: 1,
  LOBBYLIST: 2,
  LOBBY: 3,
  GAME: 4,
  ERROR: 5
};

export const testData = {
  id: "UKFY",
  users: [
    { username: "alex", connected: true },
    { username: "emily", connected: true },
    { username: "yolanda", connected: true },
    { username: "edward", connected: true }
  ],
  host: "alex",
  name: "name",
  capacity: 8,
  game: {
    players: {
      alex: {
        username: "alex",
        team: "BOSSES",
        level: 1,
        points: 15,
        cardsPlayed: [{ value: 10, suit: "CLUBS" }]
      },
      emily: {
        username: "emily",
        team: "PEASANTS",
        level: 1,
        points: 15,
        cardsPlayed: []
      },
      edward: {
        username: "edward",
        team: "PEASANTS",
        level: 1,
        points: 15,
        cardsPlayed: []
      },
      yolanda: {
        username: "yolanda",
        team: "BOSSES",
        level: 1,
        points: 15,
        cardsPlayed: []
      }
    },
    turn: "emily",
    trumpSuit: "DIAMOND",
    trumpNumber: 2,
    trumpFlipUser: "emily",
    trumpNumCardsFlipped: 2,
    banker: "edward",
    gamePhase: "PLAYING"
  }
};
