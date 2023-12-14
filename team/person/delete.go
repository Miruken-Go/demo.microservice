package person

import (
	"github.com/miruken-go/demo.microservice/team-api/commands"
	"github.com/miruken-go/miruken/handles"
)

func (h *Handler) Delete(
	_ *handles.It, delete *commands.DeletePeople,
) error {
	return nil
}
