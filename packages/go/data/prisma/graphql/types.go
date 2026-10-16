package graphql

import (
	GT "github.com/verikono/glocalize/types/go"
)

type (
	GraphQLSchema GT.GraphQLSchema
	PrismaSchema  = GT.PrismaSchema
)

// ParsablePropertyShape is the parsed form of a single Prisma model field.
type ParsablePropertyShape struct {
	Name       string
	Type       string
	Optional   bool
	Multi      bool
	Params     []string
	Directives []string
	Docblocks  []string
}

// ParsableModelShape is the parsed form of a Prisma model block.
type ParsableModelShape struct {
	Name  string
	Props []ParsablePropertyShape
	Docs  []string
}

// ParsableEnumShape is the parsed form of a Prisma enum block.
type ParsableEnumShape struct {
	Name    string
	Options []string
	Docs    []string
}

// parserContext carries the extracted models and enums through parsing.
type parserContext struct {
	models []*ParsableModelShape
	enums  []*ParsableEnumShape
}

// gqlBlock is a parsed GraphQL enum or type block.
type gqlBlock struct {
	kind    string // "model" or "enum"
	name    string
	docs    []string
	options []string
	fields  []gqlField
}

// gqlField is a single field declaration inside a GraphQL type block.
type gqlField struct {
	docs []string
	line string
}
