package team

import (
	"github.com/miruken-go/demo.microservice/team-api/data"
	"github.com/miruken-go/demo.microservice/team-api/queries"
	"github.com/miruken-go/miruken/handles"
)

//go:generate $GOPATH/bin/miruken -tests

type (
	Handler struct {
		nextId int32
	}
)

func (h *Handler) Find(
	_ *handles.It, find queries.FindTeams,
) ([]data.Team, error) {
	return []data.Team{
		{Id: 1, Name: "Breakaway", Color: data.ColorOrange},
	}, nil
}
