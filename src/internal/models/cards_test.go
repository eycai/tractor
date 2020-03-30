package models_test

import (
	"fmt"
	"testing"

	"github.com/eycai/tractor/src/internal/models"
)

func TestNextTrickWins(t *testing.T) {
	type test struct {
		prev     []models.Trick
		next     []models.Trick
		expected bool
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
						Value:     5,
						GameValue: 3,
						Suit:      models.Diamond,
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
						Value:     8,
						GameValue: 6,
						Suit:      models.Diamond,
					},
				},
			},
			expected: true,
		}, {
			prev: []models.Trick{
				{
					Pattern:  models.NOfAKind,
					NumCards: 1,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value:     5,
						GameValue: 3,
						Suit:      models.Diamond,
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
						Value:     8,
						GameValue: 6,
						Suit:      models.Spade,
					},
				},
			},
			expected: false,
		}, {
			prev: []models.Trick{
				{
					Pattern:  models.NOfAKind,
					NumCards: 2,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value:     5,
						GameValue: 3,
						Suit:      models.Diamond,
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
						Value:       2,
						GameValue:   36,
						Suit:        models.Spade,
						IsTrumpSuit: true,
					},
				},
			},
			expected: true,
		}, {
			prev: []models.Trick{
				{
					Pattern:               models.Tractor,
					NumCards:              4,
					TractorNumConsecutive: 2,
					Suit:                  models.Diamond,
					IsTrump:               false,
					LargestCard: models.Card{
						Value:     5,
						GameValue: 3,
						Suit:      models.Diamond,
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
						Value:       2,
						GameValue:   36,
						Suit:        models.Spade,
						IsTrumpSuit: true,
					},
				},
			},
			expected: false,
		}, {
			prev: []models.Trick{
				{
					Pattern:               models.Tractor,
					NumCards:              4,
					TractorNumConsecutive: 2,
					Suit:                  models.Diamond,
					IsTrump:               false,
					LargestCard: models.Card{
						Value:     5,
						GameValue: 3,
						Suit:      models.Diamond,
					},
				}, {
					Pattern:  models.NOfAKind,
					NumCards: 2,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value:     10,
						GameValue: 8,
						Suit:      models.Diamond,
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
						Value:     7,
						GameValue: 5,
						Suit:      models.Diamond,
					},
				}, {
					Pattern:  models.NOfAKind,
					NumCards: 2,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value:     2,
						GameValue: 0,
						Suit:      models.Diamond,
					},
				},
			},
			expected: false,
		}, {
			prev: []models.Trick{
				{
					Pattern:               models.Tractor,
					NumCards:              6,
					TractorNumConsecutive: 2,
					Suit:                  models.Diamond,
					IsTrump:               false,
					LargestCard: models.Card{
						Value:     10,
						GameValue: 7,
						Suit:      models.Diamond,
					},
				}, {
					Pattern:               models.Tractor,
					NumCards:              6,
					TractorNumConsecutive: 3,
					Suit:                  models.Diamond,
					IsTrump:               false,
					LargestCard: models.Card{
						Value:     6,
						GameValue: 3,
						Suit:      models.Diamond,
					},
				}, {
					Pattern:  models.NOfAKind,
					NumCards: 1,
					Suit:     models.Diamond,
					IsTrump:  false,
					LargestCard: models.Card{
						Value:     1,
						GameValue: 11,
						Suit:      models.Diamond,
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
						Value:       7,
						GameValue:   40,
						Suit:        models.Club,
						IsTrumpSuit: true,
					},
				}, {
					Pattern:               models.Tractor,
					NumCards:              6,
					TractorNumConsecutive: 3,
					Suit:                  models.Club,
					IsTrump:               true,
					LargestCard: models.Card{
						Value:       10,
						GameValue:   43,
						Suit:        models.Club,
						IsTrumpSuit: true,
					},
				}, {
					Pattern:  models.NOfAKind,
					NumCards: 1,
					Suit:     models.Diamond,
					IsTrump:  true,
					LargestCard: models.Card{
						Value:         2,
						GameValue:     48,
						Suit:          models.Spade,
						IsTrumpNumber: true,
					},
				},
			},
			expected: true,
		},
	}
	for _, tc := range tests {
		wins := models.NextTrickWins(tc.prev, tc.next)
		if wins != tc.expected {
			t.Errorf("expected %v but got %v for tricks %v, %v", tc.expected, wins, tc.prev, tc.next)
		}
	}
}
func TestParse(t *testing.T) {
	type test struct {
		input    []models.Card
		expected models.Trick
		err      error
	}
	tests := []test{
		{
			input: []models.Card{
				models.Card{
					Value:     5,
					Suit:      models.Spade,
					GameValue: 2,
				},
			},
			expected: models.Trick{
				Pattern: models.NOfAKind,
				LargestCard: models.Card{
					Value:     5,
					Suit:      models.Spade,
					GameValue: 2,
				},
				NumCards: 1,
				Suit:     models.Spade,
				IsTrump:  false,
			},
			err: nil,
		}, {
			input: []models.Card{
				models.Card{
					Value:       8,
					Suit:        models.Diamond,
					IsTrumpSuit: true,
					GameValue:   41,
				},
				models.Card{
					Value:       8,
					Suit:        models.Diamond,
					IsTrumpSuit: true,
					GameValue:   41,
				},
			},
			expected: models.Trick{
				Pattern: models.NOfAKind,
				LargestCard: models.Card{
					Value:       8,
					Suit:        models.Diamond,
					IsTrumpSuit: true,
					GameValue:   41,
				},
				NumCards: 2,
				Suit:     models.Diamond,
				IsTrump:  true,
			},
			err: nil,
		}, {
			input: []models.Card{
				models.Card{
					Value:       8,
					Suit:        models.Diamond,
					IsTrumpSuit: true,
					GameValue:   41,
				},
				models.Card{
					Value:       8,
					Suit:        models.Spade,
					IsTrumpSuit: false,
					GameValue:   29,
				},
			},
			expected: models.Trick{
				Pattern: models.NOfAKind,
				LargestCard: models.Card{
					Value:       8,
					Suit:        models.Diamond,
					IsTrumpSuit: true,
					GameValue:   41,
				},
				NumCards: 2,
				Suit:     models.Diamond,
				IsTrump:  true,
			},
			err: fmt.Errorf("only one card is trump"),
		}, {
			input: []models.Card{
				models.Card{
					Value:       8,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   17,
				},
				models.Card{
					Value:       8,
					Suit:        models.Spade,
					IsTrumpSuit: false,
					GameValue:   29,
				},
			},
			expected: models.Trick{
				Pattern: models.NOfAKind,
				LargestCard: models.Card{
					Value:       8,
					Suit:        models.Spade,
					IsTrumpSuit: false,
					GameValue:   29,
				},
				NumCards: 2,
				Suit:     models.Diamond,
				IsTrump:  false,
			},
			err: fmt.Errorf("not the same suit"),
		}, {
			input: []models.Card{
				models.Card{
					Value:       1,
					Suit:        models.Spade,
					IsTrumpSuit: true,
					GameValue:   47,
				},
				models.Card{
					Value:       1,
					Suit:        models.Spade,
					IsTrumpSuit: true,
					GameValue:   47,
				},
				models.Card{
					Value:         5,
					Suit:          models.Diamond,
					IsTrumpSuit:   false,
					IsTrumpNumber: true,
					GameValue:     48,
				},
				models.Card{
					Value:         5,
					Suit:          models.Diamond,
					IsTrumpSuit:   false,
					IsTrumpNumber: true,
					GameValue:     48,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value:         5,
					Suit:          models.Diamond,
					IsTrumpSuit:   false,
					IsTrumpNumber: true,
					GameValue:     48,
				},
				NumCards:              4,
				Suit:                  models.Diamond,
				IsTrump:               true,
				TractorNumConsecutive: 2,
			},
			err: nil,
		}, {
			input: []models.Card{
				models.Card{
					Value:       13,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   11,
				},
				models.Card{
					Value:       13,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   11,
				},
				models.Card{
					Value:       1,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   12,
				},
				models.Card{
					Value:       1,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   12,
				},
				models.Card{
					Value:       2,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   0,
				},
				models.Card{
					Value:       2,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   0,
				},
				models.Card{
					Value:       3,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   1,
				},
				models.Card{
					Value:       3,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   1,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value:       3,
					Suit:        models.Diamond,
					IsTrumpSuit: false,
					GameValue:   1,
				},
				NumCards:              8,
				Suit:                  models.Diamond,
				IsTrump:               false,
				TractorNumConsecutive: 2,
			},
			err: nil,
		}, {
			input: []models.Card{
				models.Card{
					Value:       1,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   50,
				},
				models.Card{
					Value:       1,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   50,
				},
				models.Card{
					Value:       1,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   50,
				},
				models.Card{
					Value:       2,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   51,
				},
				models.Card{
					Value:       2,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   51,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value:       2,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   51,
				},
				NumCards:              5,
				Suit:                  models.Joker,
				IsTrump:               true,
				TractorNumConsecutive: 2,
			},
			err: fmt.Errorf("tractor incorrect length"),
		}, {
			input: []models.Card{
				models.Card{
					Value:         3,
					Suit:          models.Diamond,
					IsTrumpSuit:   false,
					IsTrumpNumber: true,
					GameValue:     48,
				},
				models.Card{
					Value:         3,
					Suit:          models.Diamond,
					IsTrumpSuit:   false,
					IsTrumpNumber: true,
					GameValue:     48,
				},
				models.Card{
					Value:         3,
					Suit:          models.Diamond,
					IsTrumpSuit:   false,
					IsTrumpNumber: true,
					GameValue:     48,
				},
				models.Card{
					Value:         3,
					Suit:          models.Club,
					IsTrumpSuit:   true,
					IsTrumpNumber: true,
					GameValue:     49,
				},
				models.Card{
					Value:         3,
					Suit:          models.Club,
					IsTrumpSuit:   true,
					IsTrumpNumber: true,
					GameValue:     49,
				},
				models.Card{
					Value:         3,
					Suit:          models.Club,
					IsTrumpSuit:   true,
					IsTrumpNumber: true,
					GameValue:     49,
				},
				models.Card{
					Value:       1,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   50,
				},
				models.Card{
					Value:       1,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   50,
				},
				models.Card{
					Value:       1,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   50,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value:       1,
					Suit:        models.Joker,
					IsTrumpSuit: true,
					GameValue:   50,
				},
				NumCards:              9,
				Suit:                  models.Joker,
				IsTrump:               true,
				TractorNumConsecutive: 3,
			},
			err: nil,
		}, {
			input: []models.Card{
				models.Card{
					Value:     2,
					Suit:      models.Diamond,
					GameValue: 0,
				},
				models.Card{
					Value:     2,
					Suit:      models.Diamond,
					GameValue: 0,
				},
				models.Card{
					Value:     4,
					Suit:      models.Diamond,
					GameValue: 2,
				},
				models.Card{
					Value:     4,
					Suit:      models.Diamond,
					GameValue: 2,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value:     4,
					Suit:      models.Diamond,
					GameValue: 2,
				},
				NumCards:              4,
				Suit:                  models.Diamond,
				IsTrump:               false,
				TractorNumConsecutive: 2,
			},
			err: fmt.Errorf("tractor not consecutive"),
		},
		{
			input: []models.Card{
				models.Card{
					Value:     3,
					Suit:      models.Diamond,
					GameValue: 1,
				},
				models.Card{
					Value:     3,
					Suit:      models.Diamond,
					GameValue: 1,
				},
				models.Card{
					Value:     4,
					Suit:      models.Diamond,
					GameValue: 2,
				},
				models.Card{
					Value:     4,
					Suit:      models.Diamond,
					GameValue: 2,
				},
				models.Card{
					Value:     4,
					Suit:      models.Diamond,
					GameValue: 2,
				},
			},
			expected: models.Trick{
				Pattern: models.Tractor,
				LargestCard: models.Card{
					Value:     4,
					Suit:      models.Diamond,
					GameValue: 2,
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
