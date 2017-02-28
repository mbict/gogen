package gogen

import "github.com/mbict/gogen/generator"

type Gogen struct {
	//UserTypes stores all the user generated user types
	UserTypes map[string]*UserTypeExpr
}

var Root *Gogen

func (g *Gogen) Name() string {
	return "gogen"
}

func (g *Gogen) Generate(path string) ([]generator.FileWriter, error) {
	//todo: implement the user types / models generator here
	return nil, nil
}

func init() {
	Root = &Gogen{
		UserTypes: make(map[string]*UserTypeExpr),
	}

	generator.Register(Root)
}

// UserType looksup a usertype, if none could be found it returns a nil value
func (g *Gogen) UserType(name string) *UserTypeExpr {
	ut, ok := g.UserTypes[name]
	if !ok {
		return nil
	}
	return ut
}
