package gogen

import (
	"fmt"
	"path"
)

type Gogen struct {
	//UserTypes stores all the user generated user types
	UserTypes map[string]*UserTypeExpr

	pkg string
}

func (g *Gogen) Package() string {
	return g.pkg
}

func (g *Gogen) SetPackage(pkg string) {
	if !path.IsAbs(pkg) {
		panic(fmt.Sprintf("must provide full gopath for package name `%s`", pkg))
	}
	g.pkg = pkg[1:]
}

var Root *Gogen

func init() {
	Root = &Gogen{
		UserTypes: make(map[string]*UserTypeExpr),
	}
}

// UserType looksup a usertype, if none could be found it returns a nil value
func (g *Gogen) UserType(name string) *UserTypeExpr {
	ut, ok := g.UserTypes[name]
	if !ok {
		return nil
	}
	return ut
}

func (g *Gogen) AddUserType(ut *UserTypeExpr) {
	g.UserTypes[ut.Name()] = ut
}
