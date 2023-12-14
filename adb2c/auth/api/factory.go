package api

import "github.com/miruken-go/miruken/creates"

//go:generate $GOPATH/bin/miruken -tests

// Factory creates api from a type id.
type Factory struct{}

func (F *Factory) New(
	_ *struct {
		s creates.It `key:"api.Subject"`
		p creates.It `key:"api.Principal"`
		e creates.It `key:"api.Entitlement"`

		cs creates.It `key:"api.CreateSubject"`
		ap creates.It `key:"api.AssignPrincipals"`
		rp creates.It `key:"api.RevokePrincipals"`
		rs creates.It `key:"api.RemoveSubjects"`
		gs creates.It `key:"api.GetSubject"`
		fs creates.It `key:"api.FindSubjects"`

		cp creates.It `key:"api.CreatePrincipal"`
		ae creates.It `key:"api.AssignEntitlements"`
		re creates.It `key:"api.RevokeEntitlements"`
		dp creates.It `key:"api.RemovePrincipal"`
		gp creates.It `key:"api.GetPrincipal"`
		fp creates.It `key:"api.FindPrincipals"`

		ce creates.It `key:"api.CreateEntitlement"`
		de creates.It `key:"api.RemoveEntitlement"`
		ge creates.It `key:"api.GetEntitlement"`
		fe creates.It `key:"api.FindEntitlements"`
	}, create *creates.It,
) any {
	switch create.Key() {
	case "api.Subject":
		return new(Subject)
	case "api.Principal":
		return new(Principal)
	case "api.Entitlement":
		return new(Entitlement)

	case "api.CreateSubject":
		return new(CreateSubject)
	case "api.AssignPrincipals":
		return new(AssignPrincipals)
	case "api.RevokePrincipals":
		return new(RevokePrincipals)
	case "api.RemoveSubject":
		return new(RemoveSubject)
	case "api.GetSubject":
		return new(GetSubject)
	case "api.FindSubjects":
		return new(FindSubjects)

	case "api.CreatePrincipal":
		return new(CreatePrincipal)
	case "api.AssignEntitlements":
		return new(AssignEntitlements)
	case "api.RevokeEntitlements":
		return new(RevokeEntitlements)
	case "api.RemovePrincipal":
		return new(RemovePrincipal)
	case "api.GetPrincipal":
		return new(GetPrincipal)
	case "api.FindPrincipals":
		return new(FindPrincipals)

	case "api.CreateEntitlement":
		return new(CreateEntitlement)
	case "api.RemoveEntitlement":
		return new(RemoveEntitlement)
	case "api.GetEntitlement":
		return new(GetEntitlement)
	case "api.FindEntitlements":
		return new(FindEntitlements)
	}

	return nil
}
