package graphql

import "regexp"

// scalarOrder is the fixed emission order of the GraphQL scalar declarations.
var scalarOrder = []string{"DateTime", "BigInt", "WEI"}

// prismaGqlTypeMap is the inverse of gqlTypeMap used by ParseGQL for the field type position: a GraphQL type name to the Prisma type it maps back to.
var prismaGqlTypeMap = map[string]string{
	"Int":      "Int",
	"String":   "String",
	"Boolean":  "Boolean",
	"Float":    "Float",
	"BigInt":   "BigInt",
	"DateTime": "DateTime",
	"Json":     "Json",
	"Decimal":  "Decimal",
	"Bytes":    "Bytes",
}

// prismaArgTypeMap is the inverse of gqlTypeMap used by ParseGQL for argument types: a GraphQL argument type to the lowercase Prisma type it maps back to.
var prismaArgTypeMap = map[string]string{
	"Int":      "int",
	"String":   "string",
	"Boolean":  "boolean",
	"Float":    "float",
	"BigInt":   "bigint",
	"DateTime": "datetime",
	"Wei":      "wei",
	"WEI":      "wei",
	"Json":     "json",
	"Decimal":  "decimal",
	"Bytes":    "bytes",
}

// gqlTypeMap mirrors the reference GqlTypeMap: a lowercase Prisma type name to the GraphQL type it renders as. "bytes" is a minimal extension over the reference so the project's own schemas (which use Bytes) parse.
var gqlTypeMap = map[string]string{
	"int":      "Int",
	"string":   "String",
	"boolean":  "Boolean",
	"float":    "Float",
	"bigint":   "BigInt",
	"datetime": "DateTime",
	"json":     "String",
	"decimal":  "String",
	"bytes":    "String",
	"wei":      "WEI",
}

var (
	modelStartRe = regexp.MustCompile(`^\s*model.*`)
	enumStartRe  = regexp.MustCompile(`^\s*enum.*`)
	modelEndRe   = regexp.MustCompile(`^\s*}.*`)
	commentLine  = regexp.MustCompile(`^\s*\/\/.+$`)
	attributeRe  = regexp.MustCompile(`^\s*@@.*`)
	// a "//@@ property" line declares a field as if it were a normal property.
	propertyDirectiveRe = regexp.MustCompile(`^\s.\/\/@@ property`)
	// skip a field or enum option whose name starts with whitespace or contains "//".
	skipNameRe = regexp.MustCompile(`^\s|\/\/`)
)
