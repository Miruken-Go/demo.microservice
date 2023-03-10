package commands

import (
	"github.com/miruken-go/demo.microservice/teamapi/data"
	"time"
)

type (
	CreatePerson struct {
		FirstName string
		LastName  string
		BirthDate time.Time
	}

	UpdatePerson struct {
		Id        int32
		FirstName string
		LastName  string
		BirthDate time.Time
	}

	DeletePeople struct {
		Filter data.Person
	}

	CreateTeam struct {
		Name    string
		Color   data.Color
		Coach   data.Coach
		Manager data.Manager
		Players []data.Player
	}

	UpdateTeam struct {
		Id      int32
		Name    string
		Color   data.Color
		Coach   *data.Coach
		Manager *data.Manager
		Players []data.Player
	}
)
