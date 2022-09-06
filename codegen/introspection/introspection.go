package introspection

import (
	_ "embed"
)

// Query is the query generated by graphiql to determine type information
//
//go:embed introspection.graphql
var Query string

// Response is the introspection query response
type Response struct {
	Schema *Schema `json:"__schema"`
}

type Schema struct {
	QueryType struct {
		Name string `json:"name"`
	} `json:"queryType"`
	MutationType struct {
		Name string `json:"name"`
	} `json:"mutationType"`
	SubscriptionType struct {
		Name string `json:"name"`
	} `json:"subscriptionType"`

	Types Types `json:"types"`
}

func (s *Schema) Query() *Type {
	return s.Types.Get(s.QueryType.Name)
}

func (s *Schema) Mutation() *Type {
	return s.Types.Get(s.MutationType.Name)
}

func (s *Schema) Subscription() *Type {
	return s.Types.Get(s.SubscriptionType.Name)
}

func (s *Schema) Visit(handlers VisitHandlers) error {
	v := Visitor{
		schema:   s,
		handlers: handlers,
	}
	return v.Run()
}

type TypeKind string

const (
	TypeKindScalar      = TypeKind("SCALAR")
	TypeKindObject      = TypeKind("OBJECT")
	TypeKindInterface   = TypeKind("INTERFACE")
	TypeKindUnion       = TypeKind("UNION")
	TypeKindEnum        = TypeKind("ENUM")
	TypeKindInputObject = TypeKind("INPUT_OBJECT")
	TypeKindList        = TypeKind("LIST")
	TypeKindNonNull     = TypeKind("NON_NULL")
)

type Scalar string

const (
	ScalarInt     = Scalar("Int")
	ScalarFloat   = Scalar("Float")
	ScalarString  = Scalar("String")
	ScalarBoolean = Scalar("Boolean")
)

type Type struct {
	Kind        TypeKind     `json:"kind"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Fields      []Field      `json:"fields,omitempty"`
	InputFields []InputValue `json:"inputFields,omitempty"`
}

type Types []*Type

func (t Types) Get(name string) *Type {
	for _, i := range t {
		if i.Name == name {
			return i
		}
	}
	return nil
}

type Field struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	TypeRef     *TypeRef    `json:"type"`
	Args        InputValues `json:"args"`
}

type TypeRef struct {
	Kind   TypeKind `json:"kind"`
	Name   string   `json:"name,omitempty"`
	OfType *TypeRef `json:"ofType,omitempty"`
}

func (r TypeRef) IsScalar() bool {
	for ref := &r; ref != nil; ref = ref.OfType {
		if ref.Kind == TypeKindScalar {
			return true
		}
	}
	return false
}

func (r TypeRef) IsCustomScalar() bool {
	for ref := &r; ref != nil; ref = ref.OfType {
		if ref.Kind == TypeKindScalar {
			switch Scalar(ref.Name) {
			case ScalarInt:
				return false
			case ScalarFloat:
				return false
			case ScalarString:
				return false
			case ScalarBoolean:
				return false
			}
			return true
		}
	}
	return false
}

func (r TypeRef) IsObject() bool {
	for ref := &r; ref != nil; ref = ref.OfType {
		if ref.Kind == TypeKindObject {
			return true
		}
	}
	return false
}

type InputValues []InputValue

type InputValue struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	DefaultValue *string  `json:"defaultValue"`
	TypeRef      *TypeRef `json:"type"`
}
