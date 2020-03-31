package models_test

import (
	"fmt"
	"testing"

	"github.com/eycai/tractor/src/internal/models"
)

func TestIsValidPlay(t *testing.T) {
	type test struct {
		prev        [][]models.Card
		next        [][]models.Card
		hand        []models.Card
		trumpSuit   models.Suit
		trumpNumber int
		expected    bool
	}

	tests := []test{
		{
			prev: [][]models.Card{
				[]models.Card{
					{
						Value: 5,
						Suit:  models.Diamond,
					},
					{
						Value: 5,
						Suit:  models.Diamond,
					},
				}, []models.Card{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				},
			},
			next: [][]models.Card{
				[]models.Card{
					{
						Value: 2,
						Suit:  models.Diamond,
					},
					{
						Value: 2,
						Suit:  models.Diamond,
					},
				}, []models.Card{
					{
						Value: 1,
						Suit:  models.Diamond,
					},
				}, []models.Card{
					{
						Value: 6,
						Suit:  models.Spade,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 6,
					Suit:  models.Spade,
				},
				{
					Value: 1,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Diamond,
				},
			},
			trumpNumber: 10,
			trumpSuit:   models.Club,
			expected:    true,
		}, {
			prev: [][]models.Card{
				[]models.Card{
					{
						Value: 5,
						Suit:  models.Diamond,
					},
					{
						Value: 5,
						Suit:  models.Diamond,
					},
				}, []models.Card{
					{
						Value: 3,
						Suit:  models.Diamond,
					},
					{
						Value: 3,
						Suit:  models.Diamond,
					},
				},
			},
			next: [][]models.Card{
				[]models.Card{
					{
						Value: 2,
						Suit:  models.Diamond,
					},
					{
						Value: 2,
						Suit:  models.Diamond,
					},
				}, []models.Card{
					{
						Value: 1,
						Suit:  models.Diamond,
					},
				}, []models.Card{
					{
						Value: 6,
						Suit:  models.Spade,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 6,
					Suit:  models.Diamond,
				},
				{
					Value: 1,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Diamond,
				},
			},
			trumpNumber: 10,
			trumpSuit:   models.Club,
			expected:    false,
		}, {
			prev: [][]models.Card{
				[]models.Card{
					{
						Value: 5,
						Suit:  models.Diamond,
					},
					{
						Value: 5,
						Suit:  models.Diamond,
					},
				},
			},
			next: [][]models.Card{
				[]models.Card{
					{
						Value: 2,
						Suit:  models.Spade,
					},
				}, []models.Card{
					{
						Value: 5,
						Suit:  models.Spade,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 2,
					Suit:  models.Spade,
				},
				{
					Value: 5,
					Suit:  models.Spade,
				},
				{
					Value: 7,
					Suit:  models.Spade,
				},
			},
			trumpNumber: 2,
			trumpSuit:   models.Diamond,
			expected:    true,
		}, {
			trumpNumber: 2,
			trumpSuit:   models.Diamond,
			prev: [][]models.Card{
				[]models.Card{
					{
						Value: 5,
						Suit:  models.Diamond,
					},
					{
						Value: 5,
						Suit:  models.Diamond,
					},
				},
			},
			next: [][]models.Card{
				[]models.Card{
					{
						Value: 2,
						Suit:  models.Spade,
					},
				}, []models.Card{
					{
						Value: 5,
						Suit:  models.Spade,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 2,
					Suit:  models.Spade,
				},
				{
					Value: 8,
					Suit:  models.Diamond,
				},
				{
					Value: 7,
					Suit:  models.Spade,
				},
			},
			expected: false,
		}, {
			trumpNumber: 3,
			trumpSuit:   models.Diamond,
			prev: [][]models.Card{
				[]models.Card{
					{
						Value: 5,
						Suit:  models.Diamond,
					},
					{
						Value: 5,
						Suit:  models.Diamond,
					},
				}, []models.Card{
					{
						Value: 7,
						Suit:  models.Diamond,
					},
				},
			},
			next: [][]models.Card{
				[]models.Card{
					{
						Value: 3,
						Suit:  models.Spade,
					},
				}, []models.Card{
					{
						Value: 6,
						Suit:  models.Diamond,
					},
				},
			},
			hand: []models.Card{
				{
					Value: 3,
					Suit:  models.Spade,
				},
				{
					Value: 6,
					Suit:  models.Diamond,
				},
				{
					Value: 7,
					Suit:  models.Spade,
				},
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		for j, tr := range tc.prev {
			for k, c := range tr {
				tc.prev[j][k] = c.WithTrump(tc.trumpNumber, tc.trumpSuit)
			}
		}

		for j, tr := range tc.next {
			for k, c := range tr {
				tc.next[j][k] = c.WithTrump(tc.trumpNumber, tc.trumpSuit)
			}
		}
		for i, c := range tc.hand {
			tc.hand[i] = c.WithTrump(tc.trumpNumber, tc.trumpSuit)
		}

		valid := models.IsValidPlay(tc.prev, tc.next, tc.hand)
		if valid != tc.expected {
			t.Errorf("expected %v but got %v for tricks %v, %v", tc.expected, valid, tc.prev, tc.next)
		}
	}
}

func TestNextTrickWins(t *testing.T) {
	type test struct {
		prev        []models.Trick
		next        []models.Trick
		expected    bool
		trumpNumber int
		trumpSuit   models.Suit
	}
	tests := []test{
		{
			prev: []models.Trick{
				{
					Pattern:  models.NOfAKind,
					NumCards: 1,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value: 5,
						Suit:  models.Diamond,
					},
				},
			},
			next: []models.Trick{
				{
					Pattern:  models.NOfAKind,
					NumCards: 1,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value: 8,
						Suit:  models.Diamond,
					},
				},
			},
			expected:    true,
			trumpNumber: 2,
			trumpSuit:   models.Club,
		}, {
			prev: []models.Trick{
				{
					Pattern:  models.NOfAKind,
					NumCards: 1,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value: 5,
						Suit:  models.Diamond,
					},
				},
			},
			next: []models.Trick{
				{
					Pattern:  models.NOfAKind,
					NumCards: 1,
					Suit:     models.Spade,
					IsTrump:  false,
					LargestCard: models.Card{
						Value: 8,
						Suit:  models.Spade,
					},
				},
			},
			trumpNumber: 10,
			trumpSuit:   models.Heart,
			expected:    false,
		}, {
			prev: []models.Trick{
				{
					Pattern:  models.NOfAKind,
					NumCards: 2,
					Suit:     models.Diamond,
					LargestCard: models.Card{
						Value: 5,
						Suit:  models.Diamond,
					},
				},
			},
			next: []models.Trick{
				{
					Pattern:  models.NOfAKind,
					NumCards: 2,
					Suit:     models.Spade,
					LargestCard: models.Card{
						Value: 2,
						Suit:  models.Spade,
					},
				},
			},
			trumpNumber: 5,
			trumpSuit:   models.Heart,
			expected:    false,
		}, {
			prev: []models.Trick{
				{
					Pattern:               models.Tractor,
					NumCards:              4,
					TractorNumConsecutive: 2,
					Suit:                  models.Diamond,
					IsTrump:               false,
					LargestCard: models.Card{
						Value: 5,
						Suit:  models.Diamond,
					},
				},
			},
			next: []models.Trick{
				{
					Pattern:  models.NOfAKind,
					NumCards: 2,
					Suit:     models.Spade,
					IsTrump:  true,
					LargestCard: models.Card{
						Value: 2,
						Suit:  models.Spade,
					},
				},
			},
			trumpNumber: 10,
			trumpSuit:   models.Spade,
			expected:    false,
		}, {
			prev: []models.Trick{
				{
					Pattern:               models.Tractor,
					NumCards:              4,
					TractorNumConsecutive: 2,
					Suit:                  models.Diamond,
					IsTrump:               false,
					LargestCard: models.Card{
						Value: 5,
						Suit:  models.Diamond,
					},
				}, {
					Pattern:  models.NOfAKind,
					NumCards: 2,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value: 10,
						Suit:  models.Diamond,
					},
				},
			},
			next: []models.Trick{
				{
					Pattern:               models.Tractor,
					NumCards:              4,
					TractorNumConsecutive: 2,
					Suit:                  models.Diamond,
					IsTrump:               false,
					LargestCard: models.Card{
						Value: 7,
						Suit:  models.Diamond,
					},
				}, {
					Pattern:  models.NOfAKind,
					NumCards: 2,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value: 2,
						Suit:  models.Diamond,
					},
				},
			},
			trumpNumber: 10,
			trumpSuit:   models.Heart,
			expected:    false,
		}, {
			prev: []models.Trick{
				{
					Pattern:               models.Tractor,
					NumCards:              6,
					TractorNumConsecutive: 2,
					Suit:                  models.Diamond,
					IsTrump:               false,
					LargestCard: models.Card{
						Value: 10,
						Suit:  models.Diamond,
					},
				}, {
					Pattern:               models.Tractor,
					NumCards:              6,
					TractorNumConsecutive: 3,
					Suit:                  models.Diamond,
					IsTrump:               false,
					LargestCard: models.Card{
						Value: 6,
						Suit:  models.Diamond,
					},
				}, {
					Pattern:  models.NOfAKind,
					NumCards: 1,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value: 1,
						Suit:  models.Diamond,
					},
				},
			},
			next: []models.Trick{
				{
					Pattern:               models.Tractor,
					NumCards:              6,
					TractorNumConsecutive: 2,
					Suit:                  models.Club,
					IsTrump:               true,
					LargestCard: models.Card{
						Value: 7,
						Suit:  models.Club,
					},
				}, {
					Pattern:               models.Tractor,
					NumCards:              6,
					TractorNumConsecutive: 3,
					Suit:                  models.Club,
					IsTrump:               true,
					LargestCard: models.Card{
						Value: 10,
						Suit:  models.Club,
					},
				}, {
					Pattern:  models.NOfAKind,
					NumCards: 1,
					Suit:     models.Diamond,
					IsTrump:  true,
					LargestCard: models.Card{
						Value: 2,
						Suit:  models.Spade,
					},
				},
			},
			trumpNumber: 2,
			trumpSuit:   models.Club,
			expected:    true,
		},
	}
	for _, tc := range tests {
		game := models.Game{
			TrumpNumber: tc.trumpNumber,
			TrumpSuit:   tc.trumpSuit,
		}

		vals := game.GetCardValues()
		for i, c := range tc.prev {
			c.LargestCard = c.LargestCard.WithGameValues(vals)
			c.LargestCard = c.LargestCard.WithTrump(tc.trumpNumber, tc.trumpSuit)
			tc.prev[i] = c
		}

		for i, c := range tc.next {
			c.LargestCard = c.LargestCard.WithGameValues(vals)
			c.LargestCard = c.LargestCard.WithTrump(tc.trumpNumber, tc.trumpSuit)
			tc.next[i] = c
		}

		wins := models.NextTrickWins(tc.prev, tc.next)
		if wins != tc.expected {
			t.Errorf("expected %v but got %v for tricks %v, %v", tc.expected, wins, tc.prev, tc.next)
		}
	}
}
func TestParse(t *testing.T) {
	type test struct {
		input       []models.Card
		expected    models.Trick
		err         error
		trumpSuit   models.Suit
		trumpNumber int
	}
	tests := []test{
		{
			input: []models.Card{
				{
					Value: 5,
					Suit:  models.Spade,
				},
			},
			expected: models.Trick{
				Pattern: models.NOfAKind,
				LargestCard: models.Card{
					Value: 5,
					Suit:  models.Spade,
				},
				NumCards: 1,
				Suit:     models.Spade,
				IsTrump:  false,
			},
			err:         nil,
			trumpSuit:   models.Diamond,
			trumpNumber: 2,
		}, {
			trumpSuit:   models.Diamond,
			trumpNumber: 2,
			input: []models.Card{
				{
					Value: 8,
					Suit:  models.Diamond,
				},
				{
					Value: 8,
					Suit:  models.Diamond,
				},
			},
			expected: models.Trick{
				Pattern: models.NOfAKind,
				LargestCard: models.Card{
					Value: 8,
					Suit:  models.Diamond,
				},
				NumCards: 2,
				Suit:     models.Diamond,
				IsTrump:  true,
			},
			err: nil,
		}, {
			trumpSuit:   models.Diamond,
			trumpNumber: 2,
			input: []models.Card{
				{
					Value: 8,
					Suit:  models.Diamond,
				},
				{
					Value: 8,
					Suit:  models.Spade,
				},
			},
			expected: models.Trick{
				Pattern: models.NOfAKind,
				LargestCard: models.Card{
					Value: 8,
					Suit:  models.Diamond,
				},
				NumCards: 2,
				Suit:     models.Diamond,
				IsTrump:  true,
			},
			err: fmt.Errorf("only one card is trump"),
		}, {
			trumpSuit:   models.Club,
			trumpNumber: 2,
			input: []models.Card{
				{
					Value: 8,
					Suit:  models.Diamond,
				},
				{
					Value: 8,
					Suit:  models.Spade,
				},
			},
			expected: models.Trick{
				Pattern: models.NOfAKind,
				LargestCard: models.Card{
					Value: 8,
					Suit:  models.Spade,
				},
				NumCards: 2,
				Suit:     models.Diamond,
				IsTrump:  false,
			},
			err: fmt.Errorf("not the same suit"),
		}, {
			trumpSuit:   models.Spade,
			trumpNumber: 5,
			input: []models.Card{
				{
					Value: 1,
					Suit:  models.Spade,
				},
				{
					Value: 1,
					Suit:  models.Spade,
				},
				{
					Value: 5,
					Suit:  models.Diamond,
				},
				{
					Value: 5,
					Suit:  models.Diamond,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value: 5,
					Suit:  models.Diamond,
				},
				NumCards:              4,
				Suit:                  models.Diamond,
				IsTrump:               true,
				TractorNumConsecutive: 2,
			},
			err: nil,
		}, {
			trumpSuit:   models.Club,
			trumpNumber: 10,
			input: []models.Card{
				{
					Value: 13,
					Suit:  models.Diamond,
				},
				{
					Value: 13,
					Suit:  models.Diamond,
				},
				{
					Value: 1,
					Suit:  models.Diamond,
				},
				{
					Value: 1,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Diamond,
				},
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 3,
					Suit:  models.Diamond,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value: 3,
					Suit:  models.Diamond,
				},
				NumCards:              8,
				Suit:                  models.Diamond,
				IsTrump:               false,
				TractorNumConsecutive: 2,
			},
			err: nil,
		}, {
			trumpSuit:   models.Diamond,
			trumpNumber: 2,
			input: []models.Card{
				{
					Value: 1,
					Suit:  models.Joker,
				},
				{
					Value: 1,
					Suit:  models.Joker,
				},
				{
					Value: 1,
					Suit:  models.Joker,
				},
				{
					Value: 2,
					Suit:  models.Joker,
				},
				{
					Value: 2,
					Suit:  models.Joker,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value: 2,
					Suit:  models.Joker,
				},
				NumCards:              5,
				Suit:                  models.Joker,
				IsTrump:               true,
				TractorNumConsecutive: 2,
			},
			err: fmt.Errorf("tractor incorrect length"),
		}, {
			trumpSuit:   models.Club,
			trumpNumber: 3,
			input: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 3,
					Suit:  models.Club,
				},
				{
					Value: 3,
					Suit:  models.Club,
				},
				{
					Value: 3,
					Suit:  models.Club,
				},
				{
					Value: 1,
					Suit:  models.Joker,
				},
				{
					Value: 1,
					Suit:  models.Joker,
				},
				{
					Value: 1,
					Suit:  models.Joker,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value: 1,
					Suit:  models.Joker,
				},
				NumCards:              9,
				Suit:                  models.Joker,
				IsTrump:               true,
				TractorNumConsecutive: 3,
			},
			err: nil,
		}, {
			trumpSuit:   models.Club,
			trumpNumber: 5,
			input: []models.Card{
				{
					Value: 2,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Diamond,
				},
				{
					Value: 4,
					Suit:  models.Diamond,
				},
				{
					Value: 4,
					Suit:  models.Diamond,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value: 4,
					Suit:  models.Diamond,
				},
				NumCards:              4,
				Suit:                  models.Diamond,
				IsTrump:               false,
				TractorNumConsecutive: 2,
			},
			err: fmt.Errorf("tractor not consecutive"),
		}, {
			trumpSuit:   models.Club,
			trumpNumber: 3,
			input: []models.Card{
				{
					Value: 2,
					Suit:  models.Diamond,
				},
				{
					Value: 2,
					Suit:  models.Diamond,
				},
				{
					Value: 4,
					Suit:  models.Diamond,
				},
				{
					Value: 4,
					Suit:  models.Diamond,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value: 4,
					Suit:  models.Diamond,
				},
				NumCards:              4,
				Suit:                  models.Diamond,
				IsTrump:               false,
				TractorNumConsecutive: 2,
			},
			err: nil,
		}, {
			trumpSuit:   models.Club,
			trumpNumber: 5,
			input: []models.Card{
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 3,
					Suit:  models.Diamond,
				},
				{
					Value: 4,
					Suit:  models.Diamond,
				},
				{
					Value: 4,
					Suit:  models.Diamond,
				},
				{
					Value: 4,
					Suit:  models.Diamond,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value: 4,
					Suit:  models.Diamond,
				},
				NumCards:              4,
				Suit:                  models.Diamond,
				IsTrump:               false,
				TractorNumConsecutive: 2,
			},
			err: fmt.Errorf("tractor incorrect length"),
		},
	}

	for _, tc := range tests {
		game := models.Game{
			TrumpNumber: tc.trumpNumber,
			TrumpSuit:   tc.trumpSuit,
		}

		vals := game.GetCardValues()
		for i, c := range tc.input {
			tc.input[i] = c.WithGameValues(vals)
			tc.input[i] = tc.input[i].WithTrump(tc.trumpNumber, tc.trumpSuit)
		}

		tc.expected.LargestCard = tc.expected.LargestCard.WithGameValues(vals)
		tc.expected.LargestCard = tc.expected.LargestCard.WithTrump(tc.trumpNumber, tc.trumpSuit)

		tr, err := models.ParseTrick(tc.input)
		if err == nil && tc.err != nil {
			t.Errorf("expected error %v but got no error", tc.err)
		} else if err != nil && tc.err == nil {
			t.Errorf("expected no error but got error %v", err)
		}
		if tc.err == nil && tr != tc.expected {
			t.Errorf("incorrect parsing for %v:\nexpected %v,\ngot %v", tc.input, tc.expected, tr)
		}
	}
}
