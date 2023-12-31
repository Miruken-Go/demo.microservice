package user

import (
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/miruken-go/demo.microservice/adb2c/api"
)

func ToApi(from models.Userable, to *api.User) {
	if id := from.GetId(); id != nil {
		to.Id = *id
	}
	if firstName := from.GetGivenName(); firstName != nil {
		to.FirstName = *firstName
	}
	if lastName := from.GetSurname(); lastName != nil {
		to.LastName = *lastName
	}
	if displayName := from.GetDisplayName(); displayName != nil {
		to.DisplayName = *displayName
	}
}
