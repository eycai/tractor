export const pages = {
  LANDING: 1,
  LOBBYLIST: 2,
  LOBBY: 3,
  GAME: 4,
  ERROR: 5
};

export const testData = {
  id: "JJFK",
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
        cardsPlayed: [{ value: 1, suit: "CLUB" }]
      },
      emily: {
        username: "emily",
        team: "PEASANTS",
        level: 1,
        points: 15,
        cardsPlayed: [
          { value: 1, suit: "CLUB" },
          { value: 1, suit: "HEART" },
          { value: 1, suit: "CLUB" },
          { value: 1, suit: "HEART" },
          { value: 1, suit: "CLUB" },
          { value: 1, suit: "HEART" }
        ]
      },
      edward: {
        username: "edward",
        team: "PEASANTS",
        level: 1,
        points: 15,
        cardsPlayed: [
          { value: 1, suit: "CLUB" },
          { value: 1, suit: "HEART" },
          { value: 1, suit: "CLUB" },
          { value: 1, suit: "HEART" },
          { value: 1, suit: "CLUB" },
          { value: 1, suit: "HEART" }
        ]
      },
      yolanda: {
        username: "yolanda",
        team: "BOSSES",
        level: 1,
        points: 15,
        cardsPlayed: [
          { value: 2, suit: "SPADE" },
          { value: 1, suit: "DIAMOND" }
        ]
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
