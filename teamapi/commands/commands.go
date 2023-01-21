package commands

import "github.com/miruken-go/demo.microservice/teamapi/data"

type (
	CreatePerson struct {
		Person data.Person
	}

	UpdatePerson struct {
		Person data.Person
	}

	DeletePeople struct {
		Filter data.Person
	}

	CreateTeam struct {
		Team data.Team
	}

	UpdateTeam struct {
		Team data.Team
	}
)
