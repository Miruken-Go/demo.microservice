package queries

import "github.com/miruken-go/demo.microservice/team-api/data"

type (
	FindPeople struct {
		Filter data.Person
	}

	FindTeams struct {
		Filter data.Team
	}
)
