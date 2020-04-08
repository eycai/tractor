package models_test

import (
	"testing"

	"github.com/eycai/tractor/src/internal/models"
)

func TestPlayCard(t *testing.T) {
	type test struct {
		play         [][]models.Card
		user         *models.User
		expectedHand []models.Card
	}

	tests := []test{
		{
			user: &models.User{
				Hand: []models.Card{
					{Value: 2, Suit: models.Diamond},
					{Value: 2, Suit: models.Diamond},
					{Value: 2, Suit: models.Heart},
				},
			},
			play: [][]models.Card{
				{
					{Value: 2, Suit: models.Diamond},
				},
			},
			expectedHand: []models.Card{
				{Value: 2, Suit: models.Diamond},
				{Value: 2, Suit: models.Heart},
			},
		},
		{
			user: &models.User{
				Hand: []models.Card{
					{Value: 2, Suit: models.Diamond},
					{Value: 2, Suit: models.Diamond},
					{Value: 2, Suit: models.Heart},
					{Value: 5, Suit: models.Diamond},
					{Value: 5, Suit: models.Diamond},
				},
			},
			play: [][]models.Card{
				{
					{Value: 2, Suit: models.Diamond},
				},
				{
					{Value: 5, Suit: models.Diamond},
					{Value: 5, Suit: models.Diamond},
				},
			},
			expectedHand: []models.Card{
				{Value: 2, Suit: models.Diamond},
				{Value: 2, Suit: models.Heart},
			},
		},
	}

	for _, tc := range tests {
		tc.user.PlayCards(tc.play)
		if !cardListsEqual(tc.user.Hand, tc.expectedHand) {
			t.Errorf("incorrect hand for play %v: expected hand %v, got %v", tc.play, tc.expectedHand, tc.user.Hand)
		}
	}
}
