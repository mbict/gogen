package gogen

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

var ErrIncompatibleRestRootExpr = errors.New("incompatible rest root expression of type *design.RootExpr")

// Writers accepts the API expression and returns the file writers used to generate the output.
type WriterGenerator interface {
	Writers(interface{}) ([]FileWriter, error)
}

// A FileWriter exposes a set of Sections and the relative path to the output file.
type FileWriter interface {
	//Sections in this file
	Sections() []Section

	//Path returns the realative path to the file to be written
	Path() string

	//Write runs the template and creates the outcome
	Write( /* io */ ) (string, error)
}

// A Section consists of a template and accompaying render data.
type Section struct {
	// Template used to render section text.
	Template *template.Template

	// Data used as input of template.
	Data interface{}
}

// Generate executes the file generating proces
func (s *Section) Generate() (string, error) {
	w := bytes.NewBuffer(nil)
	err := s.Template.Execute(w, s.Data)
	if err != nil {
		return "", err
	}
	return w.String(), nil
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

func (fw *fileWriter) Write() (string, error) {
	w := bytes.NewBuffer(nil)
	for _, s := range fw.sections {
		res, err := s.Generate()
		if err != nil {
			return "", err
		}
		w.WriteString(res)
	}
	return w.String(), nil
}

var ATTRIBUTE_T = `
{{define "ATTRIBUTE"}}
{{- if eq .Type.Name "boolean" "int32" "int64" "uint32" "uint64" "float32" "float64" "string" -}}
{{.Type.Name}}
{{- else if eq .Type.Name "array" -}}
[]{{template "ATTRIBUTE" .Type.ElemType}}
{{- else if eq .Type.Name "map" -}}
map[{{template "ATTRIBUTE" .Type.KeyType}}]{{template "ATTRIBUTE" .Type.ElemType}}
{{- else if eq .Type.Name "[]byte" -}}
{{.Type.Name}}
{{- else if eq .Type.Name "any" -}}
interface{}
{{- else -}}
{{.Type.Name}}
{{- end}}
{{- end}}`

var USERTYPE_T = `
{{define "USERTYPE"}}
type {{title .TypeName}} {{if eq .Name "object"}}
struct {
{{range $prop := .Attribute}}
	{{ $prop.Name }} {{template "ATTRIBUTE" $prop.Attribute}}
{{end}}
}
{{else}}
{{template "ATTRIBUTE" .Attribute}}
{{end}}
{{end}}
`

var HANDLER_T = `
{{define "HANDLER"}}
func {{title .Resource.Name}}{{title .Name}}(ctx *{{title .Resource.Name}}{{title .Name}}Request) error {
	return nil
}
{{end}}
`

var CONTEXT_T = `
{{define "CONTEXT"}}
type {{title .Resource.Name}}{{title .Name}}Context struct {
	context.Context
	rw http.ResponseWriter
	req *http.Request
{{if .Payload}}
	Payload *{{title .Resource.Name}}Payload
{{end}}{{if .Params -}}
	{{range $prop := .Params.Type }}
	{{ $prop.Name }} {{template "ATTRIBUTE" $prop.Attribute}}
	{{- end}}
{{end}}
	Request *{{title .Resource.Name}}{{title .Name}}Request
}
{{end}}
`

var REQUEST_T = `
{{define "REQUEST"}}
type {{title .ResourceName}}{{title .Action}}Request struct {
	{{range $name, $attr := .Attributes}}
		{{$name}} {{template "ATTRIBUTE" $attr}}
	{{end}}
	StallingIds []uuid.UUID
	Names       []string
	Domains     []string
}

func New{{title .ResourceName}}{{title .Action}}Request(req *http.Request, params httprouter.Params) (*{{title .ResourceName}}{{title .Action}}Request, error) {

	r := &{{title .ResourceName}}{{title .Action}}Request{
		Context: ctx,
	}

	stallingIds := req.URL.Query()["stalling_id"]
	if len(stallingIds) > 0 {

	}

	names := req.URL.Query()["names"]
	if len(names) > 0 {
		findRequest.Names = names
	}

	domains := req.URL.Query()["domains"]
	if len(domains) > 0 {
		findRequest.Domains = domains
	}

	return r, nil
}

func (ctx *{{title .ResourceName}}{{title .Action}}Request) OK(result *{{title .ResourceName}}{{title .Action}}Response) error {
	ctx.rw.Header().Set("Content-Type", "application/stalling.collection.stalling+json")
	ctx.rw.WriteHeader(200)
	return nil //ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}
{{end}}
`

var PAYLOAD_T = `
{{define "PAYLOAD"}}
type {{.ResourceName}}{{title .Action}}Payload struct {
	{{range $name, $attr := .RawPayload}}
		{{$name}} {{template "ATTRIBUTE" $attr}}
	{{end}}
	Name        *string
	Description *string
	Domain      *string
}

func (p *{{.ResourceName}}{{title .Action}}Payload) Validate() error {

	err := error(nil)
	if p.Name == nil {

	}

	if p.Domain == nil {
	}

	return err
}


type {{title .ResourceName}}{{title .Action}}Payload struct {
	{{range $name, $attr := .RawPayload}}
		{{$name}} {{template "ATTRIBUTE" $attr}}
	{{end}}
	Name        string
	Description *string
	Domain      string
}

func (p *{{title .ResourceName}}{{title .Action}}Payload) Validate() error {

	err := error(nil)
	if p.Name == "" {

	}

	if p.Domain == "" {
	}

	return err
}
{{end}}
`

type ApiGen struct {
	template *template.Template
}

func NewGenerator() *ApiGen {
	t := template.New("base")

	t.Funcs(template.FuncMap{
		"title":   strings.Title,
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
	})

	template.Must(t.Parse(ATTRIBUTE_T))
	template.Must(t.Parse(USERTYPE_T))
	template.Must(t.Parse(CONTEXT_T))
	template.Must(t.Parse(HANDLER_T))
	template.Must(t.Parse(REQUEST_T))
	template.Must(t.Parse(PAYLOAD_T))

	return &ApiGen{
		template: t,
	}
}

func (g *ApiGen) Writers(root interface{}) ([]FileWriter, error) {
	rootExpr, ok := root.(*Root)
	if !ok {
		return nil, ErrIncompatibleRestRootExpr
	}

	res := []FileWriter{}
	for _, r := range rootExpr.Resources {
		for _, a := range r.Actions {

			s, err := g.GenerateContext(a)
			if err != nil {
				return nil, err
			}

			fwContext := &fileWriter{
				sections: s,
				path:     fmt.Sprintf("operations/%s/%s_context.go", r.Name, a.Name),
			}

			s, err = g.GenerateHandler(a)
			if err != nil {
				return nil, err
			}

			fwHandler := &fileWriter{
				sections: s,
				path:     fmt.Sprintf("operations/%s/%s.go", r.Name, a.Name),
			}

			if a.Payload != nil {
				s, err = g.GeneratePayload(a)
				if err != nil {
					return nil, err
				}

				fwPayload := &fileWriter{
					sections: s,
					path:     fmt.Sprintf("operations/%s/%s.go", r.Name, a.Name),
				}

				res = append(res, fwPayload)
			}

			res = append(res, fwContext, fwHandler)
		}
	}

	return res, nil
}

func (g *ApiGen) GenerateContext(a *Action) ([]Section, error) {
	t := g.template.Lookup("CONTEXT")
	if t == nil {
		return nil, errors.New("template not found")
	}

	s := Section{
		Template: template.Must(t.Clone()),
		Data:     a,
	}

	return []Section{s}, nil
}

func (g *ApiGen) GenerateHandler(a *Action) ([]Section, error) {

	t := g.template.Lookup("HANDLER")
	if t == nil {
		return nil, errors.New("template not found")
	}

	s := Section{
		Template: template.Must(t.Clone()),
		Data:     a,
	}

	return []Section{s}, nil
}

func (g *ApiGen) GeneratePayload(a *Action) ([]Section, error) {

	t := g.template.Lookup("PAYLOAD")
	if t == nil {
		return nil, errors.New("template not found")
	}

	s := Section{
		Template: template.Must(t.Clone()),
		Data:     a,
	}

	return []Section{s}, nil
}

// ===============
// actual test part
func main() {
	root := &Root{}
	genapi(root)

	codegen := NewGenerator()
	writers, err := codegen.Writers(root)
	if err != nil {
		panic(err)
	}

	for _, w := range writers {
		c, err := w.Write()
		if err != nil {
			fmt.Printf("Error generating file `%s` : %s\n", w.Path(), err)
		}

		fmt.Printf("\n\n[%s]\n%s\n", w.Path(), c)
	}

	fmt.Println(writers)
}

func genapi(r *Root) {
	r.Path = "/api"
	r.Resources = []*Resource{
		&Resource{
			Name:                "account",
			Path:                "/account",
			CanonicalActionName: "show",
			Actions: []*Action{
				&Action{
					Endpoint: &Endpoint{
						Name: "find",
					},
					Routes: []*Route{
						&Route{
							Method: "GET",
							Path:   "/",
						},
					},
				},
				&Action{
					Endpoint: &Endpoint{
						Name: "show",
					},
					Routes: []*Route{
						&Route{
							Method: "GET",
							Path:   "/{account_id}",
						},
					},
					Params: &AttributeExpr{
						Type: &Object{
							&Field{
								Name: "account_id",
								Attribute: &AttributeExpr{
									Type: Int32,
								},
							},
						},
					},
				},
			},
		},
		&Resource{
			Name:       "user",
			ParentName: "account",
			Path:       "/user",
			Actions: []*Action{
				&Action{
					Endpoint: &Endpoint{
						Name: "find",
					},
					Routes: []*Route{
						&Route{
							Method: "GET",
							Path:   "/",
						},
					},
				},
				&Action{
					Endpoint: &Endpoint{
						Name: "show",
					},
					Routes: []*Route{
						&Route{
							Method: "GET",
							Path:   "/{user_id}",
						},
					},
				},
			},
		},
		&Resource{
			Name: "test",
			Path: "/test",
			Actions: []*Action{
				&Action{
					Endpoint: &Endpoint{
						Name: "find",
					},
					Routes: []*Route{
						&Route{
							Method: "GET",
							Path:   "/",
						},
					},
				},
				&Action{
					Endpoint: &Endpoint{
						Name: "show",
					},
					Routes: []*Route{
						&Route{
							Method: "GET",
							Path:   "/{test_id}",
						},
					},
				},
			},
		},
	}

	//finalize
	for _, r := range r.Resources {
		for _, a := range r.Actions {
			a.Resource = r
			for _, ar := range a.Routes {
				ar.Action = a
			}
		}
	}
}
