package generator

import (
	"bytes"
	"io"
)

// A FileWriter exposes a set of Sections and the relative path to the output file.
type FileWriter interface {
	//Sections in this file
	Sections() []Section

	//Path returns the realative path to the file to be written
	Path() string

	//WriteString runs the template and returns the generated string
	Write(io io.Writer) error

	//WriteString runs the template and returns the generated string
	WriteString() (string, error)
}

func NewFileWriter(sections []Section, path string) FileWriter {
	return &fileWriter{
		sections: sections,
		path:     path,
	}
}

type fileWriter struct {
	sections []Section
	path     string
}

//Sections in this file
func (fw *fileWriter) Sections() []Section {
	return fw.sections
}

//Path returns the realative path to the file to be written
func (fw *fileWriter) Path() string {
	return fw.path
}

func (fw *fileWriter) Write(w io.Writer) error {
	for _, s := range fw.sections {
		err := s.Generate(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fw *fileWriter) WriteString() (string, error) {
	w := bytes.NewBuffer(nil)
	for _, s := range fw.sections {
		err := s.Generate(w)
		if err != nil {
			return "", err
		}
	}
	return w.String(), nil
}
