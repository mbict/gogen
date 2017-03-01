package gogen

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type GenerateSuite struct{}

var _ = Suite(&GenerateSuite{})

func (s *GenerateSuite) TestHelloWorld(c *C) {

	gen := NewGenerator()

	tests := []struct {
		name     string
		action   *Action
		expected string
	}{
		{
			name: "context with one param",
			action: &Action{
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
				Resource: &Resource{
					Name: "account",
					Path: "/account",
				},
			},
			expected: "" +
				"\n" +
				"type AccountShowContext struct {\n" +
				"\tcontext.Context\n" +
				"\trw http.ResponseWriter\n" +
				"\treq *http.Request\n" +
				"\n" +
				"\taccount_id int32\n" +
				"\n" +
				"\tRequest *AccountShowRequest\n" +
				"}\n",
		},
		{
			name: "context with 2 params",
			action: &Action{
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
						&Field{
							Name: "foo",
							Attribute: &AttributeExpr{
								Type: String,
							},
						},
						&Field{
							Name: "bar",
							Attribute: &AttributeExpr{
								Type: Boolean,
							},
						},
						&Field{
							Name: "baz",
							Attribute: &AttributeExpr{
								Type: Any,
							},
						},
						&Field{
							Name: "arraytest",
							Attribute: &AttributeExpr{
								Type: &Array{
									ElemType: &AttributeExpr{
										Type: Float32,
									},
								},
							},
						},
						&Field{
							Name: "objtest",
							Attribute: &AttributeExpr{
								Type: &UserTypeExpr{
									AttributeExpr: &AttributeExpr{
										Type: String,
									},
									TypeName: "testtype",
								},
							},
						},
					},
				},
				Resource: &Resource{
					Name: "account",
					Path: "/account",
				},
			},
			expected: "" +
				"\n" +
				"type AccountShowContext struct {\n" +
				"\tcontext.Context\n" +
				"\trw http.ResponseWriter\n" +
				"\treq *http.Request\n" +
				"\n" +
				"\taccount_id int32\n" +
				"\tfoo string\n" +
				"\tbar boolean\n" +
				"\tbaz interface{}\n" +
				"\tarraytest []float32\n" +
				"\tobjtest testtype\n" +
				"\n" +
				"\tRequest *AccountShowRequest\n" +
				"}\n",
		},
		{
			name: "context no param",
			action: &Action{
				Endpoint: &Endpoint{
					Name: "show",
				},
				Routes: []*Route{
					&Route{
						Method: "GET",
						Path:   "/{account_id}",
					},
				},
				Params: nil,
				Resource: &Resource{
					Name: "account",
					Path: "/account",
				},
			},
			expected: "" +
				"\n" +
				"type AccountShowContext struct {\n" +
				"\tcontext.Context\n" +
				"\trw http.ResponseWriter\n" +
				"\treq *http.Request\n" +
				"\n" +
				"\tRequest *AccountShowRequest\n" +
				"}\n",
		},
	}

	for _, test := range tests {
		sec, err := gen.GenerateContext(test.action)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		res, err := sec[0].Generate()

		c.Check(err, IsNil)
		c.Check(res, Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}

}
