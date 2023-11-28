package api

import "github.com/miruken-go/miruken/creates"

//go:generate $GOPATH/bin/miruken -tests

// Factory creates api from a type id.
type Factory struct{}

func (F *Factory) New(
	_*struct {
		t creates.It `key:"api.Tag"`
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
		tp creates.It `key:"api.TagPrincipal"`
		up creates.It `key:"api.UntagPrincipal"`
	    ae creates.It `key:"api.AssignEntitlements"`
	    re creates.It `key:"api.RevokeEntitlements"`
		dp creates.It `key:"api.RemovePrincipals"`
		gp creates.It `key:"api.GetPrincipal"`
		fp creates.It `key:"api.FindPrincipals"`

		ce creates.It `key:"api.CreateEntitlement"`
		te creates.It `key:"api.TagEntitlement"`
		ue creates.It `key:"api.UntagEntitlement"`
		de creates.It `key:"api.RemoveEntitlements"`
		ge creates.It `key:"api.GetEntitlement"`
		fe creates.It `key:"api.FindEntitlements"`

		ct creates.It `key:"api.CreateTag"`
		rt creates.It `key:"api.RemoveTags"`
		gt creates.It `key:"api.GetTag"`
		ft creates.It `key:"api.FindTags"`
	  }, create *creates.It,
) any {
	switch create.Key() {
	case "api.Tag":
		return new(Tag)
	case "api.Subject":
		return new(Subject)
	case "api.Principal":
		return new(Principal)
	case "api.Entitlement":
		return new(Entitlement)

	case "api.CreateSubject":
		return new(CreateSubject)
	case "api.RemoveSubjects":
		return new(RemoveSubjects)
	case "api.GetSubject":
		return new(GetSubject)
	case "api.FindSubjects":
		return new(FindSubjects)

	case "api.CreatePrincipal":
		return new(CreatePrincipal)
	case "api.TagPrincipal":
		return new(TagPrincipal)
	case "api.UntagPrincipal":
		return new(UntagPrincipal)
	case "api.AssignPrincipals":
		return new(AssignPrincipals)
	case "api.RevokePrincipals":
		return new(RevokePrincipals)
	case "api.RemovePrincipals":
		return new(RemovePrincipals)
	case "api.GetPrincipal":
		return new(GetPrincipal)
	case "api.FindPrincipals":
		return new(FindPrincipals)

	case "api.CreateEntitlement":
		return new(CreateEntitlement)
	case "api.TagEntitlement":
		return new(TagEntitlement)
	case "api.UntagEntitlement":
		return new(UntagEntitlement)
	case "api.AssignEntitlements":
		return new(AssignEntitlements)
	case "api.RevokeEntitlements":
		return new(RevokeEntitlements)
	case "api.RemoveEntitlements":
		return new(RevokeEntitlements)
	case "api.GetEntitlement":
		return new(GetEntitlement)
	case "api.FindEntitlements":
		return new(FindEntitlements)

	case "api.CreateTag":
		return new(CreateTag)
	case "api.RemoveTags":
		return new(RemoveTags)
	case "api.GetTag":
		return new(GetTag)
	case "api.FindTags":
		return new(FindTags)
	}

	return nil
}
