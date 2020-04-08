package models_test

import (
	"testing"

	"github.com/eycai/tractor/src/internal/models"
)

func TestDrawOrder(t *testing.T) {
	type test struct {
		room          *models.Room
		expectedOrder []string
	}

	tests := []test{
		{
			room: &models.Room{
				Users: []*models.UserStatus{
					{Username: "a"}, {Username: "b"}, {Username: "c"}, {Username: "d"},
				},
				Game: &models.Game{
					Turn: "a",
				},
			},
			expectedOrder: []string{"a", "b", "c", "d"},
		},
		{
			room: &models.Room{
				Users: []*models.UserStatus{
					{Username: "a"}, {Username: "b"}, {Username: "c"}, {Username: "d"},
				},
				Game: &models.Game{
					Turn: "b",
				},
			},
			expectedOrder: []string{"b", "c", "d", "a"},
		},
		{
			room: &models.Room{
				Users: []*models.UserStatus{
					{Username: "a"}, {Username: "b"}, {Username: "c"}, {Username: "d"},
				},
				Game: &models.Game{
					Turn: "d",
				},
			},
			expectedOrder: []string{"d", "a", "b", "c"},
		},
	}

	for _, tc := range tests {
		order := tc.room.DrawOrder()
		fail := false
		for i, u := range order {
			if tc.expectedOrder[i] != u {
				fail = true
				break
			}
		}
		if fail {
			t.Errorf("expected draw order %v, got %v", tc.expectedOrder, order)
		}
	}
}
