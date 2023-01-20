package data

import "time"

type (
	Color int

	Person struct {
		Id        int32
		FirstName string
		LastName  string
		BirthDate time.Time
	}

	Player struct {
		Person Person
		Number int
	}

	Coach struct {
		Person  Person
		License string
	}

	Manager struct {
		Person Person
	}

	Team struct {
		Id      int32
		Name    string
		Color   Color
		Coach   Coach
		Manager Manager
		Players []Player
	}
)

const (
	ColorBlack Color = iota
	ColorBlue
	ColorGreen
	ColorLightBlue
	ColorMaroon
	ColorOrange
	ColorRed
	ColorWhite
	ColorYellow
	ColorNone = -1
)
