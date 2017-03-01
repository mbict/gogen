package gogen

type Gogen struct {
	//UserTypes stores all the user generated user types
	UserTypes map[string]*UserTypeExpr
}

var Root *Gogen

func init() {
	Root = &Gogen{}
}

// UserType looksup a usertype, if none could be found it returns a nil value
func (g *Gogen) UserType(name string) *UserTypeExpr {
	ut, ok := g.UserTypes[name]
	if !ok {
		return nil
	}
	return ut
}

