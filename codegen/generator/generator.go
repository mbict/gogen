package generator

import (
	"fmt"
	"strings"
)

type Generator interface {
	Name() string
	Generate(string) error
}

var generators []Generator

func Register(g Generator) error {
	for _, o := range generators {
		if g.Name() == o.Name() {
			return fmt.Errorf("generator: duplicate generator %s", g.Name())
		}
	}
	//t := reflect.TypeOf(r)
	//if t.Kind() == reflect.Ptr {
	//	t = t.Elem()
	//}
	//dslPackages[t.PkgPath()] = true
	generators = append(generators, g)
	return nil
}

//Run will execute all generators in the base directory
func Run(targetPath string, names ...string) error {
	for _, generator := range getGenerators(names) {
		err := generator.Generate(targetPath)
		if err != nil {
			return fmt.Errorf("Running generator %s resulted in a error: `%s`", generator.Name(), err.Error())
		}
	}

	return nil
}

func getGenerators(names []string) []Generator {
	if len(names) == 0 {
		return generators
	}

	res := []Generator{}
	for _, gen := range generators {
		for _, name := range names {
			if strings.EqualFold(name, gen.Name()) {
				res = append(res, gen)
				break
			}
		}
	}
	return res
}

//Registered returns the names of all the generators registered
func Registered() []string {
	result := make([]string, len(generators))
	for key, generator := range generators {
		result[key] = generator.Name()
	}
	return result
}
