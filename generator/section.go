package generator

import (
	"io"
	"text/template"
)

// A Section consists of a template and accompaying render data.
type Section struct {
	// Template used to render section text.
	Template *template.Template

	// Data used as input of template.
	Data interface{}
}

// Generate executes the file generating proces
func (s *Section) Generate(buf io.Writer) error {
	err := s.Template.Execute(buf, s.Data)
	if err != nil {
		return err
	}
	return nil
}
