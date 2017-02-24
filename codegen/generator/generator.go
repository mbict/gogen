package generator

import (
	"fmt"
	"strings"
	"github.com/mbict/gogen"
	"path"
)

type Generator interface {
	Name() string
	Generate(string) ([]gogen.FileWriter, error)
}

var generators []Generator

func Register(g Generator) error {
	for _, o := range generators {
		if g.Name() == o.Name() {
			return fmt.Errorf("generator: duplicate generator %s", g.Name())
		}
	}
	generators = append(generators, g)
	return nil
}

//Run will execute all generators in the base directory
func Run(targetPath string, names ...string) error {
	files := []gogen.FileWriter{}
	for _, generator := range getGenerators(names) {
		fw, err := generator.Generate(targetPath)
		if err != nil {
			return fmt.Errorf("Running generator %s resulted in a error: `%s`", generator.Name(), err.Error())
		}
		files = append(files, fw...)
	}

	//todo: run plugins etc
	for _, fw := range files {
		fmt.Println(path.Join(targetPath, fw.Path()))
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
