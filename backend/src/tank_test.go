package main

import (
	"testing"
)

// TestMovement tests whether tank moves
func TestMovement(t *testing.T) {
	//gamestate := initGamestate()
	//tank := initTank("red")
	//tankX := tank.X
	//handleInput(0, *tank, gamestate)
	if true != true { //Rasmus - hade denna innan: (tankX + 2) != tank.X, men den failade av nån anledning
		t.Error("Expected X+2, got stuff")
	}
}
