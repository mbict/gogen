package gogen

import (
	"path"
)

//	Describer is the the interface that allows to set a helping description to a expression
//	The dsl function DescriptionExpr(string) will use this to set the the description of the expression
type Describer interface {
	Description() string
	SetDescription(string)
}

// Packager is the interface that is used to set the Package/namespace to a expression
type Packager interface {
	Package() string
	SetPackage(string)
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

// DescriptionExpr is the implementation of the describer interface used for composition
type PackageNameExpr struct {
	Name string
}

func (d *PackageNameExpr) Package() string {
	if !path.IsAbs(d.Name) {
		return path.Join(Root.Package(), d.Name)
	}
	return d.Name[1:]
}

func (d *PackageNameExpr) SetPackage(pkg string) {
	d.Name = pkg
}
