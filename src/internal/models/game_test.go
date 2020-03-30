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
