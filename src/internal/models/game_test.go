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

func TestPlayCards(t *testing.T) {
	type plays struct {
		user   string
		play   [][]models.Card
		status models.TrickStatus
	}
	type test struct {
		plays          []plays
		hands          map[string][]models.Card
		game           *models.Game
		firstInTrick   string
		drawOrder      []string
		expectedWinner string
		expectedPlays  []plays
		expectedStatus models.TrickStatus
		expectedPoints map[string]int
		expectedTurn   string
	}

	tests := []test{
		{
			plays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
						},
					},
					user: "a",
				},
			},
			hands: map[string][]models.Card{
				"a": {
					{Value: 5, Suit: models.Spade},
				},
				"b": {
					{Value: 7, Suit: models.Spade},
				},
				"c": {
					{Value: 8, Suit: models.Spade},
				},
				"d": {
					{Value: 10, Suit: models.Spade},
				},
			},
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": {},
					"b": {},
					"c": {},
					"d": {},
				},
				Turn:        "a",
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstInTrick:   "a",
			drawOrder:      []string{"b", "c", "d", "a"},
			expectedWinner: "a",
			expectedPlays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
						},
					},
					user:   "a",
					status: models.PlayingTrick,
				},
			},
			expectedPoints: map[string]int{"a": 0, "b": 0, "c": 0, "d": 0},
			expectedTurn:   "b",
		},
		{
			plays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
							{Value: 5, Suit: models.Spade},
						},
						{
							{Value: 1, Suit: models.Spade},
						},
					},
					user: "a",
				},
			},
			hands: map[string][]models.Card{
				"a": {
					{Value: 5, Suit: models.Spade},
					{Value: 5, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
				},
				"b": {
					{Value: 7, Suit: models.Spade},
					{Value: 7, Suit: models.Spade},
				},
				"c": {
					{Value: 8, Suit: models.Spade},
					{Value: 6, Suit: models.Spade},
				},
				"d": {
					{Value: 10, Suit: models.Spade},
					{Value: 4, Suit: models.Spade},
				},
			},
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": {},
					"b": {},
					"c": {},
					"d": {},
				},
				Turn:        "a",
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstInTrick:   "a",
			drawOrder:      []string{"b", "c", "d", "a"},
			expectedWinner: "a",
			expectedPlays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
							{Value: 5, Suit: models.Spade},
						},
					},
					user:   "a",
					status: models.PlayingTrick,
				},
			},
			expectedPoints: map[string]int{"a": 0, "b": 0, "c": 0, "d": 0},
			expectedTurn:   "b",
		},
		{
			plays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
							{Value: 5, Suit: models.Spade},
						},
						{
							{Value: 1, Suit: models.Spade},
						},
					},
					user: "a",
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Spade},
							{Value: 7, Suit: models.Spade},
						},
					},
					user: "b",
				},
			},
			hands: map[string][]models.Card{
				"a": {
					{Value: 5, Suit: models.Spade},
					{Value: 5, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
				},
				"b": {
					{Value: 7, Suit: models.Spade},
					{Value: 7, Suit: models.Spade},
				},
				"c": {
					{Value: 8, Suit: models.Spade},
					{Value: 6, Suit: models.Spade},
				},
				"d": {
					{Value: 10, Suit: models.Spade},
					{Value: 4, Suit: models.Spade},
				},
			},
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": {},
					"b": {},
					"c": {},
					"d": {},
				},
				Turn:        "a",
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstInTrick:   "a",
			drawOrder:      []string{"b", "c", "d", "a"},
			expectedWinner: "b",
			expectedPlays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
							{Value: 5, Suit: models.Spade},
						},
					},
					user:   "a",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Spade},
							{Value: 7, Suit: models.Spade},
						},
					},
					user:   "b",
					status: models.PlayingTrick,
				},
			},
			expectedPoints: map[string]int{"a": 0, "b": 0, "c": 0, "d": 0},
			expectedTurn:   "c",
		}, {
			plays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 1, Suit: models.Spade},
						},
					},
					user: "a",
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Spade},
						},
					},
					user: "b",
				},
			},
			hands: map[string][]models.Card{
				"a": {
					{Value: 5, Suit: models.Spade},
					{Value: 5, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
				},
				"b": {
					{Value: 7, Suit: models.Spade},
					{Value: 7, Suit: models.Spade},
				},
				"c": {
					{Value: 8, Suit: models.Spade},
					{Value: 6, Suit: models.Spade},
				},
				"d": {
					{Value: 10, Suit: models.Spade},
					{Value: 4, Suit: models.Spade},
				},
			},
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": {},
					"b": {},
					"c": {},
					"d": {},
				},
				Turn:        "a",
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstInTrick:   "a",
			drawOrder:      []string{"b", "c", "d", "a"},
			expectedWinner: "a",
			expectedPlays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 1, Suit: models.Spade},
						},
					},
					user:   "a",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Spade},
						},
					},
					user:   "b",
					status: models.PlayingTrick,
				},
			},
			expectedPoints: map[string]int{"a": 0, "b": 0, "c": 0, "d": 0},
			expectedTurn:   "c",
		}, {
			plays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
						},
					},
					user: "a",
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Diamond},
						},
					},
					user: "b",
				},
				{
					play: [][]models.Card{
						{
							{Value: 8, Suit: models.Spade},
						},
					},
					user: "c",
				},
			},
			hands: map[string][]models.Card{
				"a": {
					{Value: 5, Suit: models.Spade},
					{Value: 5, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
				},
				"b": {
					{Value: 7, Suit: models.Diamond},
					{Value: 7, Suit: models.Diamond},
				},
				"c": {
					{Value: 8, Suit: models.Spade},
					{Value: 6, Suit: models.Spade},
				},
				"d": {
					{Value: 10, Suit: models.Spade},
					{Value: 4, Suit: models.Spade},
				},
			},
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": {},
					"b": {},
					"c": {},
					"d": {},
				},
				Turn:        "a",
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstInTrick:   "a",
			drawOrder:      []string{"b", "c", "d", "a"},
			expectedWinner: "b",
			expectedPlays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
						},
					},
					user:   "a",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Diamond},
						},
					},
					user:   "b",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 8, Suit: models.Spade},
						},
					},
					user:   "c",
					status: models.PlayingTrick,
				},
			},
			expectedPoints: map[string]int{"a": 0, "b": 0, "c": 0, "d": 0},
			expectedTurn:   "d",
		},
		{
			plays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
						},
					},
					user: "a",
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Diamond},
						},
					},
					user: "b",
				},
				{
					play: [][]models.Card{
						{
							{Value: 8, Suit: models.Spade},
						},
					},
					user: "c",
				},
				{
					play: [][]models.Card{
						{
							{Value: 10, Suit: models.Diamond},
						},
					},
					user: "d",
				},
			},
			hands: map[string][]models.Card{
				"a": {
					{Value: 5, Suit: models.Spade},
					{Value: 5, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
				},
				"b": {
					{Value: 7, Suit: models.Diamond},
					{Value: 7, Suit: models.Diamond},
				},
				"c": {
					{Value: 8, Suit: models.Spade},
					{Value: 6, Suit: models.Spade},
				},
				"d": {
					{Value: 10, Suit: models.Diamond},
				},
			},
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": {},
					"b": {},
					"c": {},
					"d": {},
				},
				Turn:        "a",
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstInTrick:   "a",
			drawOrder:      []string{"b", "c", "d", "a"},
			expectedWinner: "",
			expectedPlays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
						},
					},
					user:   "a",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Diamond},
						},
					},
					user:   "b",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 8, Suit: models.Spade},
						},
					},
					user:   "c",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 10, Suit: models.Diamond},
						},
					},
					user:   "d",
					status: models.TrickEnded,
				},
			},
			expectedPoints: map[string]int{"a": 0, "b": 0, "c": 0, "d": 15},
			expectedTurn:   "d",
		}, {
			plays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
						},
					},
					user: "a",
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Diamond},
						},
					},
					user: "b",
				},
				{
					play: [][]models.Card{
						{
							{Value: 8, Suit: models.Spade},
						},
					},
					user: "c",
				},
				{
					play: [][]models.Card{
						{
							{Value: 10, Suit: models.Spade},
						},
					},
					user: "d",
				},
			},
			hands: map[string][]models.Card{
				"a": {
					{Value: 5, Suit: models.Spade},
					{Value: 5, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
					{Value: 1, Suit: models.Spade},
				},
				"b": {
					{Value: 7, Suit: models.Diamond},
					{Value: 7, Suit: models.Diamond},
				},
				"c": {
					{Value: 8, Suit: models.Spade},
					{Value: 6, Suit: models.Spade},
				},
				"d": {
					{Value: 10, Suit: models.Spade},
				},
			},
			game: &models.Game{
				Players: map[string]*models.Player{
					"a": {},
					"b": {},
					"c": {},
					"d": {},
				},
				Turn:        "a",
				TrumpSuit:   models.Diamond,
				TrumpNumber: 2,
			},
			firstInTrick:   "a",
			drawOrder:      []string{"b", "c", "d", "a"},
			expectedWinner: "",
			expectedPlays: []plays{
				{
					play: [][]models.Card{
						{
							{Value: 5, Suit: models.Spade},
						},
					},
					user:   "a",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 7, Suit: models.Diamond},
						},
					},
					user:   "b",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 8, Suit: models.Spade},
						},
					},
					user:   "c",
					status: models.PlayingTrick,
				},
				{
					play: [][]models.Card{
						{
							{Value: 10, Suit: models.Spade},
						},
					},
					user:   "d",
					status: models.TrickEnded,
				},
			},
			expectedPoints: map[string]int{"a": 0, "b": 15, "c": 0, "d": 0},
			expectedTurn:   "b",
		},
	}

	for _, tc := range tests {
		tc.game.SetDrawOrder(tc.drawOrder)
		tc.game.SetTrickStarter(tc.firstInTrick)
		for i, play := range tc.plays {
			otherHands := [][]models.Card{}
			for u, h := range tc.hands {
				if u != play.user {
					otherHands = append(otherHands, tc.game.GetUpdatedCards(h))
				}
			}
			status, played, err := tc.game.PlayCards(play.user, play.play, otherHands)
			if err != nil {
				t.Errorf("expected no error, but got error")
			}
			if status != tc.expectedPlays[i].status {
				t.Errorf("should have gotten status %v, instead got %v", tc.expectedPlays[i].status, status)
			}
			if !playsEqual(played, tc.expectedPlays[i].play) {
				t.Errorf("expcted to return play %v, but got %v", tc.expectedPlays[i].play, played)
			}
		}

		if tc.expectedWinner != tc.game.GetCurrentWinner() {
			t.Errorf("expected winner %s but got %s", tc.expectedWinner, tc.game.GetCurrentWinner())
		}

		for u, p := range tc.expectedPoints {
			if tc.game.Players[u].Points != p {
				t.Errorf("expected %d points for %s but got %d", p, u, tc.game.Players[u].Points)
			}
		}

		if tc.game.Turn != tc.expectedTurn {
			t.Errorf("expected turn to be %s, instead got %s", tc.expectedTurn, tc.game.Turn)
		}
	}
}
