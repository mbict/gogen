package gogen

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
	"bytes"
	"errors"
)

type CodeGenerator interface {
	//Returns the basic template engine
	Template() *template.Template
}

type codeGenerator struct {
	template *template.Template
}

func NewCodeGenerator(templatePath ...string) CodeGenerator {

	t := template.New("base")

	t.Funcs(template.FuncMap{
		"snake":   snakeCase,
		"title":   strings.Title,
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
		"untitle": untitle,
		//go specific functions
		"packageName": PackageName,
		"dict":        dict,
	})

	//base templates
	_, file, _, _ := runtime.Caller(0)
	pattern := filepath.Join(path.Dir(file), "templates", "*.tmpl")
	template.Must(t.ParseGlob(pattern))

	//user provided templates
	for _, path := range templatePath {
		template.Must(t.ParseGlob(path))
	}

	return &codeGenerator{
		template: t,
	}
}

func (c *codeGenerator) Template() *template.Template {
	return c.template
}

//Make the first charater lower case
func untitle(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

// SnakeCase produces the snake_case of a CamelCase stirng
func snakeCase(name string) string {
	var b bytes.Buffer
	var lastUnderscore bool
	ln := len(name)
	if ln == 0 {
		return ""
	}
	b.WriteRune(unicode.ToLower(rune(name[0])))
	for i := 1; i < ln; i++ {
		r := rune(name[i])
		nextIsLower := false
		if i < ln-1 {
			n := rune(name[i + 1])
			nextIsLower = unicode.IsLower(n) && unicode.IsLetter(n)
		}
		if unicode.IsUpper(r) {
			if !lastUnderscore && nextIsLower {
				b.WriteRune('_')
				lastUnderscore = true
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
			lastUnderscore = false
		}
	}
	return b.String()
}

// dict will provide the template engine a way to create new dictionaries key/value
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call, must be pairs (string/interface)")
	}
	dict := make(map[string]interface{}, len(values) / 2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i + 1]
	}
	return dict, nil
}
