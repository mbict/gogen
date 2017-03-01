package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// Writers accepts the API expression and returns the file writers used to generate the output.
type WriterGenerator interface {
	Writers(interface{}) ([]FileWriter, error)
}

type Generator interface {
	Name() string
	Generate(string) ([]FileWriter, error)
}

var generators []Generator

// Register registers a new generator to the list for execution
func Register(g Generator) error {
	for _, o := range generators {
		if g.Name() == o.Name() {
			return fmt.Errorf("generator: duplicate generator %s", g.Name())
		}
	}
	generators = append(generators, g)
	return nil
}

// Run will execute all generators in the base directory
func Run(targetPath string, names ...string) error {
	gopath := os.Getenv("GOPATH")

	files := []FileWriter{}
	for _, generator := range getGenerators(names) {
		log.Printf("Running generator `%s`\n", generator.Name())

		fw, err := generator.Generate(targetPath)
		if err != nil {
			return fmt.Errorf("Running generator %s resulted in a error: `%s`", generator.Name(), err.Error())
		}
		files = append(files, fw...)
	}

	//todo: run plugins etc

	for _, fw := range files {
		filename := path.Join(gopath, "src", targetPath, fw.Path())
		os.MkdirAll(path.Dir(filename), 0755)
		actionTaken := "SKIP"
		if _, err := os.Stat(filename); err == nil {
			//nop
		} else {

			buf := bytes.NewBuffer(nil)
			err = fw.Write(buf)
			if err != nil {
				return fmt.Errorf("Running filewrite %s generation resulted in a error:\n`%s`", filename, err.Error())
			}

			ioutil.WriteFile(filename, buf.Bytes(), 0644)
			actionTaken = "OK"
		}
		log.Printf("%-7s%s\n", "["+actionTaken+"]", filename)
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

// Registered returns the names of all the generators registered
func Registered() []string {
	result := make([]string, len(generators))
	for key, generator := range generators {
		result[key] = generator.Name()
	}
	return result
}
