package models_test

import (
	"testing"

	"github.com/eycai/tractor/src/internal/models"
)

func testGameValues(d models.Deck) bool {
	values := make([]int, 52)
	for _, c := range d.Cards {
		values[c.GameValue]++
	}
	for i, v := range values {
		if i == 48 {
			if v != 6 {
				return false
			}
		} else if v != 2 {
			return false
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
		TrumpNumber: 2,
		TrumpSuit:   models.Diamond,
	}
	d := g.GetDeck()
	if len(d.Cards) != 108 {
		t.Errorf("expected deck of length 108, got %d", len(d.Cards))
	}
	if !testGameValues(d) {
		t.Errorf("incorrect value assignments")
	}

	g.TrumpNumber = 5
	g.TrumpSuit = models.Club

	d = g.GetDeck()
	if len(d.Cards) != 108 {
		t.Errorf("expected deck of length 108, got %d", len(d.Cards))
	}
	if !testGameValues(d) {
		t.Errorf("incorrect value assignments")
	}

	g.TrumpNumber = 1
	g.TrumpSuit = models.Heart

	d = g.GetDeck()
	if len(d.Cards) != 108 {
		t.Errorf("expected deck of length 108, got %d", len(d.Cards))
	}
	if !testGameValues(d) {
		t.Errorf("incorrect value assignments")
	}
}
