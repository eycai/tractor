package models_test

import (
	"testing"

	"github.com/eycai/tractor/src/internal/models"
)

func testValues(d models.Deck) bool {
	values := make([]int, 13)
	suits := make(map[models.Suit]int)
	for _, c := range d.Cards {
		if c.Suit != models.Joker {
			values[c.Value-1]++
		}
		suits[c.Suit]++
	}
	numDecks := len(d.Cards) / 54
	for _, v := range values {
		if v != numDecks*4 {
			return false
		}
	}
	for s, n := range suits {
		if s == models.Joker {
			if n != 2*numDecks {
				return false
			}
		} else {
			if n != 13*numDecks {
				return false
			}
		}
	}
	return true
}

func TestKittySize(t *testing.T) {
	type test struct {
		input    models.Game
		expected int
	}

	tests := []test{
		{
			input: models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
			},
			expected: 8,
		}, {
			input: models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
					"e": &models.Player{},
				},
			},
			expected: 8,
		}, {
			input: models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
					"e": &models.Player{},
					"f": &models.Player{},
				},
			},
			expected: 6,
		},
	}

	for _, tc := range tests {
		kitty := tc.input.KittySize()
		if kitty != tc.expected {
			t.Errorf("expected kitty of size %d, got %d", tc.expected, kitty)
		}
	}
}

func TestDeck(t *testing.T) {
	g := models.Game{
		Players: map[string]*models.Player{
			"a": &models.Player{},
			"b": &models.Player{},
			"c": &models.Player{},
			"d": &models.Player{},
		},
	}
	d := g.GetDeck()
	if len(d.Cards) != 108 {
		t.Errorf("expected deck of length 108, got %d", len(d.Cards))
	}
	if !testValues(d) {
		t.Errorf("incorrect value assignments")
	}

	g = models.Game{
		Players: map[string]*models.Player{
			"a": &models.Player{},
			"b": &models.Player{},
			"c": &models.Player{},
			"d": &models.Player{},
			"e": &models.Player{},
		},
	}

	d = g.GetDeck()
	if len(d.Cards) != 108 {
		t.Errorf("expected deck of length 108, got %d", len(d.Cards))
	}
	if !testValues(d) {
		t.Errorf("incorrect value assignments")
	}

	g = models.Game{
		Players: map[string]*models.Player{
			"a": &models.Player{},
			"b": &models.Player{},
			"c": &models.Player{},
			"d": &models.Player{},
			"e": &models.Player{},
			"f": &models.Player{},
		},
	}

	d = g.GetDeck()
	if len(d.Cards) != 162 {
		t.Errorf("expected deck of length 162, got %d", len(d.Cards))
	}
	if !testValues(d) {
		t.Errorf("incorrect value assignments")
	}
}

func gamesEquivalentForFlip(a *models.Game, b *models.Game) bool {
	return a.TrumpNumber == b.TrumpNumber &&
		a.TrumpSuit == b.TrumpSuit &&
		a.TrumpFlipUser == b.TrumpFlipUser &&
		a.TrumpNumCardsFlipped == b.TrumpNumCardsFlipped &&
		a.Banker == b.Banker
}

func TestFlipCard(t *testing.T) {
	type test struct {
		game     *models.Game
		flipCard models.Card
		flipNum  int
		flipUser string
		expected *models.Game
		success  bool
	}

	tests := []test{
		{
			game: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Club,
				TrumpFlipUser:        "",
				TrumpNumCardsFlipped: 0,
				GamePhase:            models.Drawing,
			},
			flipCard: models.Card{
				Value: 2,
				Suit:  models.Diamond,
			},
			flipUser: "a",
			flipNum:  1,
			expected: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Diamond,
				TrumpFlipUser:        "a",
				TrumpNumCardsFlipped: 1,
				GamePhase:            models.Drawing,
				Banker:               "a",
			},
			success: true,
		}, {
			game: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Club,
				TrumpFlipUser:        "b",
				TrumpNumCardsFlipped: 1,
				GamePhase:            models.Drawing,
				Banker:               "b",
			},
			flipCard: models.Card{
				Value: 2,
				Suit:  models.Diamond,
			},
			flipUser: "a",
			flipNum:  1,
			expected: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Club,
				TrumpFlipUser:        "b",
				TrumpNumCardsFlipped: 1,
				GamePhase:            models.Drawing,
				Banker:               "b",
			},
			success: false,
		}, {
			game: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Club,
				TrumpFlipUser:        "a",
				TrumpNumCardsFlipped: 1,
				GamePhase:            models.Drawing,
				Banker:               "a",
			},
			flipCard: models.Card{
				Value: 2,
				Suit:  models.Club,
			},
			flipUser: "a",
			flipNum:  1,
			expected: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Club,
				TrumpFlipUser:        "a",
				TrumpNumCardsFlipped: 2,
				GamePhase:            models.Drawing,
				Banker:               "a",
			},
			success: true,
		}, {
			game: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Club,
				TrumpFlipUser:        "b",
				TrumpNumCardsFlipped: 1,
				GamePhase:            models.Drawing,
				Banker:               "b",
			},
			flipCard: models.Card{
				Value: 2,
				Suit:  models.Diamond,
			},
			flipUser: "a",
			flipNum:  2,
			expected: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Diamond,
				TrumpFlipUser:        "a",
				TrumpNumCardsFlipped: 2,
				GamePhase:            models.Drawing,
				Banker:               "a",
			},
			success: true,
		}, {
			game: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Club,
				TrumpFlipUser:        "b",
				TrumpNumCardsFlipped: 2,
				GamePhase:            models.Drawing,
				Banker:               "b",
			},
			flipCard: models.Card{
				Value: 1,
				Suit:  models.Joker,
			},
			flipUser: "a",
			flipNum:  2,
			expected: &models.Game{
				TrumpNumber:          2,
				TrumpSuit:            models.Joker,
				TrumpFlipUser:        "a",
				TrumpNumCardsFlipped: 2,
				GamePhase:            models.Drawing,
				Banker:               "a",
			},
			success: true,
		},
	}

	for _, tc := range tests {
		success := tc.game.FlipCard(tc.flipCard, tc.flipNum, tc.flipUser)
		if success != tc.success {
			t.Errorf("expected success to be %v, instead got %v", tc.success, success)
		}
		if !gamesEquivalentForFlip(tc.expected, tc.game) {
			t.Errorf("expected game %v, instead got %v", tc.expected, tc.game)
		}
	}
}

func TestIsValidPlayForGame(t *testing.T) {
	type test struct {
		game        *models.Game
		firstPlayer string
		cards       [][]models.Card
		hand        []models.Card
		success     bool
	}

	tests := []test{
		{
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
				TrumpSuit:   models.Club,
				TrumpNumber: 2,
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Spade,
				},
			},
			success: true,
		}, {
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{
						CardsPlayed: [][]models.Card{
							{
								{
									Value: 1,
									Suit:  models.Diamond,
								},
							},
						},
					},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
			},
			success: true,
		}, {
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{
						CardsPlayed: [][]models.Card{
							{
								{
									Value: 1,
									Suit:  models.Diamond,
								},
							},
						},
					},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Spade,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 3,
					Suit:  models.Spade,
				},
			},
			success: false,
		}, {
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{
						CardsPlayed: [][]models.Card{
							{
								{
									Value: 1,
									Suit:  models.Club,
								},
							},
						},
					},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 3,
					Suit:  models.Spade,
				},
			},
			success: true,
		}, {
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{
						CardsPlayed: [][]models.Card{
							{
								{
									Value: 1,
									Suit:  models.Club,
								},
							},
						},
					},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 3,
					Suit:  models.Club,
				},
			},
			success: false,
		}, {
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{
						CardsPlayed: [][]models.Card{
							{
								{
									Value: 1,
									Suit:  models.Club,
								},
							},
						},
					},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
				TrumpSuit:   models.Club,
				TrumpNumber: 2,
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Spade,
				},
			},
			success: false,
		}, {
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				}, {
					{
						Value: 1,
						Suit:  models.Diamond,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 1,
					Suit:  models.Diamond,
				},
			},
			success: false,
		}, {
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
				TrumpSuit:   models.Club,
				TrumpNumber: 2,
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				}, {
					{
						Value: 1,
						Suit:  models.Diamond,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 1,
					Suit:  models.Diamond,
				},
			},
			success: true,
		}, {
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
				TrumpSuit:   models.Club,
				TrumpNumber: 2,
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				}, {
					{
						Value: 2,
						Suit:  models.Diamond,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Diamond,
				},
			},
			success: false,
		}, {
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": &models.Player{},
					"b": &models.Player{},
					"c": &models.Player{},
					"d": &models.Player{},
				},
				TrumpSuit:   models.Spade,
				TrumpNumber: 2,
			},
			firstPlayer: "a",
			cards: [][]models.Card{
				{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				}, {
					{
						Value: 1,
						Suit:  models.Club,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 1,
					Suit:  models.Club,
				},
			},
			success: false,
		},
	}

	for _, tc := range tests {
		tc.game.SetTrickStarter(tc.firstPlayer)
		tc.hand = tc.game.GetUpdatedCards(tc.hand)
		success := tc.game.IsValidPlayForGame(tc.cards, tc.hand)
		if success != tc.success {
			t.Errorf("expected play %v, %v to have success %v, but instead got %v", tc.cards, tc.hand, tc.success, success)
		}
	}
}
