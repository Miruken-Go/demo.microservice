package api

import "github.com/miruken-go/miruken/creates"

//go:generate $GOPATH/bin/miruken -tests

// Factory creates api components from a type id.
type Factory struct{}

func (F *Factory) New(
	_ *struct {
	    u creates.It `key:"api.User"`
		s creates.It `key:"api.Subject"`
		p creates.It `key:"api.Principal"`

		cs creates.It `key:"api.CreateSubject"`
		ap creates.It `key:"api.AssignPrincipals"`
		rp creates.It `key:"api.RevokePrincipals"`
		rs creates.It `key:"api.RemoveSubjects"`
		gs creates.It `key:"api.GetSubject"`
		fs creates.It `key:"api.FindSubjects"`

		cp creates.It `key:"api.CreatePrincipal"`
		ip creates.It `key:"api.IncludePrincipals"`
		ep creates.It `key:"api.ExcludePrincipals"`
		dp creates.It `key:"api.RemovePrincipal"`
		gp creates.It `key:"api.GetPrincipal"`
		fp creates.It `key:"api.FindPrincipals"`
	    xp creates.It `key:"api.ExpandPrincipals"`
	    sp creates.It `key:"api.SatisfyPrincipals"`
	  }, create *creates.It,
) any {
	switch create.Key() {
	case "api.User": return new(User)
	case "api.Subject": return new(Subject)
	case "api.Principal": return new(Principal)

	case "api.CreateSubject": return new(CreateSubject)
	case "api.AssignPrincipals": return new(AssignPrincipals)
	case "api.RevokePrincipals": return new(RevokePrincipals)
	case "api.RemoveSubject": return new(RemoveSubject)
	case "api.GetSubject": return new(GetSubject)
	case "api.FindSubjects": return new(FindSubjects)

	case "api.CreatePrincipal": return new(CreatePrincipal)
	case "api.IncludePrincipals": return new(IncludePrincipals)
	case "api.ExcludePrincipals": return new(ExcludePrincipals)
	case "api.RemovePrincipal": return new(RemovePrincipal)
	case "api.GetPrincipal": return new(GetPrincipal)
	case "api.FindPrincipals": return new(FindPrincipals)
	case "api.ExpandPrincipals": return new(ExpandPrincipals)
	case "api.SatisfyPrincipals": return new(SatisfyPrincipals)
	}
	return nil
}
