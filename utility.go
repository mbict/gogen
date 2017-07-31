package gogen

//DescriberExpr is the the interface that allows to set a helping description to a expression
//The dsl function DescriptionExpr(string) will use this to set the the description of the expression
type Describer interface {
	Description() string
	SetDescription(string)
}

type Namespace interface {
	Namespace() string
	SetNamespace(string)
}


// DescriptionExpr is the implementation of the describer interface used for composition
type DescriptionExpr struct {
	description string
}

func (d *DescriptionExpr) Description() string {
	return d.description
}

func (d *DescriptionExpr) SetDescription(description string) {
	d.description = description
}
